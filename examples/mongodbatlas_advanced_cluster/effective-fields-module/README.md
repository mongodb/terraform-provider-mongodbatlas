# MongoDB Atlas Advanced Cluster Module with Effective Fields

This example demonstrates how to create a **reusable Terraform module** for MongoDB Atlas clusters that seamlessly supports both auto-scaling and non-auto-scaling configurations without requiring `lifecycle.ignore_changes` blocks.

## Why This Matters for Module Authors

When creating reusable Terraform modules for MongoDB Atlas clusters, you typically face a challenge:

- **Without auto-scaling**: Module works fine, no special handling needed
- **With auto-scaling**: Atlas changes instance sizes dynamically, causing Terraform plan drift unless you use `lifecycle.ignore_changes`

### The Problem with Traditional Approaches

**Before `use_effective_fields`, you needed TWO different modules or complex workarounds:**

```terraform
# ❌ Problem: Module users must add lifecycle.ignore_changes themselves
module "cluster" {
  source = "./cluster-module"
  # ... configuration

  lifecycle {
    ignore_changes = [
      # Users must know which fields to ignore
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      # ... and maintain this list
    ]
  }
}
```

**Or you'd need separate modules:**
- `cluster-module-with-autoscaling/` (includes lifecycle.ignore_changes)
- `cluster-module-without-autoscaling/` (no lifecycle blocks)

### The Solution: use_effective_fields

By setting `use_effective_fields = true` in your module, you create **ONE module that works for both scenarios**:

```terraform
# ✅ Solution: Single module works for both auto-scaling and fixed-size clusters
module "cluster" {
  source = "./cluster-module"
  enable_auto_scaling = true  # or false - module handles both!
  # ... configuration
}
# No lifecycle.ignore_changes needed!
```

## Benefits for Module Authors

1. **Single Module for All Use Cases**: Write one module that handles auto-scaling and non-auto-scaling clusters
2. **No lifecycle.ignore_changes Required**: Module users don't need to add lifecycle blocks
3. **Visibility into Scaled Values**: Module can output both configured and effective (actual) instance sizes
4. **Cleaner Module Interface**: Simpler for consumers, easier to maintain
5. **Future-Proof**: Aligns with provider version 3.0 where this will be default behavior

## Module Structure

```
effective-fields-module/
├── module/                    # The reusable module
│   ├── main.tf               # Module resources with use_effective_fields = true
│   ├── variables.tf          # Module input variables
│   └── outputs.tf            # Module outputs including effective specs
├── with-autoscaling/         # Example using the module WITH auto-scaling
│   ├── main.tf
│   ├── variables.tf
│   └── versions.tf
└── without-autoscaling/      # Example using the module WITHOUT auto-scaling
    ├── main.tf
    ├── variables.tf
    └── versions.tf
```

## How It Works

The module sets `use_effective_fields = true` internally:

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  # ... other configuration
  use_effective_fields = true  # This is the key!

  replication_specs = [
    # ... configuration that may or may not include auto-scaling
  ]
}
```

The module also includes a data source to read effective values:

```terraform
data "mongodbatlas_advanced_cluster" "this" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
  depends_on = [mongodbatlas_advanced_cluster.this]
}
```

And exposes both configured and effective specifications through outputs:

```terraform
output "configured_specs" {
  # What you configured
}

output "effective_specs" {
  # What Atlas actually allocated (after auto-scaling)
}
```

## Usage Examples

### Example 1: Without Auto-Scaling

See `without-autoscaling/main.tf`:

```terraform
module "atlas_cluster" {
  source = "../module"

  enable_auto_scaling = false

  replication_specs = [
    {
      region_configs = [{
        electable_specs = {
          instance_size = "M10"  # Fixed size
          node_count    = 3
        }
      }]
    }
  ]
}
```

### Example 2: With Auto-Scaling

See `with-autoscaling/main.tf`:

```terraform
module "atlas_cluster" {
  source = "../module"

  enable_auto_scaling = true  # Enable auto-scaling

  replication_specs = [
    {
      region_configs = [{
        electable_specs = {
          instance_size = "M10"  # Starting size
          node_count    = 3
        }
        auto_scaling = {
          compute_enabled           = true
          compute_min_instance_size = "M10"
          compute_max_instance_size = "M30"
        }
      }]
    }
  ]
}
```

**Notice:** Same module, different configuration, no `lifecycle.ignore_changes` needed!

## Running the Examples

### Prerequisites

* Terraform >= 1.0
* MongoDB Atlas Provider v2.0.0 or later
* MongoDB Atlas account with service account credentials

### Setup

1. **Configure credentials:**

```bash
export MONGODB_ATLAS_CLIENT_ID="<your-client-id>"
export MONGODB_ATLAS_CLIENT_SECRET="<your-client-secret>"
```

Or create a `terraform.tfvars` file:

```hcl
atlas_client_id     = "<your-client-id>"
atlas_client_secret = "<your-client-secret>"
atlas_org_id        = "<your-org-id>"
```

2. **Choose an example:**

```bash
cd with-autoscaling/
# or
cd without-autoscaling/
```

3. **Initialize and apply:**

```bash
terraform init
terraform plan
terraform apply
```

4. **View outputs:**

```bash
terraform output cluster_info
terraform output configured_vs_effective
```

5. **Cleanup:**

```bash
terraform destroy
```

## Key Module Features

### Input Variables

- `enable_auto_scaling`: Enable/disable compute auto-scaling for electable nodes
- `enable_analytics_auto_scaling`: Enable/disable compute auto-scaling for analytics nodes
- `replication_specs`: Full cluster topology configuration
- `tags`: Resource tags

### Outputs

- `configured_specs`: What you configured in Terraform
- `effective_specs`: What Atlas actually allocated (reflects auto-scaling changes)
- `auto_scaling_enabled`: Whether auto-scaling is enabled
- `connection_strings`: Cluster connection strings
- `project_id`, `cluster_id`: Resource identifiers

## Learn More

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields)
- [Advanced Cluster Resource Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)

## Migration from Legacy Modules

If you have existing modules that use `lifecycle.ignore_changes`:

1. Add `use_effective_fields = true` to your cluster resource
2. Remove `lifecycle.ignore_changes` blocks
3. Add a data source to read effective values
4. Update outputs to expose both configured and effective specs

Your module users won't need to change their code - it's a backward-compatible improvement!
