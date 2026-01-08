# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 1 Type

`object` ([Details](schema-properties-global_helpers-oneof-1.md))

any of

* [Untitled undefined type in sql-http-proxy configuration](schema-properties-global_helpers-oneof-1-anyof-0.md "check type definition")

* [Untitled undefined type in sql-http-proxy configuration](schema-properties-global_helpers-oneof-1-anyof-1.md "check type definition")

# 1 Properties

| Property               | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                  |
| :--------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [js](#js)              | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-global_helpers-oneof-1-properties-js.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js")             |
| [js\_files](#js_files) | `array`  | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-global_helpers-oneof-1-properties-js_files.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js_files") |

## js

Inline JavaScript helper code

`js`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-global_helpers-oneof-1-properties-js.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js")

### js Type

`string`

### js Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## js\_files

Paths to JavaScript helper files (relative to config)

`js_files`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-global_helpers-oneof-1-properties-js_files.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js_files")

### js\_files Type

`string[]`

### js\_files Constraints

**minimum number of items**: the minimum number of items for this array is: `1`
