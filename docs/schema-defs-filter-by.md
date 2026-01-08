# Filter By Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/filterBy
```

Filter data by input parameters. When set, data is filtered to match input keys.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [schema.json\*](../out/schema.json "open original schema") |

## filterBy Type

`string` ([Filter By](schema-defs-filter-by.md))

## filterBy Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
