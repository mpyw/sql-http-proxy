# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

**sql-http-proxy** is a YAML configuration-based HTTP to SQL proxy server. It allows you to define SQL queries in a YAML config file and expose them as HTTP endpoints.

## Architecture

```
sql-http-proxy/
├── cmd/sql-http-proxy/main.go     # Entry point
│
├── internal/
│   ├── cli/commands/
│   │   ├── app.go                 # urfave/cli v3 app definition
│   │   └── drivers_*.go           # Database driver imports (build tags)
│   ├── config/                    # Configuration parsing
│   └── server/                    # HTTP server and handlers
│
└── go.mod
```

## Development Commands

Requires Go 1.27 or later: the project uses `encoding/json/v2` / `encoding/json/jsontext`
and generic methods.

```bash
# Build (all drivers)
go build ./cmd/sql-http-proxy

# Build with specific drivers
go build -tags postgres ./cmd/sql-http-proxy
go build -tags postgres,mysql ./cmd/sql-http-proxy

# Run
./sql-http-proxy -c config.yaml -l :8080

# Test
go test ./...
```

## Configuration

Example `sql-http-proxy.yaml`:

```yaml
database:
  dsn: postgres://user:pass@localhost:5432/db?sslmode=disable

http:
  cors: true  # Enable permissive CORS, or use object for detailed config

queries:
  - type: one
    path: /users/by-id
    sql: SELECT * FROM users WHERE id = :id

  - type: many
    path: /users
    sql: SELECT * FROM users LIMIT :limit
```

- `database.dsn`: Database connection string (supports `${VAR}` env expansion)
- `http.cors`: CORS config - `true` for permissive, or object with `allowed_origins`, `allow_credentials`, `max_age`
- `queries[].type`: `one` for single row, `many` for multiple rows
- `queries[].path`: HTTP endpoint path
- `queries[].sql`: SQL query with named placeholders (`:name`)

Query parameters are passed via URL: `/users/by-id?id=123`

## Supported Databases

- PostgreSQL (`postgres`, `postgresql`)
- MySQL (`mysql`)
- SQLite (`file`, `sqlite`)
- SQL Server (`sqlserver`)

## Code Style

- Follow standard Go conventions
- Use `urfave/cli/v3` for CLI structure
- Error messages should be user-friendly
