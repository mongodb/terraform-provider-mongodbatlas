terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 2.0"
    }
  }
}

provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# Example: Using the module WITH auto-scaling
# Notice: No lifecycle.ignore_changes needed!
# The module uses use_effective_fields = true internally, making it work seamlessly
module "atlas_cluster" {
  source = "../module"

  atlas_org_id = var.atlas_org_id
  project_name = "EffectiveFieldsExample-WithAutoScale"
  cluster_name = "example-cluster-autoscale"
  cluster_type = "REPLICASET"

  # Auto-scaling is enabled
  enable_auto_scaling           = true
  enable_analytics_auto_scaling = true

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

# Outputs showing both configured and effective specs
# With auto-scaling enabled, effective specs will show the actual scaled values
output "cluster_info" {
  description = "Cluster information"
  value = {
    cluster_name         = module.atlas_cluster.cluster_name
    project_id           = module.atlas_cluster.project_id
    state                = module.atlas_cluster.cluster_state
    auto_scaling_enabled = module.atlas_cluster.auto_scaling_enabled
  }
}

output "configured_vs_effective" {
  description = "Comparison of configured and effective specifications - effective values show what Atlas has scaled to"
  value = {
    configured = module.atlas_cluster.configured_specs
    effective  = module.atlas_cluster.effective_specs
  }
}
