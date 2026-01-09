# Configuration Schema Reference

This document describes all configuration options for sql-http-proxy. For usage examples, see the main [README](README.md) and [sql-http-proxy.example.yaml](sql-http-proxy.example.yaml).

## Table of Contents

- [Top-Level Options](#top-level-options)
- [Database Configuration](#database-configuration)
  - [Database Init](#database-init)
  - [Driver Examples](#driver-examples)
- [HTTP Configuration](#http-configuration)
  - [CORS](#cors)
- [Global Helpers](#global-helpers)
- [CSV Config](#csv-config)
- [Queries](#queries)
  - [Query Type](#query-type)
  - [Query Method](#query-method)
- [Mutations](#mutations)
  - [Mutation Type](#mutation-type)
  - [Mutation Method](#mutation-method)
  - [MySQL-Specific Behavior](#mysql-specific-behavior)
- [Mock](#mock)
  - [Mock Sources](#mock-sources)
    - [Object Sources](#object-sources-type-one-only)
    - [Array Sources](#array-sources-type-many-or-type-one-with-filter)
  - [Filter](#filter)
  - [Delay](#delay)
  - [Mock JS Variables](#mock-js-variables)
- [Transform](#transform)
  - [Pre-Transform](#pre-transform)
  - [Post-Transform](#post-transform)
- [Error Handling](#error-handling)
- [Path Parameters](#path-parameters)
- [Named Placeholders](#named-placeholders)
- [Accepts](#accepts)
- [Request Object](#request-object)
- [Response Object](#response-object)
- [Headers API](#headers-api)

---

## Top-Level Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| [`database`](#database-configuration) | object | No* | Database connection configuration |
| [`http`](#http-configuration) | object | No | HTTP server configuration (CORS, etc.) |
| [`global_helpers`](#global-helpers) | object/string | No | JavaScript helpers for all transforms |
| [`csv`](#csv-config) | object | No | CSV parsing options |
| [`queries`](#queries) | array | No | Query endpoints (SELECT) |
| [`mutations`](#mutations) | array | No | Mutation endpoints (INSERT/UPDATE/DELETE) |

> *Required unless all endpoints use [mock](#mock)

## Database Configuration

Database connection configuration with `${VAR}`, `$VAR`, or `${VAR:-default}` environment variable expansion.

```yaml
database:
  dsn: postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST:-localhost}:${DB_PORT:-5432}/mydb
```

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `dsn` | string | Yes | Database connection string with env var expansion |
| `init` | string/object | No | SQL to execute on startup (e.g., schema creation) |

### Database Init

Execute SQL when the database connection is established. Useful for creating tables, seeding data, or SQLite in-memory databases.

```yaml
# Shorthand (inline SQL only)
database:
  dsn: "sqlite::memory:"
  init: |
    CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT);
    INSERT INTO users (id, name) VALUES (1, 'Alice');

# Full form (with sql_files)
database:
  dsn: "sqlite::memory:"
  init:
    sql: |
      CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT);
    sql_files:
      - ./migrations/001_init.sql
      - ./migrations/002_seed.sql
```

| Property | Type | Description |
|----------|------|-------------|
| `sql` | string | Inline SQL code to execute |
| `sql_files` | string[] | Paths to SQL files (relative to config file) |

> [!TIP]
> **Shorthand:** `init: |` is equivalent to `init: { sql: | }`

### Driver Examples

| Database | DSN Format |
|----------|------------|
| PostgreSQL | `postgres://user:pass@localhost:5432/db?sslmode=disable` |
| MySQL | `mysql://user:pass@tcp(localhost:3306)/db` |
| SQLite | `file:./data.db` or `sqlite:./data.db` |
| SQL Server | `sqlserver://user:pass@localhost:1433?database=db` |

## HTTP Configuration

HTTP server configuration including CORS settings.

```yaml
http:
  cors: true  # Permissive CORS (Access-Control-Allow-Origin: *)
```

Or with detailed configuration:

```yaml
http:
  cors:
    allowed_origins:
      - https://example.com
      - https://app.example.com
    allow_credentials: true
    max_age: 86400
```

### CORS

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `allowed_origins` | string[] | Yes (if object) | List of allowed origins. Use `["*"]` for all origins |
| `allow_credentials` | boolean | No | Allow credentials (cookies, auth headers). Default: `false` |
| `max_age` | integer | No | Preflight cache duration in seconds. Default: `0` |

> [!TIP]
> Use `cors: true` for permissive development mode. For production, specify `allowed_origins` explicitly.

## Global Helpers

JavaScript functions available in all [`pre`](#pre-transform), [`post`](#post-transform) transforms, [mock JS sources](#mock-js-variables), [`filter`](#filter), and [`csv.value_parser`](#csv-config).

```yaml
# Shorthand (inline JS only)
global_helpers: |
  function validate(x) { ... }

# Full form (with js_files)
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

> [!TIP]
> **Shorthand:** `global_helpers: |` is equivalent to `global_helpers: { js: | }`

## CSV Config

Custom value parsing for CSV [mock](#mock) data.

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

> [!IMPORTANT]
> Each query must have either `sql` OR [`mock`](#mock), not both.

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
| `path` | string | - | **Required.** Endpoint path. Supports [path parameters](#path-parameters) (e.g., `/users/{id}`) |
| `sql` | string | - | SQL query with [`:name` placeholders](#named-placeholders) (required if no mock) |
| `mock` | object | - | [Mock data source](#mock) (required if no sql) |
| `method` | string | `GET` | HTTP method |
| `accepts` | string/array | `[json, form]` | [Accepted Content-Types](#accepts) |
| `handle_not_found` | boolean | `false` | Pass `null` to [post](#post-transform) instead of 404 (type: one only) |
| `transform` | object | - | [Pre/post transforms](#transform) |

### Query Type

| Value | Description |
|-------|-------------|
| `one` | Returns a single row as an object. Returns 404 if not found (or `null` to post with `handle_not_found: true`) |
| `many` | Returns multiple rows as an array. Returns empty array `[]` if no rows found |

> [!TIP]
> Use `handle_not_found: true` to handle not-found cases in post-transform (e.g., return a default object instead of 404).

### Query Method

| Value | Description |
|-------|-------------|
| `GET` | (Default) Parameters from query string |
| `POST` | Parameters from request body |
| `PUT` | Parameters from request body |
| `PATCH` | Parameters from request body |
| `DELETE` | Parameters from request body |

## Mutations

Mutation endpoints for INSERT/UPDATE/DELETE operations.

> [!IMPORTANT]
> Each mutation (type: one/many) must have either `sql` OR [`mock`](#mock), not both.

> [!WARNING]
> `type: none` requires `sql` — mock is not supported.

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
| `path` | string | - | **Required.** Endpoint path. Supports [path parameters](#path-parameters) (e.g., `/users/{id}`) |
| `sql` | string | - | SQL (use `RETURNING *` for one/many) |
| `mock` | object | - | [Mock data source](#mock) (type: one/many only) |
| `accepts` | string/array | `[json, form]` | [Accepted Content-Types](#accepts) |
| `transform` | object | - | [Pre/post transforms](#transform) |


### Mutation Type

| Value | Description |
|-------|-------------|
| `one` | Returns a single row as an object (use `RETURNING *` in SQL) |
| `many` | Returns multiple rows as an array (use `RETURNING *` in SQL) |
| `none` | Returns 204 No Content with empty body. For fire-and-forget operations |

### Mutation Method

| Value | Description |
|-------|-------------|
| `POST` | (Default) Create resource |
| `PUT` | Replace resource |
| `PATCH` | Partial update |
| `DELETE` | Delete resource |

### MySQL-Specific Behavior

> [!WARNING]
> MySQL does not support `RETURNING` clause. Use `ctx.lastInsertId` and `ctx.rowsAffected` in post-transform instead.

For MySQL, mutations use `Exec` instead of `Query`. This makes `ctx.lastInsertId` and `ctx.rowsAffected` available in [post-transform](#post-transform):

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

> [!TIP]
> For PostgreSQL, SQLite, and SQL Server, use `RETURNING *` clause instead to get inserted data directly.

## Mock

Mock allows returning data without a database connection. Useful for testing, prototyping, or static data endpoints.

Mock is specified at the [query](#queries)/[mutation](#mutations) level (same level as `sql`). Use either `sql` OR `mock`, not both.

### Mock Sources

> [!IMPORTANT]
> Only one source type can be specified per mock.

#### Object Sources (type: one only)

| Property | Type | Description |
|----------|------|-------------|
| `object` | object | YAML object |
| `object_json` | string | JSON string containing object |
| `object_json_file` | string | Path to JSON file containing object |
| `object_js` | string | JavaScript returning object (see [Mock JS Variables](#mock-js-variables)) |

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
| `array_js` | string | JavaScript returning array (see [Mock JS Variables](#mock-js-variables)) |
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

The `filter` option allows filtering array data using JavaScript.

> [!IMPORTANT]
> For `type: one` with [array sources](#array-sources-type-many-or-type-one-with-filter), `filter` is **required** to select which row to return.

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

### Delay

Add artificial latency to mock responses. Useful for testing loading states and timeouts.

```yaml
mock:
  object: { id: 1, name: Alice }
  delay: 500ms
```

| Value | Description |
|-------|-------------|
| `delay` | Duration string (e.g., `100ms`, `1s`, `500µs`) |

**Supported units:** `ns`, `us`/`µs`, `ms`, `s`, `m`, `h`

```yaml
# With array source
mock:
  array:
    - { id: 1, name: Alice }
    - { id: 2, name: Bob }
  delay: 1s

# With filter
mock:
  array:
    - { id: 1, name: Alice }
    - { id: 2, name: Bob }
  filter: return row.id == parseInt(input.id)
  delay: 200ms
```

### Mock JS Variables

For `object_js` and `array_js`:

- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state
- `sql` (free variable): SQL string (read-only)
- `request` (free variable): HTTP [request object](#request-object) (read-only)

## Transform

JavaScript processing at different stages. See also [Global Helpers](#global-helpers) for reusable functions.

```yaml
transform:
  pre: |
    # Before SQL/mock
  post: |
    # After SQL/mock
```

### Pre-Transform

Validates and transforms input before SQL/[mock](#mock) execution.

**Variables:**
- `input` (parameter): Request parameters
- `ctx` (free variable): Shared state (persists to [post](#post-transform))
- `sql` (free variable): SQL string (modifiable, only meaningful when using sql)
- `request` (free variable): HTTP [request object](#request-object) (read-only)

**Returns:** Object with parameters for SQL/mock

```yaml
pre: |
  ctx.startTime = Date.now();
  if (sql) sql += ' WHERE active = true';
  return { id: parseInt(input.id) };
```

### Post-Transform

Transforms the result after SQL/[mock](#mock) execution.

**Variables:**
- `input` (parameter): Original request parameters
- `output` (parameter): Query result
- `ctx` (free variable): Shared state from [pre](#pre-transform)
  - `ctx.lastInsertId`: Auto-increment ID ([MySQL mutations](#mysql-specific-behavior) only)
  - `ctx.rowsAffected`: Number of affected rows ([MySQL mutations](#mysql-specific-behavior) only)
- `request` (free variable): HTTP [request object](#request-object) (read-only)
- `response` (free variable): HTTP [response object](#response-object) (writable)

**For type: one:**

```yaml
post: |
  return { ...output, formatted: true };
```

**For type: many:**

```yaml
# Shorthand: transform entire array
post: |
  return { data: output, count: output.length };

# Full form: each row only
post:
  each: |
    return { ...output, processed: true };

# Full form: entire array only
post:
  all: |
    return { data: output, count: output.length };

# Full form: both (each runs first, then all)
post:
  each: |
    return { ...output, upper: output.name.toUpperCase() };
  all: |
    return { items: output };
```

> [!TIP]
> **Shorthand:** `post: |` is equivalent to `post: { all: | }`

## Error Handling

Throw an object with `status` and `body` in any [transform](#transform):

```yaml
pre: |
  if (!input.token) {
    throw { status: 401, body: { error: 'unauthorized' } };
  }
```

> [!NOTE]
> If you throw without `status`, the default depends on the phase:

| Phase | Default Status |
|-------|----------------|
| [`pre`](#pre-transform) | 400 Bad Request |
| [`mock`](#mock) | 500 Internal Server Error |
| [`post`](#post-transform) | 500 Internal Server Error |

## Path Parameters

Use `{param}` syntax in paths to capture URL segments as parameters (chi router syntax).

```yaml
queries:
  - type: one
    path: /users/{id}
    sql: SELECT * FROM users WHERE id = :id

  - type: many
    path: /users/{user_id}/posts
    sql: SELECT * FROM posts WHERE user_id = :user_id
```

> [!CAUTION]
> Path parameters take precedence over query string and body parameters. If both exist, the path parameter wins.

```yaml
# Request: GET /users/42?id=999
# Result: id = "42" (path parameter wins)
- type: one
  path: /users/{id}
  mock:
    object_js: |
      return { received_id: input.id };  // "42"
```

**Multiple path parameters:**

```yaml
- type: one
  path: /users/{user_id}/posts/{post_id}
  sql: SELECT * FROM posts WHERE user_id = :user_id AND id = :post_id
```

**Regex constraints:**

Use `{param:regex}` syntax to validate path parameters with regular expressions. Requests that don't match will return 404.

```yaml
# Numeric ID only
- type: one
  path: /posts/{id:[0-9]+}
  sql: SELECT * FROM posts WHERE id = :id

# Slug format (lowercase letters, numbers, hyphens)
- type: one
  path: /articles/{slug:[a-z0-9-]+}
  sql: SELECT * FROM articles WHERE slug = :slug
```

**Shorthand patterns:**

For common patterns, use `{param:*shorthand*}` syntax:

| Shorthand | Description | Regex |
|-----------|-------------|-------|
| `*uuid*` | Any UUID (strict lowercase) | `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}` |
| `*uuid_v4*` | UUIDv4 only (version=4, variant=1) | `[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}` |
| `*uuid_v7*` | UUIDv7 only (version=7, variant=1) | `[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}` |

```yaml
# Any UUID
- type: one
  path: /users/{id:*uuid*}
  sql: SELECT * FROM users WHERE id = :id

# UUIDv4 only
- type: one
  path: /tokens/{id:*uuid_v4*}
  sql: SELECT * FROM tokens WHERE id = :id

# UUIDv7 only
- type: one
  path: /events/{id:*uuid_v7*}
  sql: SELECT * FROM events WHERE id = :id
```

> [!NOTE]
> The `*...*` syntax is unambiguous because `*` at the start of a regex is invalid.

## Named Placeholders

Use `:name` syntax for SQL parameters:

```yaml
sql: SELECT * FROM users WHERE id = :id AND status = :status
```

Parameters come from (in order of priority):
1. **[Path parameters](#path-parameters):** URL segments (e.g., `/users/{id}` → `id`)
2. **Query string:** URL parameters (`?status=active`)
3. **Request body:** JSON or form-urlencoded (see [Accepts](#accepts))

## Accepts

Control which Content-Types are accepted for request body.

| Value | Content-Type | Description |
|-------|--------------|-------------|
| `json` | `application/json` | JSON request body |
| `form` | `application/x-www-form-urlencoded` | Form data request body |

**Usage:**

```yaml
accepts: json          # Only application/json
accepts: form          # Only application/x-www-form-urlencoded
accepts: [json, form]  # Both (default)
```

> [!WARNING]
> Returns 415 Unsupported Media Type if the request Content-Type doesn't match.

## Request Object

The `request` free variable provides read-only access to HTTP request information. Available in [pre](#pre-transform), [post](#post-transform), and [mock JS](#mock-js-variables).

| Property | Type | Description |
|----------|------|-------------|
| `method` | string | HTTP method (e.g., `"GET"`, `"POST"`) |
| `url` | string | Request URL |
| `headers` | [Headers](#headers-api) | Request headers (read-only) |

```yaml
transform:
  pre: |
    // Check request method
    if (request.method !== 'POST') {
      throw { status: 405, body: { error: 'POST required' } };
    }
    // Read custom header
    const apiKey = request.headers.get('X-API-Key');
    if (!apiKey) {
      throw { status: 401, body: { error: 'API key required' } };
    }
    return input;
```

## Response Object

The `response` free variable allows modifying HTTP response properties. **Only available in [post-transform](#post-transform).**

| Property | Type | Access | Description |
|----------|------|--------|-------------|
| `status` | number | get/set | HTTP status code (100-599) |
| `statusText` | string | get/set | HTTP status text |
| `headers` | [Headers](#headers-api) | get (object is writable) | Response headers |
| `ok` | boolean | get | `true` if status is 200-299 |

```yaml
transform:
  post: |
    // Set custom response headers
    response.headers.set('X-Custom-Header', 'value');
    response.headers.set('Cache-Control', 'max-age=3600');

    // Change status code (affects HTTP response)
    response.status = 201;

    return output;
```

> [!NOTE]
> Setting `response.status` changes the HTTP status code of the response. This is different from throwing an error with `status` — the response body is still the return value of post-transform.

> [!WARNING]
> `response` is only available in post-transform. Accessing it in pre-transform or mock JS will result in `undefined`.

## Headers API

The Headers API follows the [Web API Headers](https://developer.mozilla.org/en-US/docs/Web/API/Headers) interface.

| Method | Description |
|--------|-------------|
| `get(name)` | Get header value (`null` if not found) |
| `has(name)` | Check if header exists |
| `set(name, value)` | Set header value (replaces existing) |
| `append(name, value)` | Add header value (allows multiple) |
| `delete(name)` | Remove header |
| `keys()` | Get all header names |
| `values()` | Get all header values |
| `entries()` | Get `[name, value]` pairs |
| `forEach(callback)` | Iterate with `callback(value, name)` |

> [!NOTE]
> Header names are case-insensitive. `headers.get('Content-Type')` and `headers.get('content-type')` return the same value.

> [!IMPORTANT]
> `request.headers` is read-only — `set()`, `append()`, and `delete()` are silently ignored.
> `response.headers` is writable — all methods work as expected.

```yaml
transform:
  post: |
    // Read request headers
    const userAgent = request.headers.get('User-Agent');
    const acceptLanguage = request.headers.get('Accept-Language');

    // Write response headers
    response.headers.set('Content-Language', 'en');
    response.headers.append('Set-Cookie', 'session=abc; HttpOnly');
    response.headers.append('Set-Cookie', 'user=123');

    // Iterate headers
    request.headers.forEach((value, name) => {
      console.log(name + ': ' + value);
    });

    return { ...output, userAgent };
```
