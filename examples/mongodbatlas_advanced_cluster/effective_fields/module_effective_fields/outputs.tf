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

/*
 This module uses Phase 1 approach (data source without use_effective_fields flag):
 - *_specs return actual provisioned values for backward compatibility with module_existing
 - Data source also exposes effective_*_specs attributes

 For Phase 2 approach (breaking change for v3.x preparation), see main.tf data source comments.
*/
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
