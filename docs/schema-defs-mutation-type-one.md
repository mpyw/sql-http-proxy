# Mutation (type: one) Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## mutationOne Type

`object` ([Mutation (type: one)](schema-defs-mutation-type-one.md))

# mutationOne Properties

| Property                | Type          | Required | Nullable       | Defined by                                                                                                                                                                               |
| :---------------------- | :------------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type)           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/type") |
| [method](#method)       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/method")            |
| [path](#path)           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/path")                     |
| [sql](#sql)             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/sql")                          |
| [accepts](#accepts)     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/accepts")         |
| [transform](#transform) | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/transform")           |

## type

Mutation returning single row (via RETURNING)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/type")

### type Type

unknown

### type Constraints

**constant**: the value of this property must be equal to:

```json
"one"
```

## method

HTTP method for mutation (default: POST)

`method`

* is optional

* Type: `string` ([HTTP Method (Mutation)](schema-defs-http-method-mutation.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/method")

### method Type

`string` ([HTTP Method (Mutation)](schema-defs-http-method-mutation.md))

### method Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
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

* defined in: [sql-http-proxy configuration](schema-defs-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/path")

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

* defined in: [sql-http-proxy configuration](schema-defs-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/sql")

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

* defined in: [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/accepts")

### accepts Type

merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-1.md "check type definition")

## transform



`transform`

* is optional

* Type: `object` ([Transform (type: one)](schema-defs-transform-type-one.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/transform")

### transform Type

`object` ([Transform (type: one)](schema-defs-transform-type-one.md))
