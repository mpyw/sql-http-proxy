# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## queryMany Type

`object` ([Details](schema-defs-querymany.md))

# queryMany Properties

| Property                | Type          | Required | Nullable       | Defined by                                                                                                                                                                     |
| :---------------------- | :------------ | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type)           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-querymany-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/type") |
| [method](#method)       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-querymethod.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/method")             |
| [path](#path)           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/path")                      |
| [sql](#sql)             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-sql.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/sql")                        |
| [accepts](#accepts)     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepts.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/accepts")                |
| [transform](#transform) | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transformmany.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/transform")        |

## type

Multiple rows query (returns array)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-querymany-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/type")

### type Type

unknown

### type Constraints

**constant**: the value of this property must be equal to:

```json
"many"
```

## method

HTTP method for query (default: GET)

`method`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-querymethod.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/method")

### method Type

`string`

### method Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
| `"GET"`    |             |
| `"POST"`   |             |
| `"PUT"`    |             |
| `"PATCH"`  |             |
| `"DELETE"` |             |

## path

HTTP endpoint path (must start with /)

`path`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/path")

### path Type

`string`

### path Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")

## sql

SQL query with named placeholders (:name)

`sql`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-sql.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/sql")

### sql Type

`string`

### sql Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## accepts

Accepted Content-Types for request body (default: \[json, form])

`accepts`

* is optional

* Type: merged type ([Details](schema-defs-accepts.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-accepts.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/accepts")

### accepts Type

merged type ([Details](schema-defs-accepts.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepts-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepts-oneof-1.md "check type definition")

## transform



`transform`

* is optional

* Type: `object` ([Details](schema-defs-transformmany.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transformmany.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/transform")

### transform Type

`object` ([Details](schema-defs-transformmany.md))
