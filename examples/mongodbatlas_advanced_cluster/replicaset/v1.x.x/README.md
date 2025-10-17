# Upgrade `mongodbatlas_advanced_cluster` v1.x → v2.0.0 (Replicaset)

This directory contains the deprecated v1.x schema version of the Replica Set example for migration reference.
Refer the `main.tf` in the parent directory (`../`) that shows what the corresponding configuration should look like after upgrading to provider v2.0.0+.

**For details & to learn how to migrate, review the complete [Migration Guide: Advanced Cluster (v1.x → v2.0.0)](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/migrate-to-advanced-cluster-2.0#how-to-migrate)**

## Key changes in v2.0.0
- `replication_specs` becomes a list of objects (was a block in v1.x).
- `region_configs` becomes a list of objects (was a block in v1.x).
- `electable_specs` becomes an attribute (was a nested block in v1.x).
- `tags` becomes a map attribute (was a block in v1.x).
- **If previously using** `num_shards`, it has been removed in v2.0.0; sharded layouts use multiple `replication_specs` entries instead and REPLICASET clusters don't require this attribute.
- Some references may drop `[0]` or `.0` indexing because nested objects are no longer lists.
