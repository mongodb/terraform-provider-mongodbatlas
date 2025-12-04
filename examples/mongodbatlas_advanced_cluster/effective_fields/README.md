# MongoDB Atlas Advanced Cluster - Effective Fields Module Example

This example demonstrates the `use_effective_fields` attribute for MongoDB Atlas Terraform modules, eliminating `lifecycle.ignore_changes` blocks and providing clear visibility into Atlas-managed changes.

## Quick Start

- **Creating a new module?** Go directly to [module_effective_fields](./module_effective_fields/)
- **Migrating an existing module?** Review both [module_existing](./module_existing/) and [module_effective_fields](./module_effective_fields/) to understand the changes

## What is use_effective_fields?

When auto-scaling is enabled, Atlas automatically adjusts instance sizes and disk capacity. This creates [configuration drift](https://developer.hashicorp.com/terraform/tutorials/state/resource-drift) that requires management.

### Without use_effective_fields (module_existing)

**Requires lifecycle.ignore_changes:**
```terraform
resource "mongodbatlas_advanced_cluster" "cluster" {
  # ... configuration ...

  lifecycle {
    ignore_changes = [
      # Must list all auto-scalable fields for all regions and node types
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      replication_specs[0].region_configs[0].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[0].electable_specs.disk_iops,
      # ... 50+ more fields for multi-region clusters
    ]
  }
}
```

**Limitations:**
- Complex: must list all combinations of regions, shards, and node types
- Inflexible: cannot be conditional
- Poor visibility: cannot distinguish configured values from auto-scaled values

### With use_effective_fields (module_effective_fields)

**Single attribute replaces entire lifecycle block:**
```terraform
resource "mongodbatlas_advanced_cluster" "cluster" {
  use_effective_fields = true
  # ... rest of configuration stays the same ...
}

# Data source exposes actual provisioned values
data "mongodbatlas_advanced_cluster" "cluster" {
  project_id           = mongodbatlas_advanced_cluster.cluster.project_id
  name                 = mongodbatlas_advanced_cluster.cluster.name
  use_effective_fields = true
}
```

**Benefits:**
- Simple: single attribute
- Clear: spec attributes stay constant, effective_* attributes show actual values
- No plan drift: Atlas changes don't trigger updates
- Works for both auto-scaling and non-auto-scaling scenarios

## Module Comparison

| Aspect | module_existing | module_effective_fields |
|--------|-----------------|-------------------------|
| **Drift handling** | `lifecycle.ignore_changes` | `use_effective_fields = true` |
| **Complexity** | 18-54+ ignore_changes fields | Single attribute |
| **Data source** | Not needed | Required for effective values |
| **Spec values** | Return state (may be auto-scaled) | Stay constant (match configuration) |
| **Effective values** | Not distinguishable | Available via `effective_*` attributes |
| **Visibility** | Limited | Full (configured + actual) |

## Migration Guide

### For Module Maintainers

1. **Update resource:** Add `use_effective_fields = true`, remove `lifecycle` block
2. **Add data source:** Read effective specs with `use_effective_fields = true`
3. **Update outputs:** Expose both configured and effective specs
4. **Publish:** Increment module version (e.g., 1.0 → 2.0)

See detailed comparison in [module_existing](./module_existing/) and [module_effective_fields](./module_effective_fields/).

### For Module Users

Simply update the module version:

```terraform
module "cluster" {
  source  = "your-org/cluster-module"
  version = "2.0.0"  # Only change needed

  # Same configuration - no other changes required
}
```

## Key Behavioral Change

**IMPORTANT:** How resource references work changes:

**module_existing:**
- `mongodbatlas_advanced_cluster.cluster.replication_specs` returns **state values** (may be auto-scaled)
- Cannot distinguish configured values from auto-scaled values

**module_effective_fields:**
- `data.mongodbatlas_advanced_cluster.cluster.replication_specs` returns **configuration values** (as defined in .tf files)
- `data.mongodbatlas_advanced_cluster.cluster.replication_specs[0].region_configs[0].effective_electable_specs` returns **actual provisioned values**
- Clear separation between what you configured and what Atlas provisioned

### Preserving Output Compatibility

If migrating and your outputs referenced `resource.replication_specs`:

**Option 1: Expose both (recommended)**
```terraform
output "configured_specs" {
  value = data.mongodbatlas_advanced_cluster.this.replication_specs
}

output "effective_specs" {
  # Use effective_electable_specs, effective_analytics_specs, etc.
  value = [/* map effective_* attributes */]
}
```

**Option 2: Preserve v1 behavior**
```terraform
data "mongodbatlas_advanced_cluster" "compat" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
  # Without use_effective_fields - returns state values like v1
}

output "specs" {
  value = data.mongodbatlas_advanced_cluster.compat.replication_specs
}
```

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields)
- [Advanced Cluster Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)
