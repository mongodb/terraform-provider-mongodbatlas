# This example shows how a module user would use version 1.0 of the module
# which requires lifecycle.ignore_changes for auto-scaling scenarios

module "atlas_cluster" {
  source = "../module_v1"

  atlas_org_id  = var.atlas_org_id
  project_name  = var.project_name
  cluster_name  = var.cluster_name
  cluster_type  = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"

          electable_specs = {
            instance_size = "M10"  # Initial size, but actual size may differ due to auto-scaling
            node_count    = 3
          }

          # Enable auto-scaling for compute and storage
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
    ModuleVersion = "1.0"
  }
}

# With v1, we can only see configured values, not actual provisioned values
output "cluster_info" {
  value = {
    cluster_id   = module.atlas_cluster.cluster_id
    cluster_name = module.atlas_cluster.cluster_name
    # Only configured specs available - cannot see actual instance sizes after auto-scaling
    configured_specs = module.atlas_cluster.configured_specs
  }
}
