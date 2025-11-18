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

# Example: Using the module WITHOUT auto-scaling
module "atlas_cluster" {
  source = "../module"

  atlas_org_id = var.atlas_org_id
  project_name = "EffectiveFieldsExample-NoAutoScale"
  cluster_name = "example-cluster-fixed"
  cluster_type = "REPLICASET"

  # Auto-scaling is disabled
  enable_auto_scaling           = false
  enable_analytics_auto_scaling = false

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
        }
      ]
    }
  ]

  tags = {
    environment = "development"
    example     = "without-autoscaling"
  }
}

# Outputs showing both configured and effective specs
# Even without auto-scaling, effective specs are available and match configured values
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
  description = "Comparison of configured and effective specifications"
  value = {
    configured = module.atlas_cluster.configured_specs
    effective  = module.atlas_cluster.effective_specs
  }
}
