# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3
```

Inline JSON array with filter\_by (filters to single result)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 3 Type

`object` ([Details](schema-defs-mocktransformone-oneof-3.md))

# 3 Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                             |
| :----------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [json](#json)            | Merged   | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-3-properties-json.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3/properties/json")           |
| [filter\_by](#filter_by) | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-3-properties-filter_by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3/properties/filter_by") |

## json



`json`

* is required

* Type: merged type ([Details](schema-defs-mocktransformone-oneof-3-properties-json.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-3-properties-json.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3/properties/json")

### json Type

merged type ([Details](schema-defs-mocktransformone-oneof-3-properties-json.md))

one (and only one) of

* [Untitled array in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-3-properties-json-oneof-0.md "check type definition")

* [Untitled string in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-3-properties-json-oneof-1.md "check type definition")

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-3-properties-filter_by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3/properties/filter_by")

### filter\_by Type

`string`

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
