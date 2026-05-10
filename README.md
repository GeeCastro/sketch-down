# Sketch Down

Local web app for boat DC electrical schematics. Visual canvas editor with structured electrical fields and validation warnings.

## Stack

- **Backend**: Go, SQLite (modernc.org/sqlite)
- **Frontend**: Svelte 5 + Vite (in `web/`)

## Quick Start

```bash
make dev     # run server on :8080
make build   # compile binary to bin/
make test    # run all Go tests
```

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/project` | Get active project (auto-creates default) |
| PUT | `/api/project/diagram` | Save/autosave diagram |
| POST | `/api/project/export` | Download project as JSON |
| POST | `/api/project/import` | Upload JSON, set as active |
| GET | `/api/project/validate` | Run validation on active diagram |

Data stored in `data/sketch-down.db`. Frontend served from `web/dist/`.
