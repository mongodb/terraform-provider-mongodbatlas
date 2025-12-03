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

# Configured specifications - values as defined in Terraform configuration
# These remain constant regardless of auto-scaling
output "configured_specs" {
  description = "Configured hardware specifications as defined in Terraform configuration"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
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

# MIGRATION CHANGE: This output is new in v2
# Effective specifications - actual values provisioned by Atlas, including auto-scaling changes
# Module users can now observe the real cluster state
output "effective_specs" {
  description = "Effective hardware specifications as provisioned by Atlas, including auto-scaling changes"
  value = [
    for spec in data.mongodbatlas_advanced_cluster.this.replication_specs : {
      zone_name = spec.zone_name
      regions = [
        for region in spec.region_configs : {
          region_name              = region.region_name
          provider_name            = region.provider_name
          effective_electable_size = region.effective_electable_specs.instance_size
          effective_electable_disk = try(region.effective_electable_specs.disk_size_gb, null)
          effective_electable_iops = try(region.effective_electable_specs.disk_iops, null)
          effective_analytics_size = try(region.effective_analytics_specs.instance_size, null)
          effective_analytics_disk = try(region.effective_analytics_specs.disk_size_gb, null)
          effective_analytics_iops = try(region.effective_analytics_specs.disk_iops, null)
          effective_read_only_size = try(region.effective_read_only_specs.instance_size, null)
          effective_read_only_disk = try(region.effective_read_only_specs.disk_size_gb, null)
          effective_read_only_iops = try(region.effective_read_only_specs.disk_iops, null)
        }
      ]
    }
  ]
}

output "auto_scaling_enabled" {
  description = "Flag indicating if auto-scaling is enabled for electable and read-only nodes"
  value = try(
    var.replication_specs[0].region_configs[0].auto_scaling.disk_gb_enabled,
    false
    ) || try(
    var.replication_specs[0].region_configs[0].auto_scaling.compute_enabled,
    false
  )
}

output "analytics_auto_scaling_enabled" {
  description = "Flag indicating if auto-scaling is enabled for analytics nodes"
  value = try(
    var.replication_specs[0].region_configs[0].analytics_auto_scaling.disk_gb_enabled,
    false
    ) || try(
    var.replication_specs[0].region_configs[0].analytics_auto_scaling.compute_enabled,
    false
  )
}
