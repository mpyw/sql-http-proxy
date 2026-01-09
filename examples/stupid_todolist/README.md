# Stupid Todo List

A minimal todo list application demonstrating sql-http-proxy with SQLite in-memory database.

## Requirements

- sql-http-proxy binary
- A web browser

## How to Run

### 1. Start the API server

```bash
cd examples/stupid_todolist
sql-http-proxy -c config.yaml -l :8090
```

The database is automatically initialized on startup via `database.init`.

### 2. Open the frontend

Open `index.html` in your browser:

```bash
# macOS
open index.html

# Linux
xdg-open index.html

# Windows
start index.html
```

Or use a simple HTTP server:

```bash
# Python 3
python -m http.server 3000

# Then open http://localhost:3000
```

## Architecture

```
Browser (index.html)
    |
    | fetch() with CORS
    v
sql-http-proxy (:8090)
    |
    | SQL
    v
SQLite (in-memory)
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/todos | List todos (supports `?q=search` for filtering) |
| GET | /api/todos/:id | Get single todo |
| POST | /api/todos | Create todo |
| PUT | /api/todos/:id | Update todo |
| DELETE | /api/todos/:id | Delete todo |

## Features Demonstrated

- **Database Auto-Init**: Uses `database.init` to create tables on startup
- **Dynamic LIKE Query**: Search with `?q=term` using pre-transform to build LIKE pattern
- **CORS Support**: Enabled via `http.cors: true`
- **Response Status**: Custom status codes (201 for create, 204 for delete)

## Notes

- Uses SQLite in-memory database (`file::memory:?cache=shared`)
- Data is lost when the server stops
