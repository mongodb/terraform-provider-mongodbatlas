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

// Uses data source to get real provisioned values from Atlas API.
output "replication_specs" {
  description = "Cluster replication specifications (actual values from Atlas)"
  value       = data.mongodbatlas_advanced_cluster.this.replication_specs
}
