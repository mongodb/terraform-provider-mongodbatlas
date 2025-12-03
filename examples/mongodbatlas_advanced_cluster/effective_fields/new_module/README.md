# Building Terraform Modules with Effective Fields

This example demonstrates how to build a reusable Terraform module for MongoDB Atlas clusters using `use_effective_fields`. This approach enables a single module to support both auto-scaling and non-auto-scaling configurations without `lifecycle.ignore_changes` blocks.

## Overview

The `use_effective_fields` attribute changes how the provider handles specification attributes during auto-scaling:

- **Specification attributes** (`electable_specs`, `analytics_specs`, `read_only_specs`) remain exactly as defined in your configuration
- **Effective specs** (`effective_electable_specs`, `effective_analytics_specs`, `effective_read_only_specs`) expose actual values provisioned by Atlas, including auto-scaling changes
- **No plan drift** when Atlas auto-scales your cluster

This clear separation enables modules to work seamlessly for both auto-scaling and non-auto-scaling scenarios while providing visibility into actual cluster state.

## Module Structure

```
new_module/
├── module/                    # The reusable module
│   ├── main.tf               # Cluster resource with use_effective_fields = true
│   ├── variables.tf          # Module input variables
│   ├── outputs.tf            # Exposes both configured and effective specs
│   └── versions.tf
├── with_autoscaling/         # Example with auto-scaling enabled
│   ├── main.tf
│   ├── variables.tf
│   └── versions.tf
└── without_autoscaling/      # Example without auto-scaling
    ├── main.tf
    ├── variables.tf
    └── versions.tf
```

## Key Implementation Details

### 1. Cluster Resource with use_effective_fields

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true  # Enables effective fields behavior
  replication_specs    = var.replication_specs
  tags                 = var.tags
}
```

### 2. Data Source for Reading Effective Specs

```terraform
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
```

The data source reads effective specs showing actual provisioned values, including any changes made by auto-scaling.

### 3. Module Outputs

Expose both configured and effective specifications:

```terraform
output "configured_specs" {
  description = "Specifications as defined in configuration"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      # Maps regular specs (configured values)
    }
  ]
}

output "effective_specs" {
  description = "Actual provisioned specifications including auto-scaling changes"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      # Maps effective specs (actual values)
    }
  ]
}
```

## Usage Examples

### With Auto-Scaling

```terraform
module "atlas_cluster" {
  source = "./module"

  atlas_org_id  = var.atlas_org_id
  project_name  = "my-project"
  cluster_name  = "auto-scaling-cluster"
  cluster_type  = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"

          electable_specs = {
            instance_size = "M10"  # Baseline size; actual size may differ due to auto-scaling
            node_count    = 3
          }

          auto_scaling = {
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M10"
            compute_max_instance_size  = "M30"
            disk_gb_enabled            = true
          }
        }
      ]
    }
  ]
}

# Access both configured and effective values
output "cluster_state" {
  value = {
    configured_size = module.atlas_cluster.configured_specs[0].regions[0].electable_size
    actual_size     = module.atlas_cluster.effective_specs[0].regions[0].effective_electable_size
  }
}
```

### Without Auto-Scaling

```terraform
module "atlas_cluster" {
  source = "./module"

  atlas_org_id  = var.atlas_org_id
  project_name  = "my-project"
  cluster_name  = "fixed-size-cluster"
  cluster_type  = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"

          electable_specs = {
            instance_size = "M10"  # Fixed size
            node_count    = 3
          }
          # No auto_scaling block
        }
      ]
    }
  ]
}
```

The same module works for both scenarios without modification.

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

### Test the Examples

**With auto-scaling:**
```bash
cd with_autoscaling/
terraform init
terraform plan
terraform apply
terraform output
```

**Without auto-scaling:**
```bash
cd ../without_autoscaling/
terraform init
terraform plan
terraform apply
terraform output
```

## Benefits

1. **Single module implementation**: One module supports all scenarios
2. **No lifecycle blocks**: Eliminates maintenance overhead
3. **Visibility**: Module users can observe both configured and actual cluster state
4. **Clean configuration**: Configured values represent intent; effective specs show reality
5. **Forward compatible**: Aligns with provider v3.x default behavior

## Additional Resources

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields) - Complete documentation
- [Advanced Cluster Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)
