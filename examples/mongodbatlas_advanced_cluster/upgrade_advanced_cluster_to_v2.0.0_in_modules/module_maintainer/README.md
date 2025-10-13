# Module Maintainer - Advanced Cluster v1.x → v2.0.0

If you own and maintain modules, this example shows how to upgrade a module that defines `mongodbatlas_advanced_cluster` from provider v1.x (old schema) to v2.0.0 (new schema) with minimal disruption to module users.

The example contains two module versions representing the key steps of the migration:

Step | Purpose | Resource
--- | --- | ---
[Step 1](./v1) | Baseline (v1.x schema) | `mongodbatlas_advanced_cluster` (old block-style schema)
[Step 2](./v2) | Upgrade to v2.0.0 schema without changing inputs | `mongodbatlas_advanced_cluster` (new list/object schema)

Notes:
- No `moved` block is required because the resource type stays the same; only the schema changes.
- If previously using `num_shards`, it has been removed in v2.0.0; express shards as multiple `replication_specs` entries.
- If previously using `disk_size_gb` at resource root level, it has been removed in v2.0.0 and now can be specified into the inner specs (for example, `electable_specs.disk_size_gb`).

Refer to the main migration guide for details: `docs/guides/migrate-to-advanced-cluster-2.0.0.md`.
