# Untitled object in sql-http-proxy configuration Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 1 Type

`object` ([Details](schema-defs-transformmany-properties-post-oneof-1.md))

# 1 Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                             |
| :------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [each](#each) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transformmany-properties-post-oneof-1-properties-each.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/each") |
| [all](#all)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transformmany-properties-post-oneof-1-properties-all.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/all")   |

## each

Transform each row individually. Signature: function(ctx, input, output)

`each`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transformmany-properties-post-oneof-1-properties-each.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/each")

### each Type

`string`

## all

Transform result. Signature: function(ctx, input, output)

`all`

* is optional

* Type: `string`

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transformmany-properties-post-oneof-1-properties-all.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/all")

### all Type

`string`
