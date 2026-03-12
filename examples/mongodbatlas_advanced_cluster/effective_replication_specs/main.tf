resource "mongodbatlas_advanced_cluster" "example" {
  project_id   = var.project_id
  name         = var.cluster_name
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
            instance_size = var.instance_size
          }
          auto_scaling = {
            compute_enabled           = true
            compute_max_instance_size = "M40"
          }
        }
      ]
    }
  ]
}

# Data source provides both configured and effective values
data "mongodbatlas_advanced_cluster" "example" {
  project_id = mongodbatlas_advanced_cluster.example.project_id
  name       = mongodbatlas_advanced_cluster.example.name
}

# Output the configured instance size from the resource
output "configured_instance_size" {
  description = "User-configured instance size from replication_specs"
  value       = mongodbatlas_advanced_cluster.example.replication_specs[0].region_configs[0].electable_specs.instance_size
}

# Output the effective (actual running) instance size from the data source
output "effective_instance_size" {
  description = "Actual running instance size from effective_replication_specs"
  value       = data.mongodbatlas_advanced_cluster.example.effective_replication_specs[0].region_configs[0].electable_specs.instance_size
}

output "cluster_id" {
  description = "The cluster ID"
  value       = mongodbatlas_advanced_cluster.example.cluster_id
}
