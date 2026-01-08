# sql-http-proxy

[![CI](https://github.com/mpyw/sql-http-proxy/actions/workflows/ci.yml/badge.svg)](https://github.com/mpyw/sql-http-proxy/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/mpyw/sql-http-proxy/graph/badge.svg)](https://codecov.io/gh/mpyw/sql-http-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/mpyw/sql-http-proxy)](https://goreportcard.com/report/github.com/mpyw/sql-http-proxy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A YAML-driven HTTP server that maps endpoints to SQL queries. Define your API in YAML, run the server, and get a working REST API.

> [!NOTE]
> This project was written by AI (Claude Code).

# Installation

```bash
go install github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest
```

## Build Tags

Use build tags for smaller binaries:

```bash
go install -tags postgres github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest  # PostgreSQL only
go install -tags sqlite github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest    # SQLite only
go install -tags mock github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest      # Mock only (no database)
```

# Quick Start

## 1. Create Configuration

Create `.sql-http-proxy.yaml`:

```yaml
dsn: postgres://${DB_USER}:${DB_PASSWORD}@localhost:5432/mydb

queries:
  - type: many
    path: /users
    sql: SELECT * FROM users

  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
```

## 2. Start Server

```bash
sql-http-proxy -l :8080
```

## 3. Make Requests

```bash
curl http://localhost:8080/users      # List all users
curl http://localhost:8080/user?id=1  # Get single user
```

# Mock Mode

Run without a database using mock data. Use `mock` instead of `sql`:

```yaml
queries:
  # Array source for type: many
  - type: many
    path: /users
    mock:
      array:
        - { id: 1, name: Alice }
        - { id: 2, name: Bob }

  # JavaScript for dynamic data
  - type: one
    path: /user
    mock:
      object_js: |
        if (input.id === '404') return null;
        return { id: parseInt(input.id), name: 'User ' + input.id };

  # Array with filter for type: one
  - type: one
    path: /user-by-email
    mock:
      array:
        - { id: 1, name: Alice, email: alice@example.com }
        - { id: 2, name: Bob, email: bob@example.com }
      filter: return row.email === input.email
```

> Returning `null` or `undefined` from mock JS results in 404 Not Found for `type: one`.

# Configuration Overview

See [SCHEMA.md](SCHEMA.md) for complete reference and [sql-http-proxy.example.yaml](sql-http-proxy.example.yaml) for examples.

## DSN

Supports `${VAR}` environment variable expansion:

```yaml
dsn: postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/mydb
```

## Queries & Mutations

Each endpoint uses either `sql` OR `mock`, not both:

```yaml
queries:
  - type: one|many
    path: /endpoint
    sql: SELECT * FROM table WHERE id = :id  # OR mock: { ... }
    transform: { ... }

mutations:
  - type: one|many|none
    method: POST
    path: /endpoint
    sql: INSERT INTO table (...) RETURNING *  # OR mock: { ... }
    transform: { ... }
```

## Transform Pipeline

```yaml
transform:
  pre: |
    # Validate input, modify SQL
    return { id: parseInt(input.id) };

  post: |
    # Transform output
    return { ...output, formatted: true };
```

## Global Helpers

Reusable JavaScript functions for all transforms:

```yaml
global_helpers:
  js: |
    function requireInt(val, name) {
      const n = parseInt(val);
      if (isNaN(n)) throw { status: 400, body: { error: name + ' required' } };
      return n;
    }

queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      pre: |
        return { id: requireInt(input.id, 'id') };
```

## Mock Sources

| Source | Type | Example |
|--------|------|---------|
| `object` | one | `mock: { object: { id: 1, name: Alice } }` |
| `object_js` | one | `mock: { object_js: \| return { id: 1 }; }` |
| `array` | many | `mock: { array: [{ id: 1 }, { id: 2 }] }` |
| `array_js` | many | `mock: { array_js: \| return [{ id: 1 }]; }` |
| `csv` | many | `mock: { csv: \| id,name \n 1,Alice }` |
| `jsonl` | many | `mock: { jsonl: \| {"id":1} \n {"id":2} }` |

For `type: one` with array sources, use `filter` to select the matching row.

# License

MIT
