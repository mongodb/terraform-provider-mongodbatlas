provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

resource "mongodbatlas_project" "project" {
  name   = "AutoScalingPerShardCluster"
  org_id = var.atlas_org_id
}

# Recommended approach: Using use_effective_fields to simplify auto-scaling management
resource "mongodbatlas_advanced_cluster" "test" {
  project_id           = mongodbatlas_project.project.id
  name                 = "AutoScalingCluster"
  cluster_type         = "SHARDED"
  use_effective_fields = true

  replication_specs = [
    { # first shard
      region_configs = [
        {
          auto_scaling = {
            compute_enabled           = true
            compute_max_instance_size = "M60"
          }
          analytics_auto_scaling = {
            compute_enabled           = true
            compute_max_instance_size = "M60"
          }
          electable_specs = {
            instance_size = "M40"
            node_count    = 3
          }
          analytics_specs = {
            instance_size = "M40"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
      zone_name = "Zone 1"
    },
    { # second shard
      region_configs = [
        {
          auto_scaling = {
            compute_enabled           = true
            compute_max_instance_size = "M60"
          }
          analytics_auto_scaling = {
            compute_enabled           = true
            compute_max_instance_size = "M60"
          }
          electable_specs = {
            instance_size = "M40"
            node_count    = 3
          }
          analytics_specs = {
            instance_size = "M40"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
      zone_name = "Zone 1"
    }
  ]
}

# Read effective values to see what Atlas has scaled to
data "mongodbatlas_advanced_cluster" "test" {
  project_id           = mongodbatlas_advanced_cluster.test.project_id
  name                 = mongodbatlas_advanced_cluster.test.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.test]
}

# Output to show both configured and actual (effective) sizes
output "shard_sizes" {
  description = "Configured vs actual instance sizes for each shard"
  value = [
    for idx, spec in data.mongodbatlas_advanced_cluster.test.replication_specs : {
      shard_index           = idx
      zone_name             = spec.zone_name
      configured_size       = spec.region_configs[0].electable_specs.instance_size
      actual_electable_size = spec.region_configs[0].effective_electable_specs.instance_size
      actual_analytics_size = spec.region_configs[0].effective_analytics_specs.instance_size
    }
  ]
}
