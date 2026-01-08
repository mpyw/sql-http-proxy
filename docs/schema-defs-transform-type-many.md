# Transform (type: many) Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## transformMany Type

`object` ([Transform (type: many)](schema-defs-transform-type-many.md))

# transformMany Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                                                             |
| :------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [pre](#pre)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/pre")                                  |
| [mock](#mock) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mock-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/mock")                      |
| [post](#post) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post") |

## pre

JavaScript to transform parameters. Free vars: ctx, sql. Param: input

`pre`

* is optional

* Type: `string` ([Pre-Transform](schema-defs-pre-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-pre-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/pre")

### pre Type

`string` ([Pre-Transform](schema-defs-pre-transform.md))

## mock

Mock data source for type: many (all formats supported)

`mock`

* is optional

* Type: merged type ([Mock Transform (type: many)](schema-defs-mock-transform-type-many.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mock-transform-type-many.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/mock")

### mock Type

merged type ([Mock Transform (type: many)](schema-defs-mock-transform-type-many.md))

one (and only one) of

* [JavaScript Code (string)](schema-defs-mock-transform-type-many-oneof-javascript-code-string.md "check type definition")

* [JavaScript Code (object)](schema-defs-mock-transform-type-many-oneof-javascript-code-object.md "check type definition")

* [Inline CSV](schema-defs-mock-transform-type-many-oneof-inline-csv.md "check type definition")

* [CSV File](schema-defs-mock-transform-type-many-oneof-csv-file.md "check type definition")

* [Inline JSON](schema-defs-mock-transform-type-many-oneof-inline-json.md "check type definition")

* [JSON File](schema-defs-mock-transform-type-many-oneof-json-file.md "check type definition")

* [Inline JSONL](schema-defs-mock-transform-type-many-oneof-inline-jsonl.md "check type definition")

* [JSONL File](schema-defs-mock-transform-type-many-oneof-jsonl-file.md "check type definition")

* [Array Shorthand](schema-defs-mock-transform-type-many-oneof-array-shorthand.md "check type definition")

* not

  * any of

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-0.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-1.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-2.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-3.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-4.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-5.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-6.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mock-transform-type-many-oneof-plain-object-shorthand-not-anyof-7.md "check type definition")

## post

JavaScript to transform results

`post`

* is optional

* Type: merged type ([Post-Transform](schema-defs-transform-type-many-properties-post-transform.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post")

### post Type

merged type ([Post-Transform](schema-defs-transform-type-many-properties-post-transform.md))

one (and only one) of

* [Post-Transform (all)](schema-defs-post-transform-all.md "check type definition")

* [Post-Transform (each/all)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall.md "check type definition")
