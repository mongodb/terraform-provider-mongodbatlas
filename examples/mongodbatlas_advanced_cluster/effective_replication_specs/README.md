# MongoDB Atlas Cluster with Effective Replication Specs (v3.0.0+)

This example demonstrates the new behavior in v3.0.0 of the MongoDB Atlas Terraform Provider where:

1. **Resource `replication_specs`** - Shows user-configured values (Optional only, not Computed)
2. **Data Source `effective_replication_specs`** - Shows actual running configuration as computed by Atlas

## Breaking Changes in v3.0.0

- `use_effective_fields` attribute has been removed
- `effective_electable_specs`, `effective_read_only_specs`, and `effective_analytics_specs` have been removed from region_configs
- New `effective_replication_specs` attribute available in data sources only
- `replication_specs` children (`auto_scaling`, `analytics_auto_scaling`, `electable_specs`, `read_only_specs`, `analytics_specs`) are now Optional only (no longer Computed)

## Migration Warning

If you previously removed `read_only_specs` or `analytics_specs` from your configuration while using v2.x, those specs were retained in the state. In v3.0.0, they will be deleted from your cluster. Review your configuration carefully before upgrading.

## Usage

```hcl
# Resource shows user-configured values
resource "mongodbatlas_advanced_cluster" "example" {
  project_id   = var.project_id
  name         = "example-cluster"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          electable_specs = {
            node_count    = 3
            instance_size = "M10"  # User-configured value
          }
          auto_scaling = {
            compute_enabled          = true
            compute_max_instance_size = "M40"  # Auto-scaling may change instance_size
          }
        }
      ]
    }
  ]
}

# Data source shows both configured and actual values
data "mongodbatlas_advanced_cluster" "example" {
  project_id = mongodbatlas_advanced_cluster.example.project_id
  name       = mongodbatlas_advanced_cluster.example.name
}

# Access configured values
output "configured_instance_size" {
  value = data.mongodbatlas_advanced_cluster.example.replication_specs[0].region_configs[0].electable_specs.instance_size
}

# Access actual running values
output "effective_instance_size" {
  value = data.mongodbatlas_advanced_cluster.example.effective_replication_specs[0].region_configs[0].electable_specs.instance_size
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| project_id | MongoDB Atlas Project ID | string | yes |
| cluster_name | Name of the cluster | string | yes |
| instance_size | Instance size for the cluster | string | no |

## Outputs

| Name | Description |
|------|-------------|
| cluster_id | The cluster ID |
| configured_instance_size | User-configured instance size |
| effective_instance_size | Actual running instance size |
