package project

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sketch-down/sketch-down/internal/diagram"
	_ "modernc.org/sqlite"
)

// Store defines the persistence interface for projects.
type Store interface {
	CreateProject(name string) (*diagram.Project, error)
	GetActiveProject() (*diagram.Project, error)
	SaveDiagram(projectID string, d *diagram.Diagram) error
	ListProjects() ([]diagram.Project, error)
	ExportProject(projectID string) ([]byte, error)
	ImportProject(data []byte) (*diagram.Project, error)
	Close() error
}

type sqliteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens or creates the SQLite database at dbPath.
func NewSQLiteStore(dbPath string) (Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(wal)")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return &sqliteStore{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id          TEXT PRIMARY KEY,
			name        TEXT NOT NULL,
			diagram_json TEXT NOT NULL,
			created_at  TEXT NOT NULL,
			updated_at  TEXT NOT NULL,
			is_active   INTEGER NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func (s *sqliteStore) CreateProject(name string) (*diagram.Project, error) {
	p := diagram.NewProject(name)

	djson, err := json.Marshal(p.Diagram)
	if err != nil {
		return nil, fmt.Errorf("marshal diagram: %w", err)
	}

	// Deactivate all, then insert new as active.
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE projects SET is_active = 0"); err != nil {
		return nil, fmt.Errorf("deactivate: %w", err)
	}

	_, err = tx.Exec(
		"INSERT INTO projects (id, name, diagram_json, created_at, updated_at, is_active) VALUES (?, ?, ?, ?, ?, 1)",
		p.ID, p.Name, string(djson),
		p.CreatedAt.Format(time.RFC3339),
		p.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("insert project: %w", err)
	}

	return p, tx.Commit()
}

func (s *sqliteStore) GetActiveProject() (*diagram.Project, error) {
	row := s.db.QueryRow(
		"SELECT id, name, diagram_json, created_at, updated_at FROM projects WHERE is_active = 1 LIMIT 1",
	)
	return scanProject(row)
}

func scanProject(row *sql.Row) (*diagram.Project, error) {
	var p diagram.Project
	var djson string
	var createdAt, updatedAt string
	err := row.Scan(&p.ID, &p.Name, &djson, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan project: %w", err)
	}

	p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	p.Diagram = &diagram.Diagram{}
	if err := json.Unmarshal([]byte(djson), p.Diagram); err != nil {
		return nil, fmt.Errorf("unmarshal diagram: %w", err)
	}
	return &p, nil
}

func (s *sqliteStore) SaveDiagram(projectID string, d *diagram.Diagram) error {
	djson, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("marshal diagram: %w", err)
	}
	now := time.Now().UTC().Format(time.RFC3339)

	res, err := s.db.Exec(
		"UPDATE projects SET diagram_json = ?, updated_at = ? WHERE id = ?",
		string(djson), now, projectID,
	)
	if err != nil {
		return fmt.Errorf("update diagram: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("project %s not found", projectID)
	}
	return nil
}

func (s *sqliteStore) ListProjects() ([]diagram.Project, error) {
	rows, err := s.db.Query(
		"SELECT id, name, diagram_json, created_at, updated_at FROM projects ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []diagram.Project
	for rows.Next() {
		var p diagram.Project
		var djson, createdAt, updatedAt string
		if err := rows.Scan(&p.ID, &p.Name, &djson, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		p.Diagram = &diagram.Diagram{}
		if err := json.Unmarshal([]byte(djson), p.Diagram); err != nil {
			return nil, fmt.Errorf("unmarshal diagram: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (s *sqliteStore) ExportProject(projectID string) ([]byte, error) {
	row := s.db.QueryRow(
		"SELECT id, name, diagram_json, created_at, updated_at FROM projects WHERE id = ?",
		projectID,
	)
	p, err := scanProject(row)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("project %s not found", projectID)
	}
	return json.MarshalIndent(p, "", "  ")
}

func (s *sqliteStore) ImportProject(data []byte) (*diagram.Project, error) {
	var p diagram.Project
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("unmarshal import: %w", err)
	}

	if p.ID == "" {
		p.ID = diagram.NewID()
	}
	if p.Diagram == nil {
		p.Diagram = diagram.NewDiagram()
	}
	now := time.Now().UTC()
	p.UpdatedAt = now
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}

	djson, err := json.Marshal(p.Diagram)
	if err != nil {
		return nil, fmt.Errorf("marshal diagram: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE projects SET is_active = 0"); err != nil {
		return nil, fmt.Errorf("deactivate: %w", err)
	}

	_, err = tx.Exec(
		`INSERT INTO projects (id, name, diagram_json, created_at, updated_at, is_active)
		 VALUES (?, ?, ?, ?, ?, 1)
		 ON CONFLICT(id) DO UPDATE SET
		   name = excluded.name,
		   diagram_json = excluded.diagram_json,
		   updated_at = excluded.updated_at,
		   is_active = 1`,
		p.ID, p.Name, string(djson),
		p.CreatedAt.Format(time.RFC3339),
		p.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("upsert project: %w", err)
	}

	return &p, tx.Commit()
}

func (s *sqliteStore) Close() error {
	return s.db.Close()
}
