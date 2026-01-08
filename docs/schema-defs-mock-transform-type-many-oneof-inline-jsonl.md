# Inline JSONL Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/6
```

Inline JSONL data (one JSON object per line)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 6 Type

`object` ([Inline JSONL](schema-defs-mock-transform-type-many-oneof-inline-jsonl.md))

# 6 Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                 |
| :----------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [jsonl](#jsonl)          | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-jsonl-properties-jsonl.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/6/properties/jsonl")         |
| [filter\_by](#filter_by) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-jsonl-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/6/properties/filter_by") |

## jsonl



`jsonl`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-jsonl-properties-jsonl.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/6/properties/jsonl")

### jsonl Type

`string`

### jsonl Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is optional

* Type: `string` ([Filter By](schema-defs-mock-transform-type-many-oneof-inline-jsonl-properties-filter-by.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-jsonl-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/6/properties/filter_by")

### filter\_by Type

`string` ([Filter By](schema-defs-mock-transform-type-many-oneof-inline-jsonl-properties-filter-by.md))

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
