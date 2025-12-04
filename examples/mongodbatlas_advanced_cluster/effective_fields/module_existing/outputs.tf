output "project_id" {
  description = "Atlas project id"
  value       = mongodbatlas_project.this.id
}

output "project_name" {
  description = "Atlas project name"
  value       = mongodbatlas_project.this.name
}

output "cluster_id" {
  description = "Atlas cluster id"
  value       = mongodbatlas_advanced_cluster.this.cluster_id
}

output "cluster_name" {
  description = "Atlas cluster name"
  value       = mongodbatlas_advanced_cluster.this.name
}

output "cluster_state" {
  description = "Atlas cluster state"
  value       = mongodbatlas_advanced_cluster.this.state_name
}

output "connection_strings" {
  description = "Atlas cluster connection strings"
  value       = mongodbatlas_advanced_cluster.this.connection_strings
  sensitive   = true
}

# With lifecycle.ignore_changes, cannot distinguish configured from auto-scaled values
output "configured_specs" {
  description = "Hardware specifications from state (may include auto-scaled values)"
  value = [
    for spec in mongodbatlas_advanced_cluster.this.replication_specs : {
      zone_name = spec.zone_name
      regions = [
        for region in spec.region_configs : {
          region_name     = region.region_name
          provider_name   = region.provider_name
          electable_size  = region.electable_specs.instance_size
          electable_count = region.electable_specs.node_count
          electable_disk  = try(region.electable_specs.disk_size_gb, null)
          analytics_size  = try(region.analytics_specs.instance_size, null)
          analytics_count = try(region.analytics_specs.node_count, null)
          analytics_disk  = try(region.analytics_specs.disk_size_gb, null)
          read_only_size  = try(region.read_only_specs.instance_size, null)
          read_only_count = try(region.read_only_specs.node_count, null)
          read_only_disk  = try(region.read_only_specs.disk_size_gb, null)
        }
      ]
    }
  ]
}
