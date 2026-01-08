# Post-Transform (each/all) Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [schema.json\*](../out/schema.json "open original schema") |

## 1 Type

`object` ([Post-Transform (each/all)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall.md))

# 1 Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                                 |
| :------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [each](#each) | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-each.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/each") |
| [all](#all)   | `string` | Optional | cannot be null | [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-all.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/all")   |

## each

Transform each row individually. Signature: function(ctx, input, output)

`each`

* is optional

* Type: `string` ([Post-Transform (each)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-each.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-each.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/each")

### each Type

`string` ([Post-Transform (each)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-each.md))

## all

Transform result. Signature: function(ctx, input, output)

`all`

* is optional

* Type: `string` ([Post-Transform (all)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-all.md))

* cannot be null

* defined in: [sql-http-proxy configuration](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-all.md "https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1/properties/all")

### all Type

`string` ([Post-Transform (all)](schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall-properties-post-transform-all.md))
