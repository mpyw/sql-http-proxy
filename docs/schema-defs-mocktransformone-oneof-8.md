# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/8
```

Inline JSONL data with filter\_by (required for type: one)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 8 Type

`object` ([Details](schema-defs-mocktransformone-oneof-8.md))

# 8 Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                             |
| :----------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [jsonl](#jsonl)          | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-8-properties-jsonl.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/8/properties/jsonl")         |
| [filter\_by](#filter_by) | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-8-properties-filter_by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/8/properties/filter_by") |

## jsonl



`jsonl`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-8-properties-jsonl.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/8/properties/jsonl")

### jsonl Type

`string`

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformone-oneof-8-properties-filter_by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/8/properties/filter_by")

### filter\_by Type

`string`

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
