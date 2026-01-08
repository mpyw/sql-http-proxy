# CSV Config Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/csv
```

Global CSV parsing options

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## csv Type

`object` ([CSV Config](schema-properties-csv-config.md))

# csv Properties

| Property                       | Type     | Required | Nullable       | Defined by                                                                                                                                                                                           |
| :----------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [value\_parser](#value_parser) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-csv-config-properties-value-parser.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/csv/properties/value_parser") |

## value\_parser

JavaScript code to parse CSV cell values. Signature: function(value) { return parsed }

`value_parser`

* is optional

* Type: `string` ([Value Parser](schema-properties-csv-config-properties-value-parser.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-csv-config-properties-value-parser.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/csv/properties/value_parser")

### value\_parser Type

`string` ([Value Parser](schema-properties-csv-config-properties-value-parser.md))

### value\_parser Constraints

**minimum length**: the minimum number of characters for this string is: `1`
