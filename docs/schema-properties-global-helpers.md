# Global Helpers Schema

```txt
https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers
```

Global JavaScript helpers available in all JS contexts (pre, mock, post, csv.value\_parser)

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                 |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [schema.json\*](../out/schema.json "open original schema") |

## global\_helpers Type

merged type ([Global Helpers](schema-properties-global-helpers.md))

one (and only one) of

* [Inline JavaScript (string)](schema-properties-global-helpers-oneof-inline-javascript-string.md "check type definition")

* any of

  * [Untitled undefined type in sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-anyof-0.md "check type definition")

  * [Untitled undefined type in sql-http-proxy configuration](schema-properties-global-helpers-oneof-javascript-files-object-anyof-1.md "check type definition")
