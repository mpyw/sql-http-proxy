# sql-http-proxy

[![CI](https://github.com/mpyw/sql-http-proxy/actions/workflows/ci.yml/badge.svg)](https://github.com/mpyw/sql-http-proxy/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/mpyw/sql-http-proxy/graph/badge.svg)](https://codecov.io/gh/mpyw/sql-http-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/mpyw/sql-http-proxy)](https://goreportcard.com/report/github.com/mpyw/sql-http-proxy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A YAML-driven HTTP server that maps endpoints to SQL queries. Define your API in YAML, run the server, and get a working REST API.

> [!NOTE]
> This project was written by AI (Claude Code).

## Installation

```bash
go install github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest
```

### Build Tags

By default, all database drivers are included. Use build tags for smaller binaries:

```bash
# PostgreSQL only
go build -tags postgres ./cmd/sql-http-proxy

# SQLite only
go build -tags sqlite ./cmd/sql-http-proxy

# Multiple drivers
go build -tags postgres,mysql ./cmd/sql-http-proxy

# Mock only (no database)
go build -tags mock ./cmd/sql-http-proxy
```

## Quick Start

### 1. Create Configuration

Create `.sql-http-proxy.yaml`:

```yaml
dsn: postgres://user:pass@localhost:5432/mydb?sslmode=disable

queries:
  - type: many
    path: /users
    sql: SELECT * FROM users

  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
```

### 2. Start Server

```bash
sql-http-proxy -l :8080
```

### 3. Make Requests

```bash
# List all users
curl http://localhost:8080/users

# Get single user (404 if not found)
curl http://localhost:8080/user?id=1
```

### Mock Mode (No Database)

Run without a database using mock data:

```yaml
queries:
  - type: many
    path: /users
    sql: SELECT * FROM users  # ignored
    transform:
      mock:
        json:
          - { id: 1, name: Alice }
          - { id: 2, name: Bob }
```

---

## Configuration Reference

See [sql-http-proxy.example.yaml](sql-http-proxy.example.yaml) for a comprehensive example.

### Top-Level Options

| Option | Type | Description |
|--------|------|-------------|
| `dsn` | string | Database connection string (or use `SQL_PROXY_DSN` env var) |
| `global_helpers` | object/string | JavaScript helpers available in all transforms |
| `csv` | object | CSV parsing options |
| `queries` | array | Query endpoints (SELECT) |
| `mutations` | array | Mutation endpoints (INSERT/UPDATE/DELETE) |

#### DSN Examples

| Database | DSN |
|----------|-----|
| PostgreSQL | `postgres://user:pass@localhost:5432/db?sslmode=disable` |
| MySQL | `mysql://user:pass@tcp(localhost:3306)/db` |
| SQLite | `file:./data.db` or `sqlite:./data.db` |
| SQL Server | `sqlserver://user:pass@localhost:1433?database=db` |

#### Global Helpers

Define JavaScript functions available in all `pre`, `mock`, `post` transforms and `csv.value_parser`:

```yaml
# Shorthand (inline JS only)
global_helpers: |
  function validate(x) { ... }

# Full form
global_helpers:
  js: |
    function validate(x) { ... }
  js_files:
    - ./helpers/utils.js
```

#### CSV Value Parser

Custom parsing for CSV mock data:

```yaml
csv:
  value_parser: |
    # Parse numbers, booleans, dates
    if (value === 'true') return true;
    if (value === 'false') return false;
    if (/^\d+$/.test(value)) return parseInt(value);
    return value;
```

---

### Queries

Query endpoints handle SELECT operations.

```yaml
queries:
  - type: one|many      # Required
    path: /endpoint     # Required
    sql: SELECT ...     # Required
    method: GET         # Optional (default: GET)
    accepts: json       # Optional (default: [json, form])
    handle_not_found: true  # Optional, type: one only
    transform: { ... }  # Optional
```

| Property | Type | Description |
|----------|------|-------------|
| `type` | `one` \| `many` | `one`: single row (404 if not found), `many`: array |
| `path` | string | HTTP endpoint path (must start with `/`) |
| `sql` | string | SQL query with `:name` placeholders |
| `method` | string | HTTP method (GET, POST, PUT, PATCH, DELETE) |
| `accepts` | string/array | Accepted Content-Types: `json`, `form`, or both |
| `handle_not_found` | boolean | Pass `null` to post-transform instead of 404 |
| `transform` | object | Pre/mock/post transforms |

---

### Mutations

Mutation endpoints handle INSERT/UPDATE/DELETE operations.

```yaml
mutations:
  - type: one|many|none # Required
    method: POST        # Optional (default: POST)
    path: /endpoint     # Required
    sql: INSERT ...     # Required
    accepts: json       # Optional (default: [json, form])
    transform: { ... }  # Optional
```

| Property | Type | Description |
|----------|------|-------------|
| `type` | `one` \| `many` \| `none` | Return type (`none` = 204 No Content) |
| `method` | string | HTTP method (POST, PUT, PATCH, DELETE) |
| `path` | string | HTTP endpoint path |
| `sql` | string | SQL with `:name` placeholders (use `RETURNING *` for `one`/`many`) |
| `accepts` | string/array | Accepted Content-Types |
| `transform` | object | Pre/mock/post transforms |

---

### Transform

Transforms allow JavaScript processing at different stages.

```yaml
transform:
  pre: |          # Before SQL execution
    ...
  mock: { ... }   # Replace SQL with mock data
  post: |         # After SQL execution
    ...
```

#### Pre-Transform

Validate and transform input parameters before SQL execution.

**Available variables:**
- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state object (persists to post-transform)
- `sql` (free variable): SQL string (can be modified)

**Must return:** Object with parameters for SQL

```yaml
transform:
  pre: |
    // Validate
    if (!input.email.includes('@')) {
      throw { status: 400, body: { error: 'invalid email' } };
    }

    // Store state for post-transform
    ctx.requestTime = Date.now();

    // Modify SQL dynamically
    if (input.active) {
      sql += ' AND active = true';
    }

    // Return parameters for SQL
    return { email: input.email.toLowerCase() };
```

#### Mock

Return mock data without database execution. Only one source type can be specified.

```yaml
# JavaScript (inline)
mock: |
  return { id: 1, name: 'Mock' };

# JavaScript (file)
mock:
  js_file: ./mocks/data.js

# CSV (inline)
mock:
  csv: |
    id,name,active
    1,Alice,true
    2,Bob,false

# CSV (file)
mock:
  csv_file: ./data.csv

# JSON (inline)
mock:
  json:
    - { id: 1, name: Alice }
    - { id: 2, name: Bob }

# JSON (file)
mock:
  json_file: ./data.json

# JSONL (inline)
mock:
  jsonl: |
    {"id": 1, "name": "Alice"}
    {"id": 2, "name": "Bob"}

# JSONL (file)
mock:
  jsonl_file: ./data.jsonl
```

**Mock JS variables:**
- `input` (parameter): Request parameters (after pre-transform)
- `ctx` (free variable): Shared state
- `sql` (free variable): SQL string (read-only in mock)

#### Post-Transform

Transform the result after SQL/mock execution.

**Available variables:**
- `input` (parameter): Original request parameters
- `output` (parameter): Query result (row or array)
- `ctx` (free variable): Shared state from pre-transform

**For `type: one`:**

```yaml
transform:
  post: |
    return {
      ...output,
      fullName: output.first_name + ' ' + output.last_name
    };
```

**For `type: many`:**

```yaml
# Transform each row
transform:
  post:
    each: |
      return { ...output, processed: true };

# Transform entire array
transform:
  post:
    all: |
      return { data: output, count: output.length };

# Both (each runs first)
transform:
  post:
    each: |
      return { ...output, processed: true };
    all: |
      return { items: output, total: output.length };
```

#### Error Handling

Throw an object with `status` and `body` to return a custom HTTP error:

```yaml
transform:
  pre: |
    if (!input.token) {
      throw { status: 401, body: { error: 'unauthorized' } };
    }
    return input;
```

| Phase | Default Status |
|-------|---------------|
| `pre` | 400 Bad Request |
| `mock` | 500 Internal Server Error |
| `post` | 500 Internal Server Error |

---

### Named Placeholders

Use `:name` syntax for SQL parameters:

```yaml
sql: SELECT * FROM users WHERE id = :id AND status = :status
```

Parameters come from:
- **GET requests:** Query string (`?id=1&status=active`)
- **POST/PUT/PATCH/DELETE:** Request body (JSON or form)

---

## Full Example

See [sql-http-proxy.example.yaml](sql-http-proxy.example.yaml) for a comprehensive example demonstrating all features.

## License

MIT
