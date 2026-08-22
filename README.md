<div align="center">
  <h1>sql-http-proxy</h1>
  <p><strong>Your SQL, served over HTTP — defined entirely in YAML</strong></p>

  [![Go Reference](https://pkg.go.dev/badge/github.com/mpyw/sql-http-proxy.svg)](https://pkg.go.dev/github.com/mpyw/sql-http-proxy)
  [![Test](https://github.com/mpyw/sql-http-proxy/actions/workflows/test.yml/badge.svg)](https://github.com/mpyw/sql-http-proxy/actions/workflows/test.yml)
  [![Codecov](https://codecov.io/gh/mpyw/sql-http-proxy/graph/badge.svg)](https://codecov.io/gh/mpyw/sql-http-proxy)
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
</div>

> [!NOTE]
> This project was written by AI (Claude Code).

A YAML-driven HTTP server that maps endpoints to SQL queries against <a href="https://www.postgresql.org/"><img src="https://cdn.simpleicons.org/postgresql" height="16" alt=""></a> PostgreSQL, <a href="https://www.mysql.com/"><img src="https://cdn.simpleicons.org/mysql" height="16" alt=""></a> MySQL, <a href="https://www.sqlite.org/"><img src="https://cdn.simpleicons.org/sqlite" height="16" alt=""></a> SQLite, and <a href="https://www.microsoft.com/sql-server">SQL Server</a>. Define your API in YAML, run the server, and get a working REST API — with **named parameters**, **path parameters**, a **JavaScript transform pipeline**, and a **mock mode** that needs no database at all.

## Installation

### <a href="https://mise.jdx.dev/"><img src="https://mise.jdx.dev/logo.svg" height="28" alt=""></a> Using [mise](https://mise.jdx.dev/) (macOS/Linux/Windows)

sql-http-proxy is installable directly from GitHub Releases via mise's `github` backend — no extra registry required:

```bash
mise use -g "github:mpyw/sql-http-proxy"
```

Or pin it per project in `mise.toml`:

```toml
[tools]
"github:mpyw/sql-http-proxy" = "latest"
```

### <a href="https://brew.sh/"><img src="https://cdn.simpleicons.org/homebrew" height="28" alt=""></a> Using [Homebrew](https://brew.sh/) (macOS/Linux)

```bash
brew install mpyw/tap/sql-http-proxy
```

### <a href="https://scoop.sh/"><img src="https://github.com/ScoopInstaller.png?size=64" height="28" alt=""></a> Using [Scoop](https://scoop.sh/) (Windows)

```powershell
scoop bucket add mpyw https://github.com/mpyw/scoop-bucket.git
scoop install sql-http-proxy
```

<details>
<summary><a href="https://www.linux.org/"><img src="https://upload.wikimedia.org/wikipedia/commons/a/af/Tux.png" height="20" alt=""></a> Manual .deb / .rpm install</summary>

Native packages are tracked by `apt`/`dnf`, so upgrades and removal stay clean. Download from [GitHub Releases](https://github.com/mpyw/sql-http-proxy/releases) and hand the file to your package manager:

**Debian/Ubuntu (.deb):**

```bash
export VERSION=0.0.0
export ARCH=amd64  # or arm64

curl -LO "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy_${VERSION}-1_${ARCH}.deb"
sudo apt install "./sql-http-proxy_${VERSION}-1_${ARCH}.deb"
```

**Red Hat/Fedora (.rpm):**

```bash
export VERSION=0.0.0
export ARCH=x86_64  # or aarch64

curl -LO "https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}/sql-http-proxy-${VERSION}-1.${ARCH}.rpm"
sudo dnf install "./sql-http-proxy-${VERSION}-1.${ARCH}.rpm"
```

</details>

<details>
<summary><a href="https://curl.se/"><img src="https://cdn.simpleicons.org/curl" height="20" alt=""></a> Downloading the tarball directly (macOS/Linux/Windows)</summary>

No package manager? Grab the archive for your platform from [GitHub Releases](https://github.com/mpyw/sql-http-proxy/releases):

```bash
export VERSION=0.0.0
export OS=linux    # or darwin
export ARCH=amd64  # or arm64
export BASE_URL="https://github.com/mpyw/sql-http-proxy/releases/download/v${VERSION}"

# Download the archive and the release's checksum list
curl -LO "${BASE_URL}/sql-http-proxy_${VERSION}_${OS}_${ARCH}.tar.gz"
curl -LO "${BASE_URL}/checksums.txt"

# Verify before installing (use `shasum -a 256 -c` on macOS)
sha256sum --ignore-missing -c checksums.txt

tar xzf "sql-http-proxy_${VERSION}_${OS}_${ARCH}.tar.gz"
sudo mv sql-http-proxy /usr/local/bin/
```

On Windows, download `sql-http-proxy_${VERSION}_windows_${ARCH}.zip` and extract `sql-http-proxy.exe` somewhere on your `PATH`.

</details>

<details>
<summary><a href="https://go.dev/"><img src="https://cdn.simpleicons.org/go" height="20" alt=""></a> Using <code>go install</code></summary>

Requires Go 1.27+.

```bash
go install github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest
```

**Build tags** trade drivers for a smaller binary. With no tags, every driver is included:

```bash
go install -tags postgres github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest  # PostgreSQL only
go install -tags sqlite   github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest  # SQLite only
go install -tags mock     github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest  # Mock only (no database)
```

Tags combine: `-tags postgres,mysql`.

</details>

<details>
<summary><a href="https://go.dev/"><img src="https://cdn.simpleicons.org/go" height="20" alt=""></a> Using <code>go tool</code> (Go 1.27+)</summary>

Pin sql-http-proxy to a project rather than the machine:

```bash
# Add to go.mod as a tool dependency
go get -tool github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest

# Run via go tool
go tool sql-http-proxy -c .sql-http-proxy.yaml -l :8080
```

</details>

<details>
<summary><a href="#"><img src="https://cdn.simpleicons.org/git" height="20" alt=""></a> Building from source</summary>

For platforms without pre-built packages, or to run an unreleased revision.

Requires Go 1.27+ — the project uses `encoding/json/v2`, `encoding/json/jsontext`, and generic methods.

```bash
git clone https://github.com/mpyw/sql-http-proxy.git
cd sql-http-proxy

go build ./cmd/sql-http-proxy                    # all drivers
go build -tags postgres,mysql ./cmd/sql-http-proxy  # selected drivers only
```

</details>

## Features

- **YAML-defined endpoints**: one entry per endpoint, `sql` or `mock`, nothing to compile
- **Named parameters**: `:name` placeholders bound from query string, request body, or path
- **Path parameters**: [chi](https://github.com/go-chi/chi) routing with regex shorthands like `{id:*uuid_v7*}`
- **Transform pipeline**: JavaScript `pre`/`post` hooks to validate input, rewrite SQL, and reshape output
- **Mock mode**: serve YAML, JSON, JSONL, CSV, or JavaScript fixtures with no database connection
- **Multi-database**: [PostgreSQL](https://www.postgresql.org/), [MySQL](https://www.mysql.com/), [SQLite](https://www.sqlite.org/), and [SQL Server](https://www.microsoft.com/sql-server), selectable at build time
- **Strict JSON**: request bodies are parsed per RFC 7493 — duplicate keys and invalid UTF-8 are rejected rather than silently binding

### Supported Databases

The driver is chosen from the DSN scheme, so nothing else needs configuring:

| Database | DSN scheme | Build tag |
|----------|-----------|-----------|
| [PostgreSQL](https://www.postgresql.org/) | `postgres://`, `postgresql://` | `postgres` |
| [MySQL](https://www.mysql.com/) | `mysql://` | `mysql` |
| [SQLite](https://www.sqlite.org/) | `file:`, `sqlite:` — or a bare `:memory:`, `*.db`, `*.sqlite` path | `sqlite` |
| [SQL Server](https://www.microsoft.com/sql-server) | `sqlserver://` | `mssql` |

Omit the build tag to get every driver; use `-tags mock` for a binary with none of them.

## Quick Start

### 1. Create Configuration

Create `.sql-http-proxy.yaml`:

```yaml
database:
  dsn: postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST:-localhost}:${DB_PORT:-5432}/mydb

queries:
  - type: many
    path: /users
    sql: SELECT * FROM users

  - type: one
    path: /users/{id:*uuid_v7*}
    sql: SELECT * FROM users WHERE id = :id
```

### 2. Start Server

```bash
sql-http-proxy -l :8080
```

### 3. Make Requests

```bash
curl http://localhost:8080/users                                        # List all users
curl http://localhost:8080/users/019411a5-3d7f-7000-8000-000000000001   # Get user by ID (UUIDv7)
```

## Mock Mode

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

> [!NOTE]
> Returning `null` or `undefined` from mock JS results in 404 Not Found for `type: one`.

## Configuration Overview

See [SCHEMA.md](SCHEMA.md) for the complete reference and [sql-http-proxy.example.yaml](sql-http-proxy.example.yaml) for examples.

### Database & HTTP Configuration

Database connection with `${VAR}` environment variable expansion:

```yaml
database:
  dsn: postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST:-localhost}:${DB_PORT:-5432}/mydb

http:
  cors: true  # or: { allowed_origins: [...], allow_credentials: true, max_age: 86400 }
```

### Queries & Mutations

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

### Path Parameters

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

### Transform Pipeline

```yaml
transform:
  pre: |
    # Validate input, modify SQL
    return { id: parseInt(input.id) };

  post: |
    # Transform output
    return { ...output, formatted: true };
```

### Global Helpers

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

### Mock Sources

Mock sources are divided into two categories based on the endpoint type.

#### Object Sources (for `type: one`)

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

#### Array Sources (for `type: many`)

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

#### Array Sources with Filter (for `type: one`)

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

See [SCHEMA.md](SCHEMA.md#mock) for the complete mock reference.

## Development

Requires Go 1.27+.

```bash
# Build (all drivers)
go build ./cmd/sql-http-proxy

# Build with specific drivers
go build -tags postgres,mysql ./cmd/sql-http-proxy

# Run
./sql-http-proxy -c .sql-http-proxy.yaml -l :8080

# Unit tests
go test ./internal/...

# E2E tests
go test ./e2e/...

# Everything, with the race detector
go test -race ./...

# Lint (needs golangci-lint v2.13.0+ to parse generic methods)
golangci-lint run
```

## License

MIT License
