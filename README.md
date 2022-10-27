# sql-http-proxy

**[Experimental]** sql-http-proxy is a JSON configuration-based HTTP to SQL proxy server.

## Installing

```bash
go install github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy@latest
```

## Usage

Create the following `sql-http-proxy.json`:

```json
{
  "dsn": "postgres://postgres:example@localhost:5432/postgres?sslmode=disable",
  "queries": [
    {
      "type": "many",
      "path": "/users",
      "sql": "SELECT * FROM users"
    },
    {
      "type": "one",
      "path": "/user",
      "sql": "SELECT * FROM users WHERE id = $1 LIMIT 1",
      "argc": 1
    }
  ]
}
```

Launch your server: 

```bash
sql-http-proxy serve -l :8080
```

Now it accepts HTTP requests to return query results:

```bash
curl 'http://localhost:8080/users' | jq
curl 'http://localhost:8080/user?$1=123' | jq
```
