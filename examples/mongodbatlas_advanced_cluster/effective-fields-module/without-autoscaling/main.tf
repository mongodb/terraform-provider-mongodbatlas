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

# This example demonstrates module usage with fixed cluster specifications (auto-scaling disabled).
module "atlas_cluster" {
  source = "../module"

  atlas_org_id = var.atlas_org_id
  project_name = "EffectiveFieldsExample-NoAutoScale"
  cluster_name = "example-cluster-fixed"
  cluster_type = "REPLICASET"

  # Disable auto-scaling to maintain fixed cluster specifications
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

# Module outputs expose both configured and effective specifications.
# Without auto-scaling, effective specifications match configured values.
output "cluster_info" {
  description = "Basic cluster information including operational state"
  value = {
    cluster_name         = module.atlas_cluster.cluster_name
    project_id           = module.atlas_cluster.project_id
    state                = module.atlas_cluster.cluster_state
    auto_scaling_enabled = module.atlas_cluster.auto_scaling_enabled
  }
}

output "configured_vs_effective" {
  description = "Comparison of configured specifications and effective specifications as provisioned by Atlas"
  value = {
    configured = module.atlas_cluster.configured_specs
    effective  = module.atlas_cluster.effective_specs
  }
}
