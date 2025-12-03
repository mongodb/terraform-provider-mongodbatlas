# This example shows how a module user upgrades to version 2.0 of the module
# MIGRATION: Simply change the source path from module_v1 to module_v2
# No other changes needed - the module interface remains the same!

module "atlas_cluster" {
  source = "../module_v2"  # ONLY CHANGE: Updated from module_v1 to module_v2

  # All input variables remain exactly the same
  atlas_org_id  = var.atlas_org_id
  project_name  = var.project_name
  cluster_name  = var.cluster_name
  cluster_type  = "REPLICASET"

  # Same replication_specs configuration as v1
  replication_specs = [
    {
      region_configs = [
        {
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"

          electable_specs = {
            instance_size = "M10"  # Initial size value that won't change in Terraform state
            node_count    = 3
          }

          # Same auto-scaling configuration as v1
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

  tags = {
    Environment = "development"
    ManagedBy   = "terraform"
    ModuleVersion = "2.0"  # Updated to track module version
  }
}

# With v2, we can now see BOTH configured and actual (effective) values
output "cluster_info" {
  value = {
    cluster_id   = module.atlas_cluster.cluster_id
    cluster_name = module.atlas_cluster.cluster_name
    # Configured specs - values from your Terraform configuration
    configured_specs = module.atlas_cluster.configured_specs
    # NEW IN V2: Effective specs - actual provisioned values including auto-scaling changes
    effective_specs = module.atlas_cluster.effective_specs
  }
}

# NEW IN V2: Compare configured vs effective to see auto-scaling impact
output "auto_scaling_status" {
  description = "Shows the difference between configured and actual cluster state"
  value = {
    auto_scaling_enabled = module.atlas_cluster.auto_scaling_enabled
    configured_size      = module.atlas_cluster.configured_specs[0].regions[0].electable_size
    actual_size          = module.atlas_cluster.effective_specs[0].regions[0].effective_electable_size
    configured_disk      = module.atlas_cluster.configured_specs[0].regions[0].electable_disk
    actual_disk          = module.atlas_cluster.effective_specs[0].regions[0].effective_electable_disk
  }
}
