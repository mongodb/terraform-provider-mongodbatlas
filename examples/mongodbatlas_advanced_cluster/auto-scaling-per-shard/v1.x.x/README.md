# Upgrade `mongodbatlas_advanced_cluster` v1.x → v2.0.0 (Auto-scaling per Shard)

This directory contains the deprecated v1.x schema version of the auto-scaling per shard example for migration reference.
Refer the `main.tf` in the parent directory (`../`) that shows what the corresponding configuration should look like after upgrading to provider v2.0.0+.

**For details & to learn how to migrate, review the complete [Migration Guide: Advanced Cluster (v1.x → v2.0.0)](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/migrate-to-advanced-cluster-2.0#how-to-migrate)**

## Key changes in v2.0.0
- `replication_specs` becomes a list of objects (was a block in v1.x).
- `region_configs` becomes a list of objects (was a block in v1.x).
- `auto_scaling` and `analytics_auto_scaling` become attributes (were nested blocks in v1.x).
- `electable_specs` becomes an attribute (was a nested block in v1.x).
- **If previously using** `disk_size_gb` at resource root level, it has been removed in v2.0.0 and now can be specified into the inner specs (for example, `electable_specs.disk_size_gb`).
- Some references drop `[0]` or `.0` indexing because nested objects are no longer lists.
