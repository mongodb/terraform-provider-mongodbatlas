# Upgrade `mongodbatlas_advanced_cluster` v1.x → v2.0.0 (Asymmetric Sharded Cluster)

This directory contains the v1.x schema version of the asymmetric sharded cluster example for migration reference.
Refer the `main.tf` in the parent directory (`../`) that shows what the corresponding configuration should look like for provider v2.0.0+.

**Review the complete [migration guide here](`docs/guides/migrate-to-advanced-cluster-2.0.0.md`).**

 ## Key changes in v2.0.0
 - `replication_specs` becomes a list of objects (was a block in v1.x).
 - `region_configs` becomes a list of objects (was a block in v1.x).
 - `electable_specs` becomes an attribute (was a nested block in v1.x).
 - `advanced_configuration` becomes an attribute (was a block in v1.x).
 - `tags` becomes an attribute (was a block in v1.x).
 - Per-shard configuration remains explicit and can vary (e.g., different `instance_size` per shard).
- `disk_size_gb` moves into the inner specs (for example, `electable_specs.disk_size_gb`).
- Some references drop `[0]` or `.0` indexing because objects are no longer lists but singular attributes (e.g., `replication_specs[0].region_configs[0].electable_specs.disk_size_gb`).

## Steps
1. (Optional) Use this v1.x example to reproduce a legacy state.
2. Compare with the v2 example in the parent directory and update your configuration accordingly.
3. Validate and apply once configuration matches your existing state.
