# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## transformOne Type

`object` ([Details](schema-defs-transformone.md))

# transformOne Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                                           |
| :------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [pre](#pre)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-pretransform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/pre")                  |
| [mock](#mock) | Merged   | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-mocktransformone.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/mock")             |
| [post](#post) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transformone-properties-post.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/post") |

## pre

JavaScript to transform parameters. Free vars: ctx, sql. Param: input

`pre`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-pretransform.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/pre")

### pre Type

`string`

## mock

Mock data source for type: one

`mock`

* is optional

* Type: merged type ([Details](schema-defs-mocktransformone.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-mocktransformone.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/mock")

### mock Type

merged type ([Details](schema-defs-mocktransformone.md))

one (and only one) of

* [Untitled string in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-0.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-1.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-2.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-3.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-4.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-5.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-6.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-7.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-8.md "check type definition")

* [Untitled object in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-9.md "check type definition")

* not

  * any of

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-0.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-1.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-2.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-3.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-4.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-5.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-6.md "check type definition")

    * [Untitled undefined type in sql-http-proxy configuration](schema-defs-mocktransformone-oneof-10-not-anyof-7.md "check type definition")

## post

JavaScript to transform result. Signature: function(ctx, input, output)

`post`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transformone-properties-post.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne/properties/post")

### post Type

`string`
