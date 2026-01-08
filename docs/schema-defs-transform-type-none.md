# Transform (type: none) Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone
```

Transform for mutations with no return value (pre only, no post)

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## transformNone Type

`object` ([Transform (type: none)](schema-defs-transform-type-none.md))

# transformNone Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                                       |
| :------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [pre](#pre)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/pre")            |
| [mock](#mock) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/mock") |

## pre

JavaScript to transform parameters. Free vars: ctx, sql. Param: input

`pre`

* is optional

* Type: `string` ([Pre-Transform](schema-defs-pre-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/pre")

### pre Type

`string` ([Pre-Transform](schema-defs-pre-transform.md))

## mock

Mock data source for type: one

`mock`

* is optional

* Type: merged type ([Mock Transform (type: one)](schema-defs-mock-transform-type-one.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-one.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone/properties/mock")

### mock Type

merged type ([Mock Transform (type: one)](schema-defs-mock-transform-type-one.md))

one (and only one) of

* [JavaScript Code (string)](schema-defs-mock-transform-type-one-oneof-javascript-code-string.md "check type definition")

* [JavaScript Code (object)](schema-defs-mock-transform-type-one-oneof-javascript-code-object.md "check type definition")

* [Inline JSON](schema-defs-mock-transform-type-one-oneof-inline-json.md "check type definition")

* [Inline JSON with filter_by](schema-defs-mock-transform-type-one-oneof-inline-json-with-filter_by.md "check type definition")

* [JSON File](schema-defs-mock-transform-type-one-oneof-json-file.md "check type definition")

* [JSON File with filter_by](schema-defs-mock-transform-type-one-oneof-json-file-with-filter_by.md "check type definition")

* [Inline CSV with filter_by](schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by.md "check type definition")

* [CSV File with filter_by](schema-defs-mock-transform-type-one-oneof-csv-file-with-filter_by.md "check type definition")

* [Inline JSONL with filter_by](schema-defs-mock-transform-type-one-oneof-inline-jsonl-with-filter_by.md "check type definition")

* [JSONL File with filter_by](schema-defs-mock-transform-type-one-oneof-jsonl-file-with-filter_by.md "check type definition")

* not

  * any of

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-0.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-1.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-2.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-3.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-4.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-5.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-6.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-one-oneof-plain-object-shorthand-not-anyof-7.md "check type definition")
