# Configuration Schema Reference

This document describes all configuration options for sql-http-proxy. For usage examples, see the main [README](README.md) and [sql-http-proxy.example.yaml](sql-http-proxy.example.yaml).

## Top-Level Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `dsn` | string | No* | Database connection string. Supports `${VAR}`, `${VAR:-default}` env expansion |
| `global_helpers` | object/string | No | JavaScript helpers for all transforms |
| `csv` | object | No | CSV parsing options |
| `queries` | array | No | Query endpoints (SELECT) |
| `mutations` | array | No | Mutation endpoints (INSERT/UPDATE/DELETE) |

> *Required unless all endpoints use mock

## DSN

Database connection string with `${VAR}`, `$VAR`, or `${VAR:-default}` environment variable expansion.

```yaml
dsn: postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST:-localhost}:${DB_PORT:-5432}/mydb
```

### Driver Examples

| Database | DSN Format |
|----------|------------|
| PostgreSQL | `postgres://user:pass@localhost:5432/db?sslmode=disable` |
| MySQL | `mysql://user:pass@tcp(localhost:3306)/db` |
| SQLite | `file:./data.db` or `sqlite:./data.db` |
| SQL Server | `sqlserver://user:pass@localhost:1433?database=db` |

## Global Helpers

JavaScript functions available in all `pre`, `mock`, `post` transforms and `csv.value_parser`.

```yaml
# Shorthand (inline JS)
global_helpers: |
  function validate(x) { ... }

# Full form
global_helpers:
  js: |
    function validate(x) { ... }
  js_files:
    - ./helpers/utils.js
```

| Property | Type | Description |
|----------|------|-------------|
| `js` | string | Inline JavaScript code |
| `js_files` | string[] | Paths to JavaScript files (relative to config) |

## CSV Config

Custom value parsing for CSV mock data.

```yaml
csv:
  value_parser: |
    if (value === 'true') return true;
    if (value === 'false') return false;
    if (/^\d+$/.test(value)) return parseInt(value);
    return value;
```

| Property | Type | Description |
|----------|------|-------------|
| `value_parser` | string | JavaScript code for parsing CSV cell values |

The `value_parser` function receives `value` (string) and should return the parsed value.

## Queries

Query endpoints for SELECT operations.

```yaml
queries:
  - type: one|many
    path: /endpoint
    sql: SELECT ...
    method: GET           # optional
    accepts: json         # optional
    handle_not_found: true  # optional
    transform: { ... }    # optional
```

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `type` | `one` \| `many` | - | **Required.** `one`: single row, `many`: array |
| `path` | string | - | **Required.** Endpoint path (must start with `/`) |
| `sql` | string | - | **Required.** SQL query with `:name` placeholders |
| `method` | string | `GET` | HTTP method |
| `accepts` | string/array | `[json, form]` | Accepted Content-Types |
| `handle_not_found` | boolean | `false` | Pass `null` to post instead of 404 (type: one only) |
| `transform` | object | - | Pre/mock/post transforms |

### Type Behavior

| Type | Found | Not Found |
|------|-------|-----------|
| `one` | Returns object | Returns 404 (or `null` with `handle_not_found`) |
| `many` | Returns array | Returns empty array `[]` |

## Mutations

Mutation endpoints for INSERT/UPDATE/DELETE operations.

```yaml
mutations:
  - type: one|many|none
    method: POST
    path: /endpoint
    sql: INSERT ... RETURNING *
    accepts: json
    transform: { ... }
```

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `type` | `one` \| `many` \| `none` | - | **Required.** Return type |
| `method` | string | `POST` | HTTP method |
| `path` | string | - | **Required.** Endpoint path |
| `sql` | string | - | **Required.** SQL (use `RETURNING *` for one/many) |
| `accepts` | string/array | `[json, form]` | Accepted Content-Types |
| `transform` | object | - | Pre/mock/post transforms |

### Type Behavior

| Type | Response |
|------|----------|
| `one` | Returns single row object |
| `many` | Returns array of rows |
| `none` | Returns 204 No Content |

### MySQL-Specific Behavior

For MySQL, mutations use `Exec` instead of `Query` since MySQL does not support `RETURNING` clause. This makes `ctx.lastInsertId` and `ctx.rowsAffected` available in post-transform:

```yaml
mutations:
  - type: one
    path: /users
    sql: INSERT INTO users (name) VALUES (:name)
    transform:
      post: |
        return { id: ctx.lastInsertId, name: input.name }
```

| Variable | MySQL | Other Drivers |
|----------|-------|---------------|
| `ctx.lastInsertId` | Auto-increment ID from INSERT | `undefined` |
| `ctx.rowsAffected` | Number of affected rows | `undefined` |

> For PostgreSQL, SQLite, and SQL Server, use `RETURNING *` clause instead to get inserted data directly.

## Transform

JavaScript processing at different stages.

```yaml
transform:
  pre: |
    # Before SQL
  mock: { ... }
    # Replace SQL with mock
  post: |
    # After SQL
```

### Pre-Transform

Validates and transforms input before SQL execution.

**Variables:**
- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state (persists to post)
- `sql` (free variable): SQL string (modifiable)

**Returns:** Object with parameters for SQL

```yaml
pre: |
  ctx.startTime = Date.now();
  sql += ' WHERE active = true';
  return { id: parseInt(input.id) };
```

### Mock

Returns mock data without database. **Only one format can be specified.**

| Format | type: one | type: many | Description |
|--------|-----------|------------|-------------|
| `js` | Yes | Yes | Inline JavaScript |
| `json` | Yes | Yes | YAML/JSON object or array |
| `json_file` | Yes | Yes | JSON file |
| `csv` | No | Yes | Inline CSV |
| `csv_file` | No | Yes | CSV file |
| `jsonl` | No | Yes | Inline JSON Lines |
| `jsonl_file` | No | Yes | JSON Lines file |

**Shorthand syntax:**

```yaml
# String → JavaScript
mock: |
  return { id: 1 };

# Array → JSON (type: many)
mock:
  - { id: 1, name: Alice }
  - { id: 2, name: Bob }

# Object → JSON (type: one)
mock: { id: 1, name: Alice }

# Explicit format
mock:
  json:
    - { id: 1 }
```

**Mock JS variables:**
- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state
- `sql` (free variable): SQL string (read-only)

### Post-Transform

Transforms the result after SQL/mock execution.

**Variables:**
- `input` (parameter): Original request parameters
- `output` (parameter): Query result
- `ctx` (free variable): Shared state from pre
  - `ctx.lastInsertId`: Auto-increment ID (MySQL mutations only)
  - `ctx.rowsAffected`: Number of affected rows (MySQL mutations only)

**For type: one:**

```yaml
post: |
  return { ...output, formatted: true };
```

**For type: many:**

```yaml
# Each row
post:
  each: |
    return { ...output, processed: true };

# Entire array
post:
  all: |
    return { data: output, count: output.length };

# Both (each runs first)
post:
  each: |
    return { ...output, upper: output.name.toUpperCase() };
  all: |
    return { items: output };
```

## Error Handling

Throw an object with `status` and `body`:

```yaml
pre: |
  if (!input.token) {
    throw { status: 401, body: { error: 'unauthorized' } };
  }
```

| Phase | Default Status |
|-------|----------------|
| `pre` | 400 Bad Request |
| `mock` | 500 Internal Server Error |
| `post` | 500 Internal Server Error |

## Named Placeholders

Use `:name` syntax for SQL parameters:

```yaml
sql: SELECT * FROM users WHERE id = :id AND status = :status
```

Parameters come from:
- **GET:** Query string (`?id=1&status=active`)
- **POST/PUT/PATCH/DELETE:** Request body (JSON or form-urlencoded)

## Accepts

Control which Content-Types are accepted for request body.

```yaml
accepts: json          # Only application/json
accepts: form          # Only application/x-www-form-urlencoded
accepts: [json, form]  # Both (default)
```
