# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/7
```

Path to JSONL file (relative to config file)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 7 Type

`object` ([Details](schema-defs-mocktransformmany-oneof-7.md))

# 7 Properties

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                 |
| :------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [jsonl\_file](#jsonl_file) | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-7-properties-jsonl_file.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/7/properties/jsonl_file") |
| [filter\_by](#filter_by)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-7-properties-filter_by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/7/properties/filter_by")   |

## jsonl\_file



`jsonl_file`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-7-properties-jsonl_file.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/7/properties/jsonl_file")

### jsonl\_file Type

`string`

### jsonl\_file Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-7-properties-filter_by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/7/properties/filter_by")

### filter\_by Type

`string`

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
