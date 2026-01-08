# CSV File Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/3
```

Path to CSV file (relative to config file)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 3 Type

`object` ([CSV File](schema-defs-mock-transform-type-many-oneof-csv-file.md))

# 3 Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                             |
| :----------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [csv\_file](#csv_file)   | `string` | Required | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-csv-file-properties-csv_file.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/3/properties/csv_file")   |
| [filter\_by](#filter_by) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-csv-file-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/3/properties/filter_by") |

## csv\_file



`csv_file`

* is required

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-csv-file-properties-csv_file.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/3/properties/csv_file")

### csv\_file Type

`string`

## filter\_by

Filter data by input parameters. When set, data is filtered to match input keys.

`filter_by`

* is optional

* Type: `string` ([Filter By](schema-defs-mock-transform-type-many-oneof-csv-file-properties-filter-by.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-csv-file-properties-filter-by.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/3/properties/filter_by")

### filter\_by Type

`string` ([Filter By](schema-defs-mock-transform-type-many-oneof-csv-file-properties-filter-by.md))

### filter\_by Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"input"` |             |
