# Migrating Existing Modules to Effective Fields

This example demonstrates how to migrate an existing Terraform module from `lifecycle.ignore_changes` to `use_effective_fields`, showing the upgrade path for both module maintainers and module users.

## Overview

This example provides a side-by-side comparison of module implementations:

- **module_v1**: Current approach using `lifecycle.ignore_changes`
- **module_v2**: Enhanced approach using `use_effective_fields`

The migration enables module users to upgrade by simply changing the module source, with no other configuration changes required.

## Directory Structure

```
existing_module/
├── README.md                  # This file
├── module_v1/                 # v1 module with lifecycle.ignore_changes
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── versions.tf
├── module_v2/                 # Modernized module with use_effective_fields
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── versions.tf
├── usage_v1/                  # Example using module_v1
│   ├── main.tf
│   ├── variables.tf
│   └── versions.tf
└── usage_v2/                  # Example using module_v2 (after upgrade)
    ├── main.tf
    ├── variables.tf
    └── versions.tf
```

## Migration Guide for Module Maintainers

### Step 1: Update the Cluster Resource

Add `use_effective_fields = true` and remove the `lifecycle` block:

**Before (module_v1):**
```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id        = mongodbatlas_project.this.id
  name              = var.cluster_name
  cluster_type      = var.cluster_type
  replication_specs = var.replication_specs
  tags              = var.tags

  # When auto-scaling is enabled (either compute or disk), Atlas may adjust
  # all three fields (instance_size, disk_size_gb, disk_iops) regardless of
  # which auto-scaling type is enabled. Must list all combinations of:
  # - replication_specs and region_configs (based on cluster topology)
  # - All 3 attributes for each node type used:
  #   * electable_specs (required - always include)
  #   * analytics_specs (optional - only if using analytics nodes)
  #   * read_only_specs (optional - only if using read-only nodes)
  lifecycle {
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      replication_specs[0].region_configs[0].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[0].electable_specs.disk_iops,
      # ... and many more fields for all node types, regions, and shards
    ]
  }
}
```

**After (module_v2):**
```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true  # Replaces entire lifecycle block
  replication_specs    = var.replication_specs
  tags                 = var.tags
}
```

### Step 2: Add Data Source for Effective Specs

```terraform
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
```

### Step 3: Update Module Outputs

Add outputs to expose effective specifications:

```terraform
# Existing output - remains unchanged
output "configured_specs" {
  description = "Specifications as defined in configuration"
  value       = [/* configured values */]
}

# New output - exposes actual values
output "effective_specs" {
  description = "Actual provisioned specifications including auto-scaling changes"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      # Maps effective specs
    }
  ]
}
```

### Step 4: Publish New Module Version

Increment your module version (e.g., 1.0 → 2.0) and publish.

## Upgrade Experience for Module Users

Module users upgrade by simply changing the module source or version:

**Before (module_v1):**
```terraform
module "atlas_cluster" {
  source  = "your-org/atlas-cluster/module"
  version = "1.0.0"

  # Configuration...
}
```

**After (module_v2):**
```terraform
module "atlas_cluster" {
  source  = "your-org/atlas-cluster/module"
  version = "2.0.0"  # Only change required

  # Same configuration - no changes needed
}
```

### What Changes for Users

1. **Inputs**: No changes - all variables remain the same
2. **Outputs**: New `effective_specs` output available (existing outputs unchanged)
3. **State**: No migration required - upgrade is seamless

**New capability:**
```terraform
output "cluster_state" {
  value = {
    configured_size = module.atlas_cluster.configured_specs[0].regions[0].electable_size
    actual_size     = module.atlas_cluster.effective_specs[0].regions[0].effective_electable_size
  }
}
```

## Running the Examples

### Prerequisites

- Terraform >= 1.0
- MongoDB Atlas Provider >= 2.0.0
- Atlas service account credentials

### Configure Credentials

```bash
export MONGODB_ATLAS_CLIENT_ID="<your-client-id>"
export MONGODB_ATLAS_CLIENT_SECRET="<your-client-secret>"
export TF_VAR_atlas_org_id="<your-org-id>"
```

### Test Module v1

```bash
cd usage_v1/
terraform init
terraform apply
terraform output
```

### Test Module v2

```bash
cd ../usage_v2/
terraform init
terraform apply
terraform output
```

## Key Differences

| Aspect | Module v1 | Module v2 |
|--------|-----------|-----------|
| **Drift Handling** | `lifecycle.ignore_changes` | `use_effective_fields = true` |
| **Code Complexity** | Must list 3 attributes × node types used × regions × shards (18+ fields for single-region, 54+ for multi-region with all node types) | Single attribute |
| **Data Source** | Not needed | Reads effective values |
| **Configured Specs** | Available | Available |
| **Effective Specs** | Not available | Available via new output |
| **Auto-scaling Support** | Yes | Yes |
| **Non-auto-scaling Support** | Yes | Yes (same module) |
| **Visibility** | Limited to configured values | Both configured and actual values |

## Benefits of Migration

**For Module Maintainers:**
- Simplified code: single attribute replaces complex lifecycle blocks
- Single module works for all scenarios
- Easier maintenance: auto-scalable fields handled automatically
- Forward compatible with provider v3.x

**For Module Users:**
- Seamless upgrade: only module version changes
- Better visibility: access to both configured and actual cluster state
- Same interface: no breaking changes
- Enhanced monitoring capabilities

## Important Notes

### Applying use_effective_fields

When first adding `use_effective_fields = true` to an existing cluster:

1. Apply this change separately (not combined with other cluster changes)
2. After successful apply, other cluster configuration changes can be made

This prevents potential validation errors.

### Backward Compatibility

- Module interface (inputs) remains unchanged
- Existing outputs remain unchanged
- New outputs are additive
- No breaking changes for module users

### Provider Version 3.x

In provider v3.x:
- Effective fields behavior becomes the default
- Modules using v2 approach are already compatible
- No further changes required

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields) - Complete documentation
- [Advanced Cluster Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)
- [New Module Example](../new_module/) - Building modules from scratch
