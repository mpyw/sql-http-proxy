# Inline JSON Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4
```

Inline JSON data (array or JSON string)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 4 Type

`object` ([Inline JSON](schema-defs-mock-transform-type-many-oneof-inline-json.md))

# 4 Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                |
| :----------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [json](#json)            | Merged   | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-json-properties-json.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4/properties/json")           |
| [filter\_by](#filter_by) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-json-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4/properties/filter_by") |

## json



`json`

* is required

* Type: merged type ([Details](schema-defs-mock-transform-type-many-oneof-inline-json-properties-json.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-json-properties-json.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4/properties/json")

### json Type

merged type ([Details](schema-defs-mock-transform-type-many-oneof-inline-json-properties-json.md))

one (and only one) of

* [Untitled array in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-json-properties-json-oneof-0.md "check type definition")

* [Untitled string in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-json-properties-json-oneof-1.md "check type definition")

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is optional

* Type: `string` ([Filter By](schema-defs-mock-transform-type-many-oneof-inline-json-properties-filter-by.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-json-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4/properties/filter_by")

### filter\_by Type

`string` ([Filter By](schema-defs-mock-transform-type-many-oneof-inline-json-properties-filter-by.md))

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
