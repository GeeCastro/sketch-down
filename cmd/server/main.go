package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sketch-down/sketch-down/internal/diagram"
	"github.com/sketch-down/sketch-down/internal/project"
	"github.com/sketch-down/sketch-down/internal/validation"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbDir := os.Getenv("DATA_DIR")
	if dbDir == "" {
		dbDir = "data"
	}
	dbPath := filepath.Join(dbDir, "sketch-down.db")

	store, err := project.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer store.Close()

	// Ensure active project exists.
	active, err := store.GetActiveProject()
	if err != nil {
		log.Fatalf("get active: %v", err)
	}
	if active == nil {
		if _, err := store.CreateProject("Default Project"); err != nil {
			log.Fatalf("create default project: %v", err)
		}
		log.Println("created default project")
	}

	mux := http.NewServeMux()
	registerAPI(mux, store)

	// Serve frontend static files from web/dist.
	distDir := "web/dist"
	if env := os.Getenv("WEB_DIST"); env != "" {
		distDir = env
	}
	fs := http.FileServer(http.Dir(distDir))
	mux.Handle("/", fs)

	handler := corsMiddleware(mux)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func registerAPI(mux *http.ServeMux, store project.Store) {
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/project", func(w http.ResponseWriter, r *http.Request) {
		p, err := store.GetActiveProject()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if p == nil {
			p, err = store.CreateProject("Default Project")
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
		}
		writeJSON(w, http.StatusOK, p)
	})

	mux.HandleFunc("PUT /api/project/diagram", func(w http.ResponseWriter, r *http.Request) {
		p, err := store.GetActiveProject()
		if err != nil || p == nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("no active project"))
			return
		}

		var d diagram.Diagram
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		if err := store.SaveDiagram(p.ID, &d); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
	})

	mux.HandleFunc("POST /api/project/export", func(w http.ResponseWriter, r *http.Request) {
		p, err := store.GetActiveProject()
		if err != nil || p == nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("no active project"))
			return
		}

		data, err := store.ExportProject(p.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.json"`, p.Name))
		w.Write(data)
	})

	mux.HandleFunc("POST /api/project/import", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(io.LimitReader(r.Body, 10<<20)) // 10MB limit
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		p, err := store.ImportProject(data)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, p)
	})

	mux.HandleFunc("GET /api/project/validate", func(w http.ResponseWriter, r *http.Request) {
		p, err := store.GetActiveProject()
		if err != nil || p == nil {
			writeError(w, http.StatusInternalServerError, fmt.Errorf("no active project"))
			return
		}

		warnings := validation.Validate(p.Diagram)
		if warnings == nil {
			warnings = []validation.Warning{}
		}
		writeJSON(w, http.StatusOK, map[string]any{"warnings": warnings})
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
