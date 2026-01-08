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

JavaScript functions available in all `pre`, `post` transforms, mock JS sources, `filter`, and `csv.value_parser`.

```yaml
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

Query endpoints for SELECT operations. Each query must have either `sql` OR `mock`, not both.

```yaml
queries:
  - type: one|many
    path: /endpoint
    sql: SELECT ...         # OR mock: { ... }
    method: GET             # optional
    accepts: json           # optional
    handle_not_found: true  # optional (type: one only)
    transform: { ... }      # optional
```

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `type` | `one` \| `many` | - | **Required.** `one`: single row, `many`: array |
| `path` | string | - | **Required.** Endpoint path (must start with `/`) |
| `sql` | string | - | SQL query with `:name` placeholders (required if no mock) |
| `mock` | object | - | Mock data source (required if no sql) |
| `method` | string | `GET` | HTTP method |
| `accepts` | string/array | `[json, form]` | Accepted Content-Types |
| `handle_not_found` | boolean | `false` | Pass `null` to post instead of 404 (type: one only) |
| `transform` | object | - | Pre/post transforms |

### Type Behavior

| Type | Found | Not Found |
|------|-------|-----------|
| `one` | Returns object | Returns 404 (or `null` with `handle_not_found`) |
| `many` | Returns array | Returns empty array `[]` |

## Mutations

Mutation endpoints for INSERT/UPDATE/DELETE operations. Each mutation (type: one/many) must have either `sql` OR `mock`, not both. Type: none requires `sql`.

```yaml
mutations:
  - type: one|many|none
    method: POST
    path: /endpoint
    sql: INSERT ... RETURNING *  # OR mock: { ... } for type: one/many
    accepts: json
    transform: { ... }
```

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `type` | `one` \| `many` \| `none` | - | **Required.** Return type |
| `method` | string | `POST` | HTTP method |
| `path` | string | - | **Required.** Endpoint path |
| `sql` | string | - | SQL (use `RETURNING *` for one/many) |
| `mock` | object | - | Mock data source (type: one/many only) |
| `accepts` | string/array | `[json, form]` | Accepted Content-Types |
| `transform` | object | - | Pre/post transforms |

> Note: `type: none` requires `sql` - mock is not supported for type: none.

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

## Mock

Mock allows returning data without a database connection. Useful for testing, prototyping, or static data endpoints.

Mock is specified at the query/mutation level (same level as `sql`). Use either `sql` OR `mock`, not both.

### Mock Sources

**Only one source type can be specified per mock.**

#### Object Sources (type: one only)

| Property | Type | Description |
|----------|------|-------------|
| `object` | object | YAML object |
| `object_json` | string | JSON string containing object |
| `object_json_file` | string | Path to JSON file containing object |
| `object_js` | string | JavaScript returning object |

```yaml
# YAML object
mock:
  object: { id: 1, name: Alice }

# JSON string
mock:
  object_json: '{"id": 1, "name": "Alice"}'

# JSON file
mock:
  object_json_file: ./data/user.json

# JavaScript
mock:
  object_js: |
    return { id: parseInt(input.id), name: "User " + input.id };
```

#### Array Sources (type: many, or type: one with filter)

| Property | Type | Description |
|----------|------|-------------|
| `array` | array | YAML array of objects |
| `array_json` | string | JSON string containing array |
| `array_json_file` | string | Path to JSON file containing array |
| `array_js` | string | JavaScript returning array |
| `csv` | string | Inline CSV data with header row |
| `csv_file` | string | Path to CSV file |
| `jsonl` | string | Inline JSONL (one JSON object per line) |
| `jsonl_file` | string | Path to JSONL file |

```yaml
# YAML array
mock:
  array:
    - { id: 1, name: Alice }
    - { id: 2, name: Bob }

# JSON string
mock:
  array_json: '[{"id": 1}, {"id": 2}]'

# CSV
mock:
  csv: |
    id,name,role
    1,Alice,admin
    2,Bob,user

# JavaScript
mock:
  array_js: |
    return [{ id: 1 }, { id: 2 }];
```

### Filter

The `filter` option allows filtering array data using JavaScript. For `type: one` with array sources, `filter` is **required** to select which row to return.

**Filter variables:**
- `row` (parameter): Current row being evaluated
- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state

**Returns:** Boolean (true to include the row)

```yaml
# type: one with filter - returns first matching row (or 404)
- type: one
  path: /user
  mock:
    array:
      - { id: 1, name: Alice }
      - { id: 2, name: Bob }
    filter: return row.id == parseInt(input.id)

# type: many with filter - returns all matching rows
- type: many
  path: /users
  mock:
    array:
      - { id: 1, name: Alice, role: admin }
      - { id: 2, name: Bob, role: user }
      - { id: 3, name: Charlie, role: admin }
    filter: return row.role === input.role
```

### Mock JS Variables

For `object_js` and `array_js`:

- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state
- `sql` (free variable): SQL string (read-only)

## Transform

JavaScript processing at different stages.

```yaml
transform:
  pre: |
    # Before SQL/mock
  post: |
    # After SQL/mock
```

### Pre-Transform

Validates and transforms input before SQL/mock execution.

**Variables:**
- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state (persists to post)
- `sql` (free variable): SQL string (modifiable, only meaningful when using sql)

**Returns:** Object with parameters for SQL/mock

```yaml
pre: |
  ctx.startTime = Date.now();
  if (sql) sql += ' WHERE active = true';
  return { id: parseInt(input.id) };
```

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
