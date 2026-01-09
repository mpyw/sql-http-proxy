# Stupid Todo List

A minimal todo list application demonstrating sql-http-proxy with SQLite in-memory database.

## Requirements

- sql-http-proxy binary
- A web browser

## How to Run

### 1. Start the API server

```bash
cd examples/stupid_todolist
sql-http-proxy -c config.yaml -l :8080
```

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

### 3. Initialize the database

Click the "Initialize Database" button to create the todos table.

## Architecture

```
Browser (index.html)
    |
    | fetch() with CORS
    v
sql-http-proxy (:8080)
    |
    | SQL
    v
SQLite (in-memory)
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/init | Initialize database (create table) |
| GET | /api/todos | List all todos |
| GET | /api/todos/:id | Get single todo |
| POST | /api/todos | Create todo |
| PUT | /api/todos/:id | Update todo |
| DELETE | /api/todos/:id | Delete todo |

## Notes

- Uses SQLite in-memory database (`file::memory:?cache=shared`)
- Data is lost when the server stops
- CORS headers are added via `global_helpers` and `response.headers.set()`
