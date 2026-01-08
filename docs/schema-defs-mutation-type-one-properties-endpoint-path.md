# Endpoint Path Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/path
```

HTTP endpoint path (must start with /)

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [schema.json\*](../out/schema.json "open original schema") |

## path Type

`string` ([Endpoint Path](schema-defs-mutation-type-one-properties-endpoint-path.md))

## path Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")
