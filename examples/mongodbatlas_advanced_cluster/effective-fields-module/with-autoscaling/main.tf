provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# This example demonstrates module usage with auto-scaling enabled.
# The module utilizes use_effective_fields = true internally, eliminating the need for lifecycle.ignore_changes blocks.
module "atlas_cluster" {
  source = "../module"

  atlas_org_id = var.atlas_org_id
  project_name = "EffectiveFieldsExample-WithAutoScale"
  cluster_name = "example-cluster-autoscale"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"

          electable_specs = {
            instance_size = "M10" # Initial size value that won't change in Terraform state, actual size in Atlas may differ due to auto-scaling
            node_count    = 3
          }

          analytics_specs = {
            instance_size = "M10"
            node_count    = 1
          }

          # Auto-scaling configuration for electable nodes
          auto_scaling = {
            disk_gb_enabled            = true
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M10"
            compute_max_instance_size  = "M30"
          }

          # Auto-scaling configuration for analytics nodes
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
    environment = "production"
    example     = "with-autoscaling"
  }
}

# Module outputs expose both configured and effective specifications.
# With auto-scaling enabled, effective specifications reflect actual scaled values.
output "cluster_info" {
  description = "Basic cluster information including operational state"
  value = {
    cluster_name                     = module.atlas_cluster.cluster_name
    project_id                       = module.atlas_cluster.project_id
    state                            = module.atlas_cluster.cluster_state
    auto_scaling_enabled             = module.atlas_cluster.auto_scaling_enabled
    analytics_auto_scaling_enabled   = module.atlas_cluster.analytics_auto_scaling_enabled
  }
}

output "configured_vs_effective" {
  description = "Comparison of configured specifications and effective specifications as provisioned by Atlas"
  value = {
    configured = module.atlas_cluster.configured_specs
    effective  = module.atlas_cluster.effective_specs
  }
}
