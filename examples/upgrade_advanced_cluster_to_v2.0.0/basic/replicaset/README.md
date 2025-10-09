# Upgrade `mongodbatlas_advanced_cluster` v1.x → v2.0.0 (REPLICASET)

This example demonstrates the specific configuration changes required to update a REPLICASET `mongodbatlas_advanced_cluster` from provider v1.x.x to v2.0.0 or later.

Review the complete [migration guide](`docs/guides/migrate-to-advanced-cluster-2.0.0.md`) here.

## Layout

- `v1x/`: Old schema (block-based) configuration for v1.x
- `v2/`: New schema (list/object-based) configuration for v2.0.0

Key changes reflected here:
- Blocks like `replication_specs`, `region_configs`, `electable_specs`, `analytics_specs` become attributes (lists/objects) in v2.0.0.
- Disk size moves into the inner specs (`electable_specs.disk_size_gb`, `analytics_specs.disk_size_gb`).
- Many outputs drop `[0]` or `.0` indexing because objects are no longer lists.

## Steps

1. Apply the v1.x configuration (this is optional if you already have a v1.x state):
   - Edit `v1x/versions.tf` to match your current provider version if needed (<= 1.41.0).
   - Provide variables (see `v1x/variables.tf`).
   - Run `terraform init && terraform apply`.

2. Upgrade to MongoDB Atlas provider v2.0.0 or later and:
   - Update the `mongodbatlas_advanced_cluster` configuration with the required updates as demonstrated in `v2/main.tf`. 
    - Additionally, remove any deprecated attributes as mentioned in the [migration guide](`docs/guides/migrate-to-advanced-cluster-2.0.0.md`).
   - Update references and indexing following the guide (for example, `connection_strings.standard` instead of `[0].standard`).

4. Validate and apply:
   - Run `terraform validate` and `terraform plan` should show no changes.
   - Run `terraform apply` only if no changes planned in previous step. This will update the `mongodbatlas_advanced_cluster` state to support the new schema.
