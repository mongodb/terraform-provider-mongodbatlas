output "project_id" {
  description = "The ID of the created Atlas project"
  value       = mongodbatlas_project.this.id
}

output "project_name" {
  description = "The name of the created Atlas project"
  value       = mongodbatlas_project.this.name
}

output "cluster_id" {
  description = "The ID of the created cluster"
  value       = mongodbatlas_advanced_cluster.this.cluster_id
}

output "cluster_name" {
  description = "The name of the created cluster"
  value       = mongodbatlas_advanced_cluster.this.name
}

output "cluster_state" {
  description = "Current state of the cluster"
  value       = mongodbatlas_advanced_cluster.this.state_name
}

output "connection_strings" {
  description = "Connection strings for the cluster"
  value       = mongodbatlas_advanced_cluster.this.connection_strings
  sensitive   = true
}

output "configured_specs" {
  description = "Configured specifications for each shard/replication spec"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      zone_name = spec.zone_name
      regions = [
        for region in spec.region_configs : {
          region_name       = region.region_name
          provider_name     = region.provider_name
          electable_size    = region.electable_specs.instance_size
          electable_count   = region.electable_specs.node_count
          analytics_size    = try(region.analytics_specs.instance_size, null)
          analytics_count   = try(region.analytics_specs.node_count, null)
          read_only_size    = try(region.read_only_specs.instance_size, null)
          read_only_count   = try(region.read_only_specs.node_count, null)
        }
      ]
    }
  ]
}

output "effective_specs" {
  description = "Effective (actual) specifications after Atlas auto-scaling. Always available regardless of whether auto-scaling is enabled"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      zone_name = spec.zone_name
      regions = [
        for region in spec.region_configs : {
          region_name              = region.region_name
          provider_name            = region.provider_name
          effective_electable_size = region.effective_electable_specs.instance_size
          effective_analytics_size = try(region.effective_analytics_specs.instance_size, null)
          effective_read_only_size = try(region.effective_read_only_specs.instance_size, null)
        }
      ]
    }
  ]
}

output "auto_scaling_enabled" {
  description = "Whether auto-scaling is enabled on this cluster"
  value       = var.enable_auto_scaling || var.enable_analytics_auto_scaling
}
