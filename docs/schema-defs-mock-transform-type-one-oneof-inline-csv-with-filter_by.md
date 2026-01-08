# Inline CSV with filter\_by Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/6
```

Inline CSV data with filter\_by (required for type: one)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 6 Type

`object` ([Inline CSV with filter\_by](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by.md))

# 6 Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                            |
| :----------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [csv](#csv)              | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by-properties-csv.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/6/properties/csv")             |
| [filter\_by](#filter_by) | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/6/properties/filter_by") |

## csv



`csv`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by-properties-csv.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/6/properties/csv")

### csv Type

`string`

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is required

* Type: `string` ([Filter By](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by-properties-filter-by.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/6/properties/filter_by")

### filter\_by Type

`string` ([Filter By](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by-properties-filter-by.md))

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
