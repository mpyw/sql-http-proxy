# sql-http-proxy

[![Test](https://github.com/mpyw/sql-http-proxy/actions/workflows/test.yml/badge.svg)](https://github.com/mpyw/sql-http-proxy/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/mpyw/sql-http-proxy/graph/badge.svg)](https://codecov.io/gh/mpyw/sql-http-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/mpyw/sql-http-proxy)](https://goreportcard.com/report/github.com/mpyw/sql-http-proxy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A YAML-driven HTTP server that maps endpoints to SQL queries. Define your API in YAML, run the server, and get a working REST API.

> [!NOTE]
> This project was written by AI (Claude Code).

# Installation

## Homebrew (macOS / Linux)

```bash
brew install mpyw/tap/sql-http-proxy
```

## Scoop (Windows)

```powershell
scoop bucket add mpyw https://github.com/mpyw/scoop-bucket
scoop install sql-http-proxy
```

## Debian / Ubuntu

```bash
# Download the latest .deb package from GitHub Releases
curl -LO https://github.com/mpyw/sql-http-proxy/releases/latest/download/sql-http-proxy_VERSION-1_amd64.deb
sudo dpkg -i sql-http-proxy_VERSION-1_amd64.deb
```

## RHEL / Fedora

```bash
# Download the latest .rpm package from GitHub Releases
curl -LO https://github.com/mpyw/sql-http-proxy/releases/latest/download/sql-http-proxy-VERSION-1.x86_64.rpm
sudo rpm -i sql-http-proxy-VERSION-1.x86_64.rpm
```

## Binary Download (Linux / macOS)

```bash
# Set version and architecture
export VERSION=1.0.0
export ARCH=amd64  # or arm64

# Linux
curl -L "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy_${VERSION}_linux_${ARCH}.tar.gz" | tar xz
sudo mv sql-http-proxy /usr/local/bin/

# macOS
curl -L "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy_${VERSION}_darwin_${ARCH}.tar.gz" | tar xz
sudo mv sql-http-proxy /usr/local/bin/
```

Or download manually from [GitHub Releases](https://github.com/mpyw/sql-http-proxy/releases).

## Go Install

```bash
go install github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest
```

### Build Tags

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
dsn: postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST:-localhost}:${DB_PORT:-5432}/mydb

queries:
  - type: many
    path: /users
    sql: SELECT * FROM users

  - type: one
    path: /users/{id:*uuid_v7*}
    sql: SELECT * FROM users WHERE id = :id
```

## 2. Start Server

```bash
sql-http-proxy -l :8080
```

## 3. Make Requests

```bash
curl http://localhost:8080/users                                        # List all users
curl http://localhost:8080/users/019411a5-3d7f-7000-8000-000000000001   # Get user by ID (UUIDv7)
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
dsn: postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST:-localhost}:${DB_PORT:-5432}/mydb
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

## Path Parameters

Use `{param}` syntax in paths to capture URL segments (chi router syntax):

```yaml
queries:
  - type: one
    path: /users/{id:*uuid_v7*}
    sql: SELECT * FROM users WHERE id = :id

  - type: many
    path: /users/{user_id:*uuid_v7*}/posts
    sql: SELECT * FROM posts WHERE user_id = :user_id

mutations:
  - type: one
    method: PUT
    path: /users/{id:*uuid_v7*}
    sql: UPDATE users SET name = :name WHERE id = :id RETURNING *
```

Path parameters take priority over query string and body parameters.

**Regex shorthands** for common validation patterns:

| Shorthand | Description |
|-----------|-------------|
| `*uuid*` | Any UUID (lowercase) |
| `*uuid_v4*` | UUIDv4 only |
| `*uuid_v7*` | UUIDv7 only |

**Custom regex** (returns 404 if not matched):

```yaml
# Numeric ID only
path: /posts/{id:[0-9]+}

# Slug format
path: /articles/{slug:[a-z0-9-]+}
```

See [SCHEMA.md](SCHEMA.md#path-parameters) for full regex syntax.

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
global_helpers: |
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

Mock sources are divided into two categories based on the endpoint type.

### Object Sources (for `type: one`)

Return a single object directly:

```yaml
mock:
  object: { id: 1, name: Alice }          # Inline YAML object
  # or
  object_json: '{"id": 1, "name": "Alice"}'  # JSON string
  # or
  object_json_file: ./data/user.json      # External JSON file
  # or
  object_js: |                            # JavaScript (dynamic)
    return { id: parseInt(input.id), name: 'User ' + input.id };
```

### Array Sources (for `type: many`)

Return multiple rows:

```yaml
mock:
  array:                                  # Inline YAML array
    - { id: 1, name: Alice }
    - { id: 2, name: Bob }
  # or
  array_json: '[{"id": 1}, {"id": 2}]'   # JSON string
  array_json_file: ./data/users.json     # External JSON file
  # or
  csv: |                                  # CSV with header
    id,name
    1,Alice
    2,Bob
  csv_file: ./data/users.csv             # External CSV file
  # or
  jsonl: |                                # JSON Lines
    {"id": 1, "name": "Alice"}
    {"id": 2, "name": "Bob"}
  jsonl_file: ./data/users.jsonl         # External JSONL file
  # or
  array_js: |                             # JavaScript (dynamic)
    return [{ id: 1, name: 'Alice' }, { id: 2, name: 'Bob' }];
```

### Array Sources with Filter (for `type: one`)

Use `filter` to select a single row from array data:

```yaml
- type: one
  path: /user
  mock:
    array:
      - { id: 1, name: Alice }
      - { id: 2, name: Bob }
    filter: return row.id === parseInt(input.id)
```

The filter receives `row` and `input`, returns `true` to include. First matching row is returned (404 if none match).

See [SCHEMA.md](SCHEMA.md#mock) for complete mock reference.

# License

MIT
