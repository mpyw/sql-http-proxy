# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## transformMany Type

`object` ([Details](schema-defs-transformmany.md))

# transformMany Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                                             |
| :------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [pre](#pre)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-pretransform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/pre")                   |
| [mock](#mock) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformmany.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/mock")             |
| [post](#post) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transformmany-properties-post.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post") |

## pre

JavaScript to transform parameters. Free vars: ctx, sql. Param: input

`pre`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-pretransform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/pre")

### pre Type

`string`

## mock

Mock data source for type: many (all formats supported)

`mock`

* is optional

* Type: merged type ([Details](schema-defs-mocktransformmany.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformmany.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/mock")

### mock Type

merged type ([Details](schema-defs-mocktransformmany.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-0.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-1.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-2.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-3.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-4.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-5.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-6.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-7.md "check type definition")

* [Untitled array in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-8.md "check type definition")

* not

  * any of

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-0.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-1.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-2.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-3.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-4.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-5.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-6.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformmany-oneof-9-not-anyof-7.md "check type definition")

## post

JavaScript to transform results

`post`

* is optional

* Type: merged type ([Details](schema-defs-transformmany-properties-post.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transformmany-properties-post.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post")

### post Type

merged type ([Details](schema-defs-transformmany-properties-post.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-posttransformall.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-transformmany-properties-post-oneof-1.md "check type definition")
