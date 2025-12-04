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
#
# PHASE 1 - BACKWARD COMPATIBLE (current implementation, data source without use_effective_fields):
# - replication_specs returns ACTUAL provisioned values from Atlas API
# - This matches module_existing behavior for seamless migration
# - Module users see no difference in outputs when migrating
#
# The data source also exposes effective_* attributes (always available for dedicated clusters):
# - effective_electable_specs.instance_size, effective_electable_specs.disk_size_gb, etc.
# - effective_analytics_specs.instance_size, effective_analytics_specs.disk_size_gb, etc.
# - effective_read_only_specs.instance_size, effective_read_only_specs.disk_size_gb, etc.
# These also return actual values when use_effective_fields is not set.
#
# PHASE 2 - BREAKING CHANGE (if data source had use_effective_fields = true, prepares for v3.x):
# - replication_specs would return CONFIGURED values (client-provided intent)
# - effective_* attributes would return ACTUAL values (Atlas-managed reality)
# - BREAKING: Module users would need to switch from replication_specs to effective_*_specs for actual values
# - Prepares for provider v3.x where this becomes default behavior
output "replication_specs" {
  description = "Cluster replication specifications (actual values for backward compatibility)"
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
