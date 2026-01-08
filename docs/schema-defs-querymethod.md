# Untitled string in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMethod
```

HTTP method for query (default: GET)

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [schema.json\*](../out/schema.json "open original schema") |

## queryMethod Type

`string`

## queryMethod Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
| `"GET"`    |             |
| `"POST"`   |             |
| `"PUT"`    |             |
| `"PATCH"`  |             |
| `"DELETE"` |             |
