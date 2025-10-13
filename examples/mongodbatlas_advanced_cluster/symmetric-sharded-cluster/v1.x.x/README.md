# Upgrade `mongodbatlas_advanced_cluster` v1.x → v2.0.0 (Symmetric Sharded Cluster)

This directory contains the deprecated v1.x schema version of the symmetric sharded (2 shards) cluster example for migration reference.
Refer the `main.tf` in the parent directory (`../`) that shows what the corresponding configuration should look like after upgrading to provider v2.0.0+.

 **For details & to learn how to migrate, review the complete [Migration Guide: Advanced Cluster (v1.x → v2.0.0)](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/migrate-to-advanced-cluster-2.0#how-to-migrate)**

 ## Key changes in v2.0.0
 - `replication_specs` becomes a list of objects (was a block in v1.x).
 - `region_configs` becomes a list of objects (was a block in v1.x).
 - `electable_specs` becomes an attribute (was a nested block in v1.x).
 - `advanced_configuration` becomes an attribute (was a block in v1.x).
 - `tags` becomes an attribute (was a block in v1.x).
 - **If previously using** `disk_size_gb` at resource root level, it has been removed in v2.0.0 and now can be specified into the inner specs (for example, `electable_specs.disk_size_gb`).
 - **If previously using** `disk_size_gb` at resource root level, it has been removed in v2.0.0 and now can be specified into the inner specs (for example, `electable_specs.disk_size_gb`).
 - Some references drop `[0]` or `.0` indexing because objects are no longer lists but singular attributes (e.g., `replication_specs[0].region_configs[0].electable_specs.disk_size_gb`).

## Steps
1. (Optional) Use this v1.x example to reproduce a legacy state.
2. Compare with the v2 example in the parent directory and update your configuration accordingly.
3. Validate and apply once configuration matches your existing state.
