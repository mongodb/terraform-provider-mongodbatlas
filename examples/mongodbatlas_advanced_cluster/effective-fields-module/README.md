# MongoDB Atlas Advanced Cluster Module with Effective Fields

This example demonstrates a reusable Terraform module for MongoDB Atlas clusters that supports both auto-scaling and non-auto-scaling configurations without requiring `lifecycle.ignore_changes` blocks.

## Overview

When creating reusable Terraform modules for MongoDB Atlas clusters with auto-scaling capabilities, module authors face a significant challenge: auto-scaling changes made by Atlas cause plan drift unless `lifecycle.ignore_changes` blocks are used. However, incorporating lifecycle blocks into modules creates inflexibility and maintenance complexity.

The `use_effective_fields` attribute addresses this by enabling a single module implementation that handles both auto-scaling and non-auto-scaling scenarios without lifecycle block requirements.

## The Challenge for Module Authors

### Without use_effective_fields

Traditional approaches require one of the following unsatisfactory solutions:

**Option 1: Separate module implementations**
- `cluster-module-with-autoscaling/` - Includes lifecycle.ignore_changes blocks for auto-scaling scenarios
- `cluster-module-without-autoscaling/` - No lifecycle blocks for fixed-size clusters

```terraform
# Inside the auto-scaling module's cluster resource
resource "mongodbatlas_advanced_cluster" "this" {
  # ... configuration

  lifecycle {
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      replication_specs[0].region_configs[0].electable_specs.disk_size_gb,
      # ... additional fields
    ]
  }
}
```

**Limitation**: Lifecycle blocks cannot be conditional, requiring two separate modules with duplicated code and increased maintenance burden.

**Option 2: Module users manage resources directly**
```terraform
# Users cannot use the module and must manage resources directly
resource "mongodbatlas_advanced_cluster" "this" {
  # ... direct configuration instead of module abstraction

  lifecycle {
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs.instance_size,
    ]
  }
}
```
**Limitation**: Eliminates the benefits of module abstraction and reusability.

### With use_effective_fields

By incorporating `use_effective_fields = true` in the module's cluster resource, a single module implementation supports both scenarios:

```terraform
# Single module works for both auto-scaling and non-auto-scaling clusters
module "cluster" {
  source = "./cluster-module"
  # Auto-scaling is automatically detected from replication_specs configuration
  # ... configuration
}
# No lifecycle.ignore_changes required
```

## Benefits

1. **Unified module implementation**: Single codebase supports all use cases
2. **Operational visibility**: Module outputs expose both configured and effective (actual) specifications
3. **Forward compatibility**: Aligns with provider v3.x where this behavior becomes default

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

## Implementation Details

### Module Resource Configuration

The module resource incorporates `use_effective_fields = true`:

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  # ... configuration
  use_effective_fields = true

  replication_specs = [
    # Configuration supports both auto-scaling and non-auto-scaling scenarios
  ]
}
```

### Data Source Integration

A companion data source reads effective (actual) values:

```terraform
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
```

### Output Exposure

The module exposes both configured and effective specifications:

```terraform
output "configured_specs" {
  description = "Hardware specifications as defined in configuration"
  # Returns user-configured values
}

output "effective_specs" {
  description = "Hardware specifications as provisioned by Atlas"
  # Returns actual operational values including auto-scaling changes
}
```

## Usage Examples

### Without Auto-Scaling

The `without-autoscaling/` directory demonstrates module usage with fixed cluster specifications:

```terraform
module "atlas_cluster" {
  source = "../module"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"  # Fixed size
            node_count    = 3
          }
          # No auto_scaling block - cluster maintains fixed specifications
        }
      ]
    }
  ]
}
```

In this configuration, cluster specifications remain constant. The effective specifications returned by the module match the configured specifications.

### With Auto-Scaling

The `with-autoscaling/` directory demonstrates module usage with auto-scaling enabled:

```terraform
module "atlas_cluster" {
  source = "../module"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10" # Initial size value that won't change in Terraform state, actual size in Atlas may differ due to auto-scaling
            node_count    = 3
          }
          auto_scaling = {
            compute_enabled           = true
            compute_min_instance_size = "M10"
            compute_max_instance_size = "M30"
          }
        }
      ]
    }
  ]
}
```

With auto-scaling enabled, Atlas adjusts instance sizes based on workload. The module automatically detects auto-scaling configuration from the `replication_specs` and exposes effective specifications that reflect Atlas-managed changes. Both configurations utilize the same module implementation without requiring `lifecycle.ignore_changes` blocks.

## Running the Examples

### Prerequisites

The examples require the following:

* Terraform >= 1.0
* MongoDB Atlas Provider v2.0.0 or later
* MongoDB Atlas service account credentials

### Configuration

Configure authentication credentials using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<your-client-id>"
export MONGODB_ATLAS_CLIENT_SECRET="<your-client-secret>"
```

Alternatively, create a `terraform.tfvars` file:

```hcl
atlas_client_id     = "<your-client-id>"
atlas_client_secret = "<your-client-secret>"
atlas_org_id        = "<your-org-id>"
```

### Execution

Navigate to the desired example directory:

```bash
cd with-autoscaling/
```

or

```bash
cd without-autoscaling/
```

Initialize and apply the configuration:

```bash
terraform init
terraform plan
terraform apply
```

Review the module outputs to observe both configured and effective specifications:

```bash
terraform output cluster_info
terraform output configured_vs_effective
```

To remove all created resources:

```bash
terraform destroy
```

## Module Interface

### Input Variables

The module accepts the following key input variables:

- `atlas_org_id` - Atlas organization identifier
- `project_name` - Atlas project name
- `cluster_name` - Atlas cluster name
- `cluster_type` - Cluster type (REPLICASET, SHARDED, or GEOSHARDED)
- `replication_specs` - Cluster topology and hardware specifications (auto-scaling is automatically detected from this configuration)
- `tags` - Key-value pairs for resource categorization

### Output Values

The module exposes the following outputs:

- `configured_specs` - Hardware specifications as defined in the configuration
- `effective_specs` - Hardware specifications as provisioned by Atlas (includes auto-scaling changes)
- `auto_scaling_enabled` - Indicates whether auto-scaling is enabled for electable and read-only nodes
- `analytics_auto_scaling_enabled` - Indicates whether auto-scaling is enabled for analytics nodes
- `connection_strings` - Connection strings for cluster access
- `project_id`, `cluster_id` - Resource identifiers

## Additional Resources

For comprehensive information about the effective fields functionality:

- [Auto-Scaling with Effective Fields](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#auto-scaling-with-effective-fields)
- [Advanced Cluster Resource Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster)
- [Advanced Cluster Data Source Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_cluster)

## Migration from Legacy Implementations

Existing modules utilizing `lifecycle.ignore_changes` can be updated through the following steps:

1. Add `use_effective_fields = true` to the cluster resource definition
2. Remove existing `lifecycle.ignore_changes` blocks
3. Incorporate a data source to retrieve effective values
4. Update module outputs to expose both configured and effective specifications

This migration maintains backward compatibility, requiring no changes to module consumer configurations.
