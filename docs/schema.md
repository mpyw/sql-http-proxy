# sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json
```

Configuration file schema for sql-http-proxy

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json](../out/schema.json "open original schema") |

## sql-http-proxy configuration Type

`object` ([sql-http-proxy configuration](schema.md))

any of

* [Untitled undefined type in sql-http-proxy configuration](schema-anyof-0.md "check type definition")

* [Untitled undefined type in sql-http-proxy configuration](schema-anyof-1.md "check type definition")

# sql-http-proxy configuration Properties

| Property                           | Type     | Required | Nullable       | Defined by                                                                                                                                                          |
| :--------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [global\_helpers](#global_helpers) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-global-helpers.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers") |
| [csv](#csv)                        | `object` | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-csv-config.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/csv")                |
| [dsn](#dsn)                        | `string` | Optional | can be null    | [sql-http-proxy configuration](schema-properties-dsn.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/dsn")                       |
| [queries](#queries)                | `array`  | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-queries.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/queries")               |
| [mutations](#mutations)            | `array`  | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-mutations.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/mutations")           |

## global\_helpers

Global JavaScript helpers available in all JS contexts (pre, mock, post, csv.value\_parser)

`global_helpers`

* is optional

* Type: merged type ([Global Helpers](schema-properties-global-helpers.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-global-helpers.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers")

### global\_helpers Type

merged type ([Global Helpers](schema-properties-global-helpers.md))

one (and only one) of

* [Inline JavaScript (string)](schema-properties-global-helpers-oneof-inline-javascript-string.md "check type definition")

* any of

  * [Untitled undefined type in sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-anyof-0.md "check type definition")

  * [Untitled undefined type in sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-anyof-1.md "check type definition")

## csv

Global CSV parsing options

`csv`

* is optional

* Type: `object` ([CSV Config](schema-properties-csv-config.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-csv-config.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/csv")

### csv Type

`object` ([CSV Config](schema-properties-csv-config.md))

## dsn

Database connection string. Supports ${VAR}, $VAR, ${VAR:-default} environment variable expansion. Required unless all operations use mock

`dsn`

* is optional

* Type: `string` ([DSN](schema-properties-dsn.md))

* can be null

* defined in: [sql-http-proxy configuration](schema-properties-dsn.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/dsn")

### dsn Type

`string` ([DSN](schema-properties-dsn.md))

## queries

List of query endpoints (SELECT)

`queries`

* is optional

* Type: an array of merged types ([Query](schema-properties-queries-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-queries.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/queries")

### queries Type

an array of merged types ([Query](schema-properties-queries-query.md))

### queries Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

## mutations

List of mutation endpoints (INSERT/UPDATE/DELETE)

`mutations`

* is optional

* Type: an array of merged types ([Mutation](schema-properties-mutations-mutation.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-mutations.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/mutations")

### mutations Type

an array of merged types ([Mutation](schema-properties-mutations-mutation.md))

### mutations Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

# sql-http-proxy configuration Definitions

## Definitions group path

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/path"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group sql

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/sql"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group method

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/method"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group queryMethod

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMethod"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group accepts

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/accepts"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group preTransform

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/preTransform"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group filterBy

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/filterBy"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group mockTransformOne

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group mockTransformMany

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group postTransformAll

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/postTransformAll"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group queryOne

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne"}
```

| Property                                | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                 |
| :-------------------------------------- | :------------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type)                           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-one-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/type")                         |
| [method](#method)                       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-one-properties-http-method-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/method")          |
| [path](#path)                           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-one-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/path")                |
| [sql](#sql)                             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-one-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/sql")                     |
| [accepts](#accepts)                     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/accepts")                              |
| [handle\_not\_found](#handle_not_found) | `boolean`     | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-one-properties-handle_not_found.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/handle_not_found") |
| [transform](#transform)                 | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/transform")                                |

### type

Single row query (404 if not found)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-one-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/type")

#### type Type

unknown

#### type Constraints

**constant**: the value of this property must be equal to:

```json
"one"
```

### method

HTTP method for query (default: GET)

`method`

* is optional

* Type: `string` ([HTTP Method (Query)](schema-defs-query-type-one-properties-http-method-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-one-properties-http-method-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/method")

#### method Type

`string` ([HTTP Method (Query)](schema-defs-query-type-one-properties-http-method-query.md))

#### method Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
| `"GET"`    |             |
| `"POST"`   |             |
| `"PUT"`    |             |
| `"PATCH"`  |             |
| `"DELETE"` |             |

### path

HTTP endpoint path (must start with /)

`path`

* is required

* Type: `string` ([Endpoint Path](schema-defs-query-type-one-properties-endpoint-path.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-one-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/path")

#### path Type

`string` ([Endpoint Path](schema-defs-query-type-one-properties-endpoint-path.md))

#### path Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")

### sql

SQL query with named placeholders (:name)

`sql`

* is required

* Type: `string` ([SQL Query](schema-defs-query-type-one-properties-sql-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-one-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/sql")

#### sql Type

`string` ([SQL Query](schema-defs-query-type-one-properties-sql-query.md))

#### sql Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### accepts

Accepted Content-Types for request body (default: \[json, form])

`accepts`

* is optional

* Type: merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/accepts")

#### accepts Type

merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-1.md "check type definition")

### handle\_not\_found

Pass null to post-transform instead of returning 404

`handle_not_found`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-one-properties-handle_not_found.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/handle_not_found")

#### handle\_not\_found Type

`boolean`

### transform



`transform`

* is optional

* Type: `object` ([Transform (type: one)](schema-defs-transform-type-one.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne/properties/transform")

#### transform Type

`object` ([Transform (type: one)](schema-defs-transform-type-one.md))

## Definitions group queryMany

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany"}
```

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                          |
| :------------------------ | :------------ | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type-1)           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-many-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/type")                |
| [method](#method-1)       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-many-properties-http-method-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/method") |
| [path](#path-1)           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-many-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/path")       |
| [sql](#sql-1)             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-query-type-many-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/sql")            |
| [accepts](#accepts-1)     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/accepts")                      |
| [transform](#transform-1) | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/transform")                       |

### type

Multiple rows query (returns array)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-many-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/type")

#### type Type

unknown

#### type Constraints

**constant**: the value of this property must be equal to:

```json
"many"
```

### method

HTTP method for query (default: GET)

`method`

* is optional

* Type: `string` ([HTTP Method (Query)](schema-defs-query-type-many-properties-http-method-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-many-properties-http-method-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/method")

#### method Type

`string` ([HTTP Method (Query)](schema-defs-query-type-many-properties-http-method-query.md))

#### method Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
| `"GET"`    |             |
| `"POST"`   |             |
| `"PUT"`    |             |
| `"PATCH"`  |             |
| `"DELETE"` |             |

### path

HTTP endpoint path (must start with /)

`path`

* is required

* Type: `string` ([Endpoint Path](schema-defs-query-type-many-properties-endpoint-path.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-many-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/path")

#### path Type

`string` ([Endpoint Path](schema-defs-query-type-many-properties-endpoint-path.md))

#### path Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")

### sql

SQL query with named placeholders (:name)

`sql`

* is required

* Type: `string` ([SQL Query](schema-defs-query-type-many-properties-sql-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-query-type-many-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/sql")

#### sql Type

`string` ([SQL Query](schema-defs-query-type-many-properties-sql-query.md))

#### sql Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### accepts

Accepted Content-Types for request body (default: \[json, form])

`accepts`

* is optional

* Type: merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/accepts")

#### accepts Type

merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-1.md "check type definition")

### transform



`transform`

* is optional

* Type: `object` ([Transform (type: many)](schema-defs-transform-type-many.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany/properties/transform")

#### transform Type

`object` ([Transform (type: many)](schema-defs-transform-type-many.md))

## Definitions group mutationOne

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne"}
```

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                 |
| :------------------------ | :------------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type-2)           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/type")                   |
| [method](#method-2)       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/method") |
| [path](#path-2)           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/path")          |
| [sql](#sql-2)             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/sql")               |
| [accepts](#accepts-2)     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/accepts")                           |
| [transform](#transform-2) | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/transform")                             |

### type

Mutation returning single row (via RETURNING)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/type")

#### type Type

unknown

#### type Constraints

**constant**: the value of this property must be equal to:

```json
"one"
```

### method

HTTP method for mutation (default: POST)

`method`

* is optional

* Type: `string` ([HTTP Method (Mutation)](schema-defs-mutation-type-one-properties-http-method-mutation.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/method")

#### method Type

`string` ([HTTP Method (Mutation)](schema-defs-mutation-type-one-properties-http-method-mutation.md))

#### method Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
| `"POST"`   |             |
| `"PUT"`    |             |
| `"PATCH"`  |             |
| `"DELETE"` |             |

### path

HTTP endpoint path (must start with /)

`path`

* is required

* Type: `string` ([Endpoint Path](schema-defs-mutation-type-one-properties-endpoint-path.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/path")

#### path Type

`string` ([Endpoint Path](schema-defs-mutation-type-one-properties-endpoint-path.md))

#### path Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")

### sql

SQL query with named placeholders (:name)

`sql`

* is required

* Type: `string` ([SQL Query](schema-defs-mutation-type-one-properties-sql-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-one-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/sql")

#### sql Type

`string` ([SQL Query](schema-defs-mutation-type-one-properties-sql-query.md))

#### sql Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### accepts

Accepted Content-Types for request body (default: \[json, form])

`accepts`

* is optional

* Type: merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/accepts")

#### accepts Type

merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-1.md "check type definition")

### transform



`transform`

* is optional

* Type: `object` ([Transform (type: one)](schema-defs-transform-type-one.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne/properties/transform")

#### transform Type

`object` ([Transform (type: one)](schema-defs-transform-type-one.md))

## Definitions group mutationMany

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany"}
```

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                   |
| :------------------------ | :------------ | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type-3)           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/type")                   |
| [method](#method-3)       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/method") |
| [path](#path-3)           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/path")          |
| [sql](#sql-3)             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/sql")               |
| [accepts](#accepts-3)     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/accepts")                            |
| [transform](#transform-3) | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/transform")                             |

### type

Mutation returning multiple rows (via RETURNING)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/type")

#### type Type

unknown

#### type Constraints

**constant**: the value of this property must be equal to:

```json
"many"
```

### method

HTTP method for mutation (default: POST)

`method`

* is optional

* Type: `string` ([HTTP Method (Mutation)](schema-defs-mutation-type-many-properties-http-method-mutation.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/method")

#### method Type

`string` ([HTTP Method (Mutation)](schema-defs-mutation-type-many-properties-http-method-mutation.md))

#### method Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
| `"POST"`   |             |
| `"PUT"`    |             |
| `"PATCH"`  |             |
| `"DELETE"` |             |

### path

HTTP endpoint path (must start with /)

`path`

* is required

* Type: `string` ([Endpoint Path](schema-defs-mutation-type-many-properties-endpoint-path.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/path")

#### path Type

`string` ([Endpoint Path](schema-defs-mutation-type-many-properties-endpoint-path.md))

#### path Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")

### sql

SQL query with named placeholders (:name)

`sql`

* is required

* Type: `string` ([SQL Query](schema-defs-mutation-type-many-properties-sql-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-many-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/sql")

#### sql Type

`string` ([SQL Query](schema-defs-mutation-type-many-properties-sql-query.md))

#### sql Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### accepts

Accepted Content-Types for request body (default: \[json, form])

`accepts`

* is optional

* Type: merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/accepts")

#### accepts Type

merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-1.md "check type definition")

### transform



`transform`

* is optional

* Type: `object` ([Transform (type: many)](schema-defs-transform-type-many.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany/properties/transform")

#### transform Type

`object` ([Transform (type: many)](schema-defs-transform-type-many.md))

## Definitions group mutationNone

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone"}
```

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                   |
| :------------------------ | :------------ | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [type](#type-4)           | Not specified | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/type")                   |
| [method](#method-4)       | `string`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/method") |
| [path](#path-4)           | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/path")          |
| [sql](#sql-4)             | `string`      | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/sql")               |
| [accepts](#accepts-4)     | Merged        | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/accepts")                            |
| [transform](#transform-4) | `object`      | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-none.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/transform")                             |

### type

Mutation with no return value (204 No Content)

`type`

* is required

* Type: unknown

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-type.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/type")

#### type Type

unknown

#### type Constraints

**constant**: the value of this property must be equal to:

```json
"none"
```

### method

HTTP method for mutation (default: POST)

`method`

* is optional

* Type: `string` ([HTTP Method (Mutation)](schema-defs-mutation-type-none-properties-http-method-mutation.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-http-method-mutation.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/method")

#### method Type

`string` ([HTTP Method (Mutation)](schema-defs-mutation-type-none-properties-http-method-mutation.md))

#### method Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value      | Explanation |
| :--------- | :---------- |
| `"POST"`   |             |
| `"PUT"`    |             |
| `"PATCH"`  |             |
| `"DELETE"` |             |

### path

HTTP endpoint path (must start with /)

`path`

* is required

* Type: `string` ([Endpoint Path](schema-defs-mutation-type-none-properties-endpoint-path.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-endpoint-path.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/path")

#### path Type

`string` ([Endpoint Path](schema-defs-mutation-type-none-properties-endpoint-path.md))

#### path Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")

### sql

SQL query with named placeholders (:name)

`sql`

* is required

* Type: `string` ([SQL Query](schema-defs-mutation-type-none-properties-sql-query.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mutation-type-none-properties-sql-query.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/sql")

#### sql Type

`string` ([SQL Query](schema-defs-mutation-type-none-properties-sql-query.md))

#### sql Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### accepts

Accepted Content-Types for request body (default: \[json, form])

`accepts`

* is optional

* Type: merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-accepted-content-types.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/accepts")

#### accepts Type

merged type ([Accepted Content-Types](schema-defs-accepted-content-types.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-0.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-accepted-content-types-oneof-1.md "check type definition")

### transform

Transform for mutations with no return value (pre only, no post)

`transform`

* is optional

* Type: `object` ([Transform (type: none)](schema-defs-transform-type-none.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-none.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone/properties/transform")

#### transform Type

`object` ([Transform (type: none)](schema-defs-transform-type-none.md))

## Definitions group transformOne

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne"}
```

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                                                           |
| :------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [pre](#pre)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-one-properties-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/pre")   |
| [mock](#mock) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/mock")                      |
| [post](#post) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-one-properties-post-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/post") |

### pre

JavaScript to transform parameters. Free vars: ctx, sql. Param: input

`pre`

* is optional

* Type: `string` ([Pre-Transform](schema-defs-transform-type-one-properties-pre-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-one-properties-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/pre")

#### pre Type

`string` ([Pre-Transform](schema-defs-transform-type-one-properties-pre-transform.md))

### mock

Mock data source for type: one

`mock`

* is optional

* Type: merged type ([Mock Transform (type: one)](schema-defs-mock-transform-type-one.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/mock")

#### mock Type

merged type ([Mock Transform (type: one)](schema-defs-mock-transform-type-one.md))

one (and only one) of

* [JavaScript Code (string)](schema-defs-mock-transform-type-one-oneof-javascript-code-string.md "check type definition")

* [JavaScript Code (object)](schema-defs-mock-transform-type-one-oneof-javascript-code-object.md "check type definition")

* [Inline JSON](schema-defs-mock-transform-type-one-oneof-inline-json.md "check type definition")

* [Inline JSON with filter_by](schema-defs-mock-transform-type-one-oneof-inline-json-with-filter_by.md "check type definition")

* [JSON File](schema-defs-mock-transform-type-one-oneof-json-file.md "check type definition")

* [JSON File with filter_by](schema-defs-mock-transform-type-one-oneof-json-file-with-filter_by.md "check type definition")

* [Inline CSV with filter_by](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by.md "check type definition")

* [CSV File with filter_by](schema-defs-mock-transform-type-one-oneof-csv-file-with-filter_by.md "check type definition")

* [Inline JSONL with filter_by](schema-defs-mock-transform-type-one-oneof-inline-jsonl-with-filter_by.md "check type definition")

* [JSONL File with filter_by](schema-defs-mock-transform-type-one-oneof-jsonl-file-with-filter_by.md "check type definition")

* not

  * any of

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-0.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-1.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-2.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-3.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-4.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-5.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-6.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-7.md "check type definition")

### post

JavaScript to transform result. Signature: function(ctx, input, output)

`post`

* is optional

* Type: `string` ([Post-Transform](schema-defs-transform-type-one-properties-post-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-one-properties-post-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/post")

#### post Type

`string` ([Post-Transform](schema-defs-transform-type-one-properties-post-transform.md))

## Definitions group transformMany

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany"}
```

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                                                                             |
| :-------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [pre](#pre-1)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many-properties-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/pre")   |
| [mock](#mock-1) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/mock")                      |
| [post](#post-1) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post") |

### pre

JavaScript to transform parameters. Free vars: ctx, sql. Param: input

`pre`

* is optional

* Type: `string` ([Pre-Transform](schema-defs-transform-type-many-properties-pre-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many-properties-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/pre")

#### pre Type

`string` ([Pre-Transform](schema-defs-transform-type-many-properties-pre-transform.md))

### mock

Mock data source for type: many (all formats supported)

`mock`

* is optional

* Type: merged type ([Mock Transform (type: many)](schema-defs-mock-transform-type-many.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/mock")

#### mock Type

merged type ([Mock Transform (type: many)](schema-defs-mock-transform-type-many.md))

one (and only one) of

* [JavaScript Code (string)](schema-defs-mock-transform-type-many-oneof-javascript-code-string.md "check type definition")

* [JavaScript Code (object)](schema-defs-mock-transform-type-many-oneof-javascript-code-object.md "check type definition")

* [Inline CSV](schema-defs-mock-transform-type-many-oneof-inline-csv.md "check type definition")

* [CSV File](schema-defs-mock-transform-type-many-oneof-csv-file.md "check type definition")

* [Inline JSON](schema-defs-mock-transform-type-many-oneof-inline-json.md "check type definition")

* [JSON File](schema-defs-mock-transform-type-many-oneof-json-file.md "check type definition")

* [Inline JSONL](schema-defs-mock-transform-type-many-oneof-inline-jsonl.md "check type definition")

* [JSONL File](schema-defs-mock-transform-type-many-oneof-jsonl-file.md "check type definition")

* [Array Shorthand](schema-defs-mock-transform-type-many-oneof-array-shorthand.md "check type definition")

* not

  * any of

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-0.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-1.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-2.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-3.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-4.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-5.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-6.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-7.md "check type definition")

### post

JavaScript to transform results

`post`

* is optional

* Type: merged type ([Post-Transform](schema-defs-transform-type-many-properties-post-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post")

#### post Type

merged type ([Post-Transform](schema-defs-transform-type-many-properties-post-transform.md))

one (and only one) of

* [Post-Transform (all)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-all.md "check type definition")

* [Post-Transform (each/all)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall.md "check type definition")

## Definitions group transformNone

Reference this group by using

```json
{"$ref":"https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone"}
```

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                                                                           |
| :-------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [pre](#pre-2)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-none-properties-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/pre") |
| [mock](#mock-2) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/mock")                     |

### pre

JavaScript to transform parameters. Free vars: ctx, sql. Param: input

`pre`

* is optional

* Type: `string` ([Pre-Transform](schema-defs-transform-type-none-properties-pre-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-none-properties-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/pre")

#### pre Type

`string` ([Pre-Transform](schema-defs-transform-type-none-properties-pre-transform.md))

### mock

Mock data source for type: one

`mock`

* is optional

* Type: merged type ([Mock Transform (type: one)](schema-defs-mock-transform-type-one.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/mock")

#### mock Type

merged type ([Mock Transform (type: one)](schema-defs-mock-transform-type-one.md))

one (and only one) of

* [JavaScript Code (string)](schema-defs-mock-transform-type-one-oneof-javascript-code-string.md "check type definition")

* [JavaScript Code (object)](schema-defs-mock-transform-type-one-oneof-javascript-code-object.md "check type definition")

* [Inline JSON](schema-defs-mock-transform-type-one-oneof-inline-json.md "check type definition")

* [Inline JSON with filter_by](schema-defs-mock-transform-type-one-oneof-inline-json-with-filter_by.md "check type definition")

* [JSON File](schema-defs-mock-transform-type-one-oneof-json-file.md "check type definition")

* [JSON File with filter_by](schema-defs-mock-transform-type-one-oneof-json-file-with-filter_by.md "check type definition")

* [Inline CSV with filter_by](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by.md "check type definition")

* [CSV File with filter_by](schema-defs-mock-transform-type-one-oneof-csv-file-with-filter_by.md "check type definition")

* [Inline JSONL with filter_by](schema-defs-mock-transform-type-one-oneof-inline-jsonl-with-filter_by.md "check type definition")

* [JSONL File with filter_by](schema-defs-mock-transform-type-one-oneof-jsonl-file-with-filter_by.md "check type definition")

* not

  * any of

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-0.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-1.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-2.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-3.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-4.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-5.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-6.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-7.md "check type definition")
