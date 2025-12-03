# Migrating Existing Modules to use Effective Fields

This example demonstrates how module maintainers can migrate an existing Terraform module from using `lifecycle.ignore_changes` to using the `use_effective_fields` attribute, and shows the seamless upgrade experience for module users.

## Overview

If you maintain a Terraform module for MongoDB Atlas clusters with auto-scaling, you're likely using `lifecycle.ignore_changes` blocks to prevent Terraform from detecting drift when Atlas scales your clusters. The `use_effective_fields` attribute eliminates this need while providing better visibility into your cluster's actual state.

### The Problem with lifecycle.ignore_changes

Module version 1.0 (shown in `module_v1/`) demonstrates the legacy approach:

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  # ... configuration

  lifecycle {
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs[0].instance_size,
      replication_specs[0].region_configs[0].electable_specs[0].disk_size_gb,
      # ... many more fields to ignore
    ]
  }
}
```

**Limitations:**
- Cannot be conditional based on auto-scaling configuration
- Requires listing all auto-scalable attributes explicitly
- Module users cannot see actual provisioned values
- Must maintain separate modules for auto-scaling vs non-auto-scaling scenarios

### The Solution: use_effective_fields

Module version 2.0 (shown in `module_v2/`) demonstrates the modern approach:

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  # ... configuration
  use_effective_fields = true  # Single flag replaces all ignore_changes
}

data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true  # Read actual provisioned values
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
```

**Benefits:**
- Single flag replaces entire `lifecycle.ignore_changes` block
- Works seamlessly for both auto-scaling and non-auto-scaling clusters
- Module users can access both configured and effective (actual) values
- Cleaner, more maintainable code

## Directory Structure

```
existing_module/
├── README.md                  # This file
├── module_v1/                 # Legacy module using lifecycle.ignore_changes
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── versions.tf
├── module_v2/                 # Modernized module using use_effective_fields
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

In your module's `main.tf`, add `use_effective_fields = true` and remove the `lifecycle` block:

**Before (v1):**
```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id        = mongodbatlas_project.this.id
  name              = var.cluster_name
  cluster_type      = var.cluster_type
  replication_specs = var.replication_specs
  tags              = var.tags

  lifecycle {
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs[0].instance_size,
      replication_specs[0].region_configs[0].electable_specs[0].disk_size_gb,
      replication_specs[0].region_configs[0].electable_specs[0].disk_iops,
      # ... more fields
    ]
  }
}
```

**After (v2):**
```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true  # NEW: Replaces entire lifecycle block
  replication_specs    = var.replication_specs
  tags                 = var.tags
}
```

### Step 2: Add a Data Source

Add a data source to read the effective (actual) values:

```terraform
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
```

### Step 3: Update Module Outputs

Add outputs to expose both configured and effective specifications:

```terraform
# Existing output - remains unchanged
output "configured_specs" {
  description = "Configured hardware specifications"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      # ... spec mapping using regular specs
    }
  ]
}

# NEW: Add effective specs output
output "effective_specs" {
  description = "Actual provisioned specifications including auto-scaling changes"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      zone_name = spec.zone_name
      regions = [
        for region in spec.region_configs : {
          region_name              = region.region_name
          provider_name            = region.provider_name
          effective_electable_size = region.effective_electable_specs.instance_size
          effective_electable_disk = region.effective_electable_specs.disk_size_gb
          # ... additional effective fields
        }
      ]
    }
  ]
}
```

### Step 4: Update Module Version

Increment your module's version number (e.g., from 1.0 to 2.0) and publish the new version.

## Upgrade Experience for Module Users

### The Best Part: Zero Changes Required!

Module users can upgrade by simply changing the module source or version:

**Before (using v1):**
```terraform
module "atlas_cluster" {
  source  = "your-org/atlas-cluster/module"
  version = "1.0.0"  # or: source = "../module_v1"

  # All your existing configuration...
}
```

**After (upgrading to v2):**
```terraform
module "atlas_cluster" {
  source  = "your-org/atlas-cluster/module"
  version = "2.0.0"  # or: source = "../module_v2"

  # Same configuration - no changes needed!
}
```

### What Changes for Users?

1. **No Input Changes**: All input variables remain exactly the same
2. **New Outputs Available**: Users can now access `effective_specs` to see actual provisioned values
3. **No Terraform State Migration**: The upgrade is seamless - just run `terraform apply`

### Example: Before and After

**With module v1 (see `usage_v1/`):**
```terraform
output "cluster_info" {
  value = {
    configured_specs = module.atlas_cluster.configured_specs
    # effective_specs is NOT available
  }
}
```

**With module v2 (see `usage_v2/`):**
```terraform
output "cluster_info" {
  value = {
    configured_specs = module.atlas_cluster.configured_specs
    effective_specs  = module.atlas_cluster.effective_specs  # NEW!
  }
}

# NEW: Compare configured vs actual to see auto-scaling impact
output "auto_scaling_status" {
  value = {
    auto_scaling_enabled = module.atlas_cluster.auto_scaling_enabled
    configured_size      = module.atlas_cluster.configured_specs[0].regions[0].electable_size
    actual_size          = module.atlas_cluster.effective_specs[0].regions[0].effective_electable_size
  }
}
```

## Running the Examples

### Prerequisites

- Terraform >= 1.0
- MongoDB Atlas Provider v2.0.0 or later
- MongoDB Atlas service account credentials

### Configure Credentials

```bash
export MONGODB_ATLAS_CLIENT_ID="<your-client-id>"
export MONGODB_ATLAS_CLIENT_SECRET="<your-client-secret>"
export TF_VAR_atlas_org_id="<your-org-id>"
```

### Test Module v1 (Legacy Approach)

```bash
cd usage_v1/
terraform init
terraform plan
terraform apply
terraform output cluster_info
```

### Test Module v2 (Modern Approach)

```bash
cd ../usage_v2/
terraform init
terraform plan
terraform apply
terraform output cluster_info
terraform output auto_scaling_status
```

## Key Differences Between v1 and v2

| Aspect | Module v1 (legacy) | Module v2 (modern) |
|--------|-------------------|-------------------|
| **Drift Handling** | `lifecycle.ignore_changes` blocks | `use_effective_fields = true` |
| **Code Complexity** | Must list all auto-scalable fields | Single flag |
| **Data Source** | Not needed | Reads effective values |
| **Configured Specs** | Available | Available |
| **Effective Specs** | Not available | Available via outputs |
| **Auto-scaling Support** | Yes | Yes |
| **Non-auto-scaling Support** | Yes (but needs separate module) | Yes (same module) |
| **Visibility** | Limited to configured values | Both configured and actual values |
| **Maintenance** | High (must update ignore_changes) | Low (automatic) |

## Benefits of Migration

### For Module Maintainers

1. **Simplified Code**: Remove complex `lifecycle.ignore_changes` blocks
2. **Single Module**: One implementation works for all scenarios
3. **Better Outputs**: Expose both configured and effective specifications
4. **Easier Maintenance**: No need to update ignore_changes when new auto-scalable fields are added
5. **Future-Proof**: Aligns with provider v3.x where `use_effective_fields` becomes default

### For Module Users

1. **Seamless Upgrade**: Just update module version - no code changes required
2. **Better Visibility**: See actual cluster state including auto-scaling changes
3. **Same Interface**: All inputs remain unchanged
4. **Enhanced Monitoring**: Track the difference between configured and actual specifications
5. **Flexibility**: Same module works for auto-scaling and non-auto-scaling clusters

## Important Notes

### Applying the use_effective_fields Change

When first adding `use_effective_fields = true` to an existing cluster resource:

1. **Apply it separately**: Make this the only change in your first apply
2. **Then make other changes**: After successfully applying `use_effective_fields`, you can make other cluster configuration changes

This prevents potential validation errors that can occur when combining `use_effective_fields` with other modifications.

### Backward Compatibility

- Module interface (inputs) remains the same
- Existing outputs remain unchanged
- New outputs are additive (existing outputs still work)
- No breaking changes for module users

### Version 3.x

In provider version 3.x:
- `use_effective_fields` will become the default behavior
- Modules using v2 approach are already compatible
- Modules using v1 approach will need migration

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields)
- [Advanced Cluster Resource Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)
- [New Module Example](../new_module/) - Building modules from scratch with effective fields

## Support

For questions or issues with this migration pattern:
1. Review the complete examples in `module_v1/` and `module_v2/`
2. Test with the usage examples in `usage_v1/` and `usage_v2/`
3. Consult the [provider documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs)
