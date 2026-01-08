# Inline CSV Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/2
```

Inline CSV data with header row

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 2 Type

`object` ([Inline CSV](schema-defs-mock-transform-type-many-oneof-inline-csv.md))

# 2 Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                               |
| :----------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [csv](#csv)              | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-csv-properties-csv.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/2/properties/csv")             |
| [filter\_by](#filter_by) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-csv-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/2/properties/filter_by") |

## csv



`csv`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-csv-properties-csv.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/2/properties/csv")

### csv Type

`string`

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is optional

* Type: `string` ([Filter By](schema-defs-mock-transform-type-many-oneof-inline-csv-properties-filter-by.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-inline-csv-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/2/properties/filter_by")

### filter\_by Type

`string` ([Filter By](schema-defs-mock-transform-type-many-oneof-inline-csv-properties-filter-by.md))

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
