# README

## Top-level Schemas

* [sql-http-proxy configuration](./schema.md "Configuration file schema for sql-http-proxy") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json`

## Other Schemas

### Objects

* [CSV Config](./schema-properties-csv-config.md "Global CSV parsing options") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/csv`

* [CSV File](./schema-defs-mock-transform-type-many-oneof-csv-file.md "Path to CSV file (relative to config file)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/3`

* [CSV File with filter\_by](./schema-defs-mock-transform-type-one-oneof-csv-file-with-filter_by.md "CSV file with filter_by (required for type: one)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/7`

* [Inline CSV](./schema-defs-mock-transform-type-many-oneof-inline-csv.md "Inline CSV data with header row") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/2`

* [Inline CSV with filter\_by](./schema-defs-mock-transform-type-one-oneof-inline-csv-with-filter_by.md "Inline CSV data with filter_by (required for type: one)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/6`

* [Inline JSON](./schema-defs-mock-transform-type-one-oneof-inline-json.md "Inline JSON data (object or JSON string for type: one without filter_by)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/2`

* [Inline JSON](./schema-defs-mock-transform-type-many-oneof-inline-json.md "Inline JSON data (array or JSON string)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4`

* [Inline JSON with filter\_by](./schema-defs-mock-transform-type-one-oneof-inline-json-with-filter_by.md "Inline JSON array with filter_by (filters to single result)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3`

* [Inline JSONL](./schema-defs-mock-transform-type-many-oneof-inline-jsonl.md "Inline JSONL data (one JSON object per line)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/6`

* [Inline JSONL with filter\_by](./schema-defs-mock-transform-type-one-oneof-inline-jsonl-with-filter_by.md "Inline JSONL data with filter_by (required for type: one)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/8`

* [JSON File](./schema-defs-mock-transform-type-one-oneof-json-file.md "Path to JSON file (for type: one without filter_by)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/4`

* [JSON File](./schema-defs-mock-transform-type-many-oneof-json-file.md "Path to JSON file (relative to config file)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/5`

* [JSON File with filter\_by](./schema-defs-mock-transform-type-one-oneof-json-file-with-filter_by.md "Path to JSON file with filter_by (filters to single result)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/5`

* [JSONL File](./schema-defs-mock-transform-type-many-oneof-jsonl-file.md "Path to JSONL file (relative to config file)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/7`

* [JSONL File with filter\_by](./schema-defs-mock-transform-type-one-oneof-jsonl-file-with-filter_by.md "JSONL file with filter_by (required for type: one)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/9`

* [JavaScript Code (object)](./schema-defs-mock-transform-type-one-oneof-javascript-code-object.md "JavaScript code") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/1`

* [JavaScript Code (object)](./schema-defs-mock-transform-type-many-oneof-javascript-code-object.md "JavaScript code") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/1`

* [JavaScript Files (object)](./schema-properties-global-helpers-oneof-javascript-files-object.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1`

* [Mutation (type: many)](./schema-defs-mutation-type-many.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationMany`

* [Mutation (type: none)](./schema-defs-mutation-type-none.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationNone`

* [Mutation (type: one)](./schema-defs-mutation-type-one.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mutationOne`

* [Plain Object Shorthand](./schema-defs-mock-transform-type-one-oneof-plain-object-shorthand.md "Plain object shorthand for JSON data (e") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/10`

* [Plain Object Shorthand](./schema-defs-mock-transform-type-many-oneof-plain-object-shorthand.md "Plain object shorthand for JSON data (e") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/9`

* [Post-Transform (each/all)](./schema-defs-transform-type-many-properties-post-transform-oneof-post-transform-eachall.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany/properties/post/oneOf/1`

* [Query (type: many)](./schema-defs-query-type-many.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryMany`

* [Query (type: one)](./schema-defs-query-type-one.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/queryOne`

* [Transform (type: many)](./schema-defs-transform-type-many.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformMany`

* [Transform (type: none)](./schema-defs-transform-type-none.md "Transform for mutations with no return value (pre only, no post)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformNone`

* [Transform (type: one)](./schema-defs-transform-type-one.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/transformOne`

* [Untitled object in sql-http-proxy configuration](./schema-defs-mock-transform-type-one-oneof-inline-json-properties-json-oneof-0.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/2/properties/json/oneOf/0`

* [Untitled object in sql-http-proxy configuration](./schema-defs-mock-transform-type-one-oneof-inline-json-with-filter_by-properties-json-oneof-0-items.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3/properties/json/oneOf/0/items`

* [Untitled object in sql-http-proxy configuration](./schema-defs-mock-transform-type-many-oneof-inline-json-properties-json-oneof-0-items.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4/properties/json/oneOf/0/items`

* [Untitled object in sql-http-proxy configuration](./schema-defs-mock-transform-type-many-oneof-array-shorthand-items.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/8/items`

### Arrays

* [Array Shorthand](./schema-defs-mock-transform-type-many-oneof-array-shorthand.md "Array shorthand for JSON data (e") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/8`

* [JavaScript File Paths](./schema-properties-global-helpers-oneof-javascript-files-object-properties-javascript-file-paths.md "Paths to JavaScript helper files (relative to config)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/global_helpers/oneOf/1/properties/js_files`

* [Mutations](./schema-properties-mutations.md "List of mutation endpoints (INSERT/UPDATE/DELETE)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/mutations`

* [Queries](./schema-properties-queries.md "List of query endpoints (SELECT)") – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/properties/queries`

* [Untitled array in sql-http-proxy configuration](./schema-defs-accepted-content-types-oneof-1.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/accepts/oneOf/1`

* [Untitled array in sql-http-proxy configuration](./schema-defs-mock-transform-type-one-oneof-inline-json-with-filter_by-properties-json-oneof-0.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformOne/oneOf/3/properties/json/oneOf/0`

* [Untitled array in sql-http-proxy configuration](./schema-defs-mock-transform-type-many-oneof-inline-json-properties-json-oneof-0.md) – `https://github.com/mpyw/sql-http-proxy/internal/config/schema.json#/$defs/mockTransformMany/oneOf/4/properties/json/oneOf/0`

## Version Note

The schemas linked above follow the JSON Schema Spec version: `https://json-schema.org/draft/2020-12/schema`
