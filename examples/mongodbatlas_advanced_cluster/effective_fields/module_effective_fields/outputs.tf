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

# Replication specifications from data source
# With use_effective_fields, regular spec attributes (electable_specs, analytics_specs,
# read_only_specs) return configuration values that stay constant.
# To get actual provisioned values (which may differ due to auto-scaling),
# use the effective_* attributes available in data source:
# - effective_electable_specs.instance_size, effective_electable_specs.disk_size_gb, etc.
# - effective_analytics_specs.instance_size, effective_analytics_specs.disk_size_gb, etc.
# - effective_read_only_specs.instance_size, effective_read_only_specs.disk_size_gb, etc.
output "replication_specs" {
  description = "Cluster replication specifications (contains both configured and effective_* specs)"
  value       = data.mongodbatlas_advanced_cluster.this.replication_specs
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
