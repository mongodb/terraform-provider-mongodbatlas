output "cluster_id" {
  description = "The ID of the created Atlas cluster."
  value       = mongodbatlas_advanced_cluster.this.id
}

output "connection_strings" {
  description = "Connection strings for the created cluster."
  value       = mongodbatlas_advanced_cluster.this.connection_strings
}
