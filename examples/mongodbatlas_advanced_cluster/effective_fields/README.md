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

- **Want to see how to use these modules?** Start with [module_user](./module_user/) - shows that migration only requires changing the module source
- **Creating a new module?** Go directly to [module_effective_fields](./module_effective_fields/)
- **Migrating an existing module?** Review both [module_existing](./module_existing/) and [module_effective_fields](./module_effective_fields/) to understand the changes

## What is use_effective_fields?

When auto-scaling is enabled, Atlas automatically adjusts instance sizes and disk capacity. This creates [configuration drift](https://developer.hashicorp.com/terraform/tutorials/state/resource-drift) that requires management.

| Aspect | module_existing | module_effective_fields |
|--------|-----------------|-------------------------|
| **Drift handling** | `lifecycle.ignore_changes` | `use_effective_fields = true` |
| **Complexity** | Significant number of ignore_changes fields | Single attribute |
| **Data source** | Optional | Required |
| **Configured values** | Mixed with state values | Always match configuration |
| **Actual provisioned values** | Via data source (no use_effective_fields) | Via data source `effective_*` attributes |

### module_existing approach

Uses `mongodbatlas_advanced_cluster` resource with `lifecycle.ignore_changes` block listing all auto-scalable fields (instance_size, disk_size_gb, disk_iops) for all node types across regions and replication specs. Includes `mongodbatlas_advanced_cluster` data source to query actual provisioned values from Atlas API.

See [module_existing/main.tf](./module_existing/main.tf) for implementation.

### module_effective_fields approach

Uses `mongodbatlas_advanced_cluster` resource with `use_effective_fields = true` attribute. Requires `mongodbatlas_advanced_cluster` data source with `use_effective_fields = true` to expose both configured values (via `replication_specs`) and actual provisioned values (via `effective_electable_specs`, `effective_analytics_specs`, `effective_read_only_specs` attributes).

Note: `effective_*_specs` attributes are ONLY available on the data source, not the resource.

See [module_effective_fields/main.tf](./module_effective_fields/main.tf) for implementation.

## Migration Guide

### For Module Maintainers

1. **Update resource:** Add `use_effective_fields = true`, remove `lifecycle` block
2. **Add data source:** Required to read effective specs with `use_effective_fields = true`
3. **Update outputs:** Reference data source for replication specs (see [module_effective_fields/outputs.tf](./module_effective_fields/outputs.tf))

See detailed implementation in [module_existing](./module_existing/) and [module_effective_fields](./module_effective_fields/).

### For Module Users

Update the module source or version - no configuration changes required. See [module_user](./module_user/) for example.

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields)
- [Advanced Cluster Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)
