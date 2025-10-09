# Upgrade `mongodbatlas_advanced_cluster` v1.x → v2.0.0 (Asymmetric Sharded Cluster)

This example shows how to update an asymmetric sharded cluster when upgrading from provider v1.x to v2.0.0.

Review the complete [migration guide](`docs/guides/migrate-to-advanced-cluster-2.0.0.md`) here.

## Layout
- `v1x/`: v1.x schema using blocks (symmetric-era style; limited per-shard expressiveness)
- `v2/`: v2.0.0 schema using lists/objects with per-shard specs (asymmetric shards)

## Key changes in v2
- `replication_specs` and `region_configs` become lists of objects (instead of blocks).
- Per-shard configuration is explicit (you can vary instance size or other attributes per shard).
- `disk_size_gb` moves into the inner specs (for example, `electable_specs.disk_size_gb`).
- Some references drop `[0]` or `.0` indexing because objects are no longer lists.

## Steps
1. (Optional) Apply `v1x/` if you need to reproduce a v1.x state locally.
2. Switch to `v2/` and review the updated schema and per-shard definitions.
3. Validate and apply once configuration matches your existing state.
