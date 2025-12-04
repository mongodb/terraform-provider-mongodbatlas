provider "mongodbatlas" {
  # Credentials from environment variables:
  # - MONGODB_ATLAS_CLIENT_ID
  # - MONGODB_ATLAS_CLIENT_SECRET
}

# Module usage - works with BOTH module_existing and module_effective_fields
# Only the source path needs to change when migrating
module "atlas_cluster" {
  source = "../module_effective_fields" # Switch to ../module_existing to use the old approach

  atlas_org_id = var.atlas_org_id
  project_name = var.project_name
  cluster_name = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"

          electable_specs = {
            instance_size = "M10"
            node_count    = 3
          }

          analytics_specs = {
            instance_size = "M10"
            node_count    = 1
          }

          # Auto-scaling configuration
          auto_scaling = {
            disk_gb_enabled            = true
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M10"
            compute_max_instance_size  = "M30"
          }

          analytics_auto_scaling = {
            disk_gb_enabled            = true
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M10"
            compute_max_instance_size  = "M30"
          }
        }
      ]
    }
  ]

  tags = {
    environment = "development"
    managed_by  = "terraform"
  }
}

# Access cluster information
output "cluster_id" {
  description = "Atlas cluster ID"
  value       = module.atlas_cluster.cluster_id
}

output "cluster_state" {
  description = "Atlas cluster state"
  value       = module.atlas_cluster.cluster_state
}

output "connection_strings" {
  description = "Atlas cluster connection strings"
  value       = module.atlas_cluster.connection_strings
  sensitive   = true
}

# When using module_effective_fields, you can access both configured and effective specs
output "configured_specs" {
  description = "Hardware specifications as defined in configuration"
  value       = module.atlas_cluster.configured_specs
}

# This output is only available when using module_effective_fields
# When using module_existing, this will show the same values as configured_specs
output "effective_specs" {
  description = "Actual hardware specifications provisioned by Atlas (may differ due to auto-scaling)"
  value       = try(module.atlas_cluster.effective_specs, "Not available with module_existing")
}
