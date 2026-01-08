# DSN Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/dsn
```

Database connection string. Supports ${VAR}, $VAR, ${VAR:-default} environment variable expansion. Required unless all operations use mock

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [schema.json\*](../out/schema.json "open original schema") |

## dsn Type

`string` ([DSN](schema-properties-dsn.md))
