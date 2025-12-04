# MongoDB Atlas Advanced Cluster - Effective Fields Module Example

This example demonstrates the `use_effective_fields` attribute for MongoDB Atlas Terraform modules, eliminating `lifecycle.ignore_changes` blocks and providing clear visibility into Atlas-managed changes.

## Directory Structure

```
effective_fields/
├── module_existing/          # Module using lifecycle.ignore_changes
├── module_effective_fields/  # Module using use_effective_fields
└── module_user/             # Usage example (works with both modules)
```

## Quick Start

- **Migrating an existing module?** Review both [module_existing](./module_existing/) and [module_effective_fields](./module_effective_fields/) to understand the changes.
- **Creating a new module?** Go directly to [module_effective_fields](./module_effective_fields/).
- **How to use these modules?** See [module_user](./module_user/) - shows that migration only requires changing the module source.

## What is use_effective_fields?

When auto-scaling is enabled, Atlas automatically adjusts instance sizes and disk capacity. This creates [configuration drift](https://developer.hashicorp.com/terraform/tutorials/state/resource-drift) that requires management.

**Future direction:** In provider v3.x, `use_effective_fields = true` will become the default behavior and the flag will be removed. Migrating now is recommended to prepare for this transition.

### module_existing approach

Uses `mongodbatlas_advanced_cluster` resource with `lifecycle.ignore_changes` block listing all auto-scalable fields (instance_size, disk_size_gb, disk_iops) for all node types across regions and replication specs. When auto-scaling is enabled, Atlas may adjust all three fields regardless of which auto-scaling type is enabled (for optimal performance). Includes `mongodbatlas_advanced_cluster` data source to query actual provisioned values from Atlas API, as ignored fields are got from the Terraform state.

**Known issues with this approach:**
- `lifecycle.ignore_changes` blocks cannot be conditional based on auto-scaling configuration (Terraform limitation)
- Requires maintaining separate modules: one for auto-scaling clusters, one for non-auto-scaling clusters.
- Complex to maintain: must list `lifecycle.ignore_changes` lines for the biggest expected cluster topology (18-54+ fields).
- Complicates provider upgrades and cluster configuration migrations.
- Has caused issues in practice, for example: [Issue #3427](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/3427).

See [module_existing/main.tf](./module_existing/main.tf) for implementation.

### module_effective_fields approach

Uses `mongodbatlas_advanced_cluster` resource with `use_effective_fields = true` to eliminate lifecycle blocks and prevent plan drift. Requires `mongodbatlas_advanced_cluster` data source to access actual provisioned values.

**Key insight for migration:** The resource and data source flags are independent. For backward compatible migration, use `use_effective_fields = true` on the **resource** (to eliminate lifecycle blocks) but omit or set to `false` on the **data source** (to maintain output compatibility):
- **Data source without flag (default, backward compatible):**  `*_specs` returns actual provisioned values, matching module_existing behavior. Perfect for seamless migration.
- **Data source with `use_effective_fields = true` (recommended for new modules):** `*_specs` returns configured values, while `effective_*_specs` attributes return actual values. Provides clear separation between intent and reality.

Note: `effective_*_specs` attributes (effective_electable_specs, effective_analytics_specs, effective_read_only_specs) are always available on the data source for dedicated clusters, regardless of the flag value.

See [module_effective_fields/main.tf](./module_effective_fields/main.tf) for implementation with detailed comments explaining both options.

## Migration Guide

### For Module Maintainers

**Phase 1: Migrate with backward compatibility (recommended first step)**

1. **Update resource:** Add `use_effective_fields = true` and remove `lifecycle.ignore_changes` blocks.
2. **Add data source:** Add `mongodbatlas_advanced_cluster` data source WITHOUT `use_effective_fields` flag (defaults to false).
3. **Update outputs:** Reference data source for replication specs.
4. **Result:** Eliminates lifecycle blocks, prevents drift, maintains output compatibility.

**Phase 2: Enhanced visibility (prepares for provider v3.x)**

This breaking change prepares for provider v3.x where effective fields will be the default behavior:

1. **Update data source:** Add `use_effective_fields = true` to data source.
2. **Update outputs:** Expose both configured specs and effective specs separately, or document that clients must use `effective_*_specs` for actual values.
3. **Update documentation:** Clearly communicate the breaking change - data source now returns both client-provided specs (via `*_specs`) and actual provisioned specs (via `effective_*_specs`). Clients must switch from using normal specs (which previously returned actual values) to using `effective_*_specs` to get actual values.
4. **Result:** Clear separation between configured intent and actual provisioned values, aligned with future v3.x behavior.

**Breaking change impact:** Module users accessing `*_specs` for actual provisioned values must switch to using `effective_*_specs` attributes (effective_electable_specs, effective_analytics_specs, effective_read_only_specs).

See detailed implementation in [module_existing](./module_existing/) and [module_effective_fields](./module_effective_fields/).

### For Module Users

Update the module source or version - no configuration changes required. Outputs remain compatible during Phase 1 migration. See [module_user](./module_user/) for example.

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields)
- [Advanced Cluster Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)
