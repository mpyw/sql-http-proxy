# JavaScript Files (object) Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 1 Type

`object` ([JavaScript Files (object)](schema-properties-global-helpers-oneof-javascript-files-object.md))

any of

* [Untitled undefined type in sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-anyof-0.md "check type definition")

* [Untitled undefined type in sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-anyof-1.md "check type definition")

# 1 Properties

| Property               | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                     |
| :--------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [js](#js)              | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-properties-inline-javascript.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js")           |
| [js\_files](#js_files) | `array`  | Optional | cannot be null | [sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-properties-javascript-file-paths.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js_files") |

## js

Inline JavaScript helper code

`js`

* is optional

* Type: `string` ([Inline JavaScript](schema-properties-global-helpers-oneof-javascript-files-object-properties-inline-javascript.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-properties-inline-javascript.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js")

### js Type

`string` ([Inline JavaScript](schema-properties-global-helpers-oneof-javascript-files-object-properties-inline-javascript.md))

### js Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## js\_files

Paths to JavaScript helper files (relative to config)

`js_files`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-properties-javascript-file-paths.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js_files")

### js\_files Type

`string[]`

### js\_files Constraints

**minimum number of items**: the minimum number of items for this array is: `1`
