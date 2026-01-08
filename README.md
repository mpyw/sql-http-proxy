# sql-http-proxy

[![CI](https://github.com/mpyw/sql-http-proxy/actions/workflows/ci.yml/badge.svg)](https://github.com/mpyw/sql-http-proxy/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/mpyw/sql-http-proxy/graph/badge.svg)](https://codecov.io/gh/mpyw/sql-http-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/mpyw/sql-http-proxy)](https://goreportcard.com/report/github.com/mpyw/sql-http-proxy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> [!NOTE]
> This project was written by AI (Claude Code).

**[Experimental]** sql-http-proxy is a YAML configuration-based HTTP to SQL proxy server.

## Installation

```bash
go install github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest
```

## Usage

Create the following `.sql-http-proxy.yaml`:

```yaml
dsn: postgres://postgres:example@localhost:5432/postgres?sslmode=disable

queries:
  - type: many
    path: /users
    sql: SELECT * FROM users

  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id LIMIT 1
```

Launch your server:

```bash
sql-http-proxy -l :8080
```

Now it accepts HTTP requests to return query results:

```bash
curl 'http://localhost:8080/users' | jq
curl 'http://localhost:8080/user?id=123' | jq
```

## Supported Dialects

| Database   | DSN Example                                              | Build Tag    |
|------------|----------------------------------------------------------|--------------|
| PostgreSQL | `postgres://user:pass@localhost:5432/db?sslmode=disable` | `"postgres"` |
| MySQL      | `mysql://user:pass@tcp(localhost:3306)/db`               | `"mysql"`    |
| SQLite     | `file:./data.db` or `sqlite:./data.db`                   | `"sqlite"`   |
| SQL Server | `sqlserver://user:pass@localhost:1433?database=db`       | `"mssql"`    |
| (None)     | (Not required when all queries use `transform.mock`)     | `"mock"`     |

### Building with Specific Drivers

By default, all drivers are included. Use build tags to include only specific drivers:

```bash
# All drivers (default)
go build ./cmd/sql-http-proxy

# PostgreSQL only
go build -tags postgres ./cmd/sql-http-proxy

# PostgreSQL + MySQL
go build -tags postgres,mysql ./cmd/sql-http-proxy

# SQLite only
go build -tags sqlite ./cmd/sql-http-proxy

# No drivers (mock-only mode)
go build -tags mock ./cmd/sql-http-proxy
```

## Configuration

### DSN

The database connection string can be specified in two ways:

1. In the YAML config file: `dsn: postgres://...`
2. Via environment variable: `SQL_PROXY_DSN`

### Query Types

- `one`: Returns a single row (404 if not found)
- `many`: Returns an array of rows

### Mutation Types

Mutations are used for INSERT, UPDATE, DELETE operations. They accept JSON request body instead of URL query parameters.

- `one`: Returns a single row (via RETURNING clause)
- `many`: Returns an array of rows (via RETURNING clause)
- `none`: Returns 204 No Content (no RETURNING)

```yaml
mutations:
  # INSERT with RETURNING - returns created row
  - type: one
    method: POST
    path: /user
    sql: INSERT INTO users (name, email) VALUES (:name, :email) RETURNING *

  # UPDATE with RETURNING - returns updated row
  - type: one
    method: PUT
    path: /user
    sql: UPDATE users SET name = :name WHERE id = :id RETURNING *

  # DELETE without RETURNING - returns 204 No Content
  - type: none
    method: DELETE
    path: /user
    sql: DELETE FROM users WHERE id = :id

  # Bulk operation with RETURNING - returns multiple rows
  - type: many
    method: POST
    path: /users/bulk-update
    sql: UPDATE users SET status = :status WHERE dept = :dept RETURNING *
```

**HTTP Methods:** `POST` (default), `PUT`, `PATCH`, `DELETE`

**Request Body:** JSON object with named parameters

```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'
```

#### Handling Not Found (type: one)

By default, `type: one` returns 404 when no row is found. Use `handle_not_found: true` to pass `null` to post-transform instead:

```yaml
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    handle_not_found: true
    transform:
      post: |
        if (output === null) {
          return { found: false, default: "guest" };
        }
        return { found: true, ...output };
```

### Named Placeholders

Use named placeholders (`:name`) in SQL queries. Parameters are passed via URL query parameters:

```yaml
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id AND status = :status
```

```bash
curl 'http://localhost:8080/user?id=123&status=active'
```

Named placeholders are automatically converted to the appropriate format for each database driver.

### Transform (JavaScript)

Use JavaScript to transform query parameters (pre) and results (post):

```yaml
queries:
  # For "one" queries: output is a single row object
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      post: |
        return {
          id: output.id,
          fullName: output.first_name + ' ' + output.last_name
        }

  # For "many" queries: output is the entire array (default)
  - type: many
    path: /users
    sql: SELECT * FROM users LIMIT :limit
    transform:
      pre: |
        return { ...input, limit: input.limit || 10 }
      post: |
        return output.map(row => ({
          id: row.id,
          name: row.name.toUpperCase()
        }))
```

#### Function Signatures

**Pre-transform** - Transform query parameters and optionally modify SQL before execution:

```javascript
function(input) {
  // Free variables (mutable):
  //   ctx: Context object (shared with post-transform)
  //   sql: SQL query string (can be modified for dynamic queries)
  // Parameter:
  //   input: Query parameters object

  // Optionally modify SQL dynamically
  sql = sql + ' WHERE status = :status';

  // Store state for post-transform
  ctx.startTime = Date.now();

  return input // Must return the transformed parameters object
}
```

**Post-transform for `one` queries** - Transform the single result row:

```javascript
function(ctx, input, output) {
  // ctx: Context object (shared with pre-transform)
  // input: Original query parameters
  // output: Single row object (or null if handle_not_found: true)
  /* your code goes here */
  return output // Must return the transformed row object
}
```

**Post-transform for `many` queries** - Transform the entire result array (default):

```javascript
function(ctx, input, output) {
  // ctx: Context object (shared with pre-transform)
  // input: Original query parameters
  // output: Array of row objects
  /* your code goes here */
  return output // Must return the transformed array
}
```

#### Each vs All Mode (for `many` queries)

By default, `post` receives the entire result array. Use `post.each` to transform each row individually:

```yaml
queries:
  # Default: "all" mode - receives entire array
  - type: many
    path: /users
    sql: SELECT * FROM users
    transform:
      post: |
        return output.map(row => ({ ...row, processed: true }))

  # Explicit "each" mode - receives each row
  - type: many
    path: /users-each
    sql: SELECT * FROM users
    transform:
      post:
        each: |
          return { ...output, processed: true }

  # Both modes can be combined (each is applied first, then all)
  - type: many
    path: /users-combined
    sql: SELECT * FROM users
    transform:
      post:
        each: |
          return { ...output, normalized: true }
        all: |
          return { total: output.length, items: output }
```

#### Sharing State Between Pre and Post

The `ctx` object can be used to share state between pre and post transforms:

```yaml
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      pre: |
        ctx.requestTime = Date.now()
        return input
      post: |
        return { ...output, requestTime: ctx.requestTime }
```

#### Error Handling

Use `throw` with Lambda-style format `{ status, body }` to return custom error responses:

```yaml
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      pre: |
        if (!input.id) {
          throw { status: 400, body: { message: "id is required" } };
        }
        return input
      post: |
        if (output.role === "admin") {
          throw { status: 403, body: { message: "access denied" } };
        }
        return output
```

**Throw Format:**

```javascript
// Object body (JSON response)
throw { status: 400, body: { message: "error", code: "INVALID_INPUT" } };
// Response: {"message":"error","code":"INVALID_INPUT"}

// Array body (JSON response)
throw { status: 422, body: [{ field: "id", message: "required" }] };
// Response: [{"field":"id","message":"required"}]

// String body (JSON string response)
throw { status: 400, body: "simple error message" };
// Response: "simple error message"

// Primitive values (null, number, boolean)
throw { status: 204, body: null };
// Response: null
```

**Default Status Codes:**

When `status` is omitted, default status codes are used:
- Pre-transform errors: `400 Bad Request`
- Post-transform errors: `500 Internal Server Error`

```javascript
// In pre-transform: defaults to 400
throw { body: { message: "validation error" } };

// In post-transform: defaults to 500
throw { body: { message: "processing error" } };
```

**Native Error Objects:**

Native JavaScript `Error` objects (and any thrown object without `status` or `body` properties) always return `500 Internal Server Error`:

```javascript
// Always returns 500, regardless of pre/post context
throw new Error("something went wrong");

// Objects without status/body are treated as native errors (500)
throw { message: "error", code: "ERR_001" };
```

Supported JavaScript features include ES6 arrow functions, destructuring, spread operator, `Object.entries/keys/values`, and array methods.

### Mock (Database-less Mode)

Use `transform.mock` to return mock data without querying the database. This is useful for testing, development, or creating API stubs:

```yaml
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      mock: |
        return { id: parseInt(input.id), name: "Mock User", email: "mock@example.com" };
```

**Mock Function Signature:**

```javascript
function(input) {
  // Free variables (mutable):
  //   ctx: Context object (shared with post-transform)
  //   sql: SQL query string (read-only for reference)
  // Parameter:
  //   input: Query parameters object

  return data // Must return mock data (object for "one", array for "many")
}
```

**Mock for `one` queries** - Return an object or `null`:

```yaml
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      mock: |
        if (input.id === "999") {
          return null;  // Returns 404 (or passes null to post if handle_not_found: true)
        }
        return { id: parseInt(input.id), name: "User " + input.id };
```

**Mock for `many` queries** - Return an array:

```yaml
queries:
  - type: many
    path: /users
    sql: SELECT * FROM users
    transform:
      mock: |
        return [
          { id: 1, name: "User 1" },
          { id: 2, name: "User 2" }
        ];
```

**Mock with Pre/Post Transform:**

Mock can be combined with pre and post transforms. The `ctx` object is shared across all transforms:

```yaml
queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      pre: |
        ctx.requestedId = input.id;
        return input;
      mock: |
        ctx.mockTime = Date.now();
        return { id: parseInt(input.id), name: "Mock User" };
      post: |
        return { ...output, requestedId: ctx.requestedId, timestamp: ctx.mockTime };
```

**Running Without Database:**

When all queries use mock, the database connection is not required. You can omit `dsn` or set it to `null`:

```yaml
# No dsn needed when all queries use mock
dsn: null  # or simply omit this line

queries:
  - type: one
    path: /user
    sql: SELECT * FROM users WHERE id = :id
    transform:
      mock: |
        return { id: parseInt(input.id), name: "Mock User" };
```

The server will log "All endpoints use mock - skipping database connection" and run without connecting to any database.

## Full Example

See [examples/full-example.yaml](examples/full-example.yaml) for a comprehensive example showcasing all features.
