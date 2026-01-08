# Query (type: many) Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## queryMany Type

`object` ([Query (type: many)](schema-defs-query-type-many.md))

# queryMany Properties

| Property                | Type          | Required | Nullable       | Defined by                                                                                                                                                                           |
| :---------------------- | :------------ | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type)           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-many-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/type") |
| [method](#method)       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-http-method-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/method")             |
| [path](#path)           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/path")                   |
| [sql](#sql)             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/sql")                        |
| [accepts](#accepts)     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/accepts")       |
| [transform](#transform) | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/transform")        |

## type

Multiple rows query (returns array)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-many-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/type")

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

* Type: `string` ([HTTP Method (Query)](schema-defs-http-method-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-http-method-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/method")

### method Type

`string` ([HTTP Method (Query)](schema-defs-http-method-query.md))

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

* Type: `string` ([Endpoint Path](schema-defs-endpoint-path.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/path")

### path Type

`string` ([Endpoint Path](schema-defs-endpoint-path.md))

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

* Type: `string` ([SQL Query](schema-defs-sql-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/sql")

### sql Type

`string` ([SQL Query](schema-defs-sql-query.md))

### sql Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## accepts

Accepted Content-Types for request body (default: \[json, form])

`accepts`

* is optional

* Type: merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/accepts")

### accepts Type

merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-1.md "check type definition")

## transform



`transform`

* is optional

* Type: `object` ([Transform (type: many)](schema-defs-transform-type-many.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/transform")

### transform Type

`object` ([Transform (type: many)](schema-defs-transform-type-many.md))
