output "mongodb_connection_strings" {
  value       = mongodbatlas_advanced_cluster.this.connection_strings
  description = "Map of Uniform Resource Locators that point to the MongoDB database."
}

output "cluster_name" {
  value       = mongodbatlas_advanced_cluster.this.name
  description = "MongoDB Atlas cluster name"
}

output "project_id" {
  value       = mongodbatlas_advanced_cluster.this.project_id
  description = "MongoDB Atlas project id"
}

output "replication_specs" {
  value       = mongodbatlas_advanced_cluster.this.replication_specs
  description = "Replication Specs for the cluster"
}

output "mongodbatlas_advanced_cluster" {
  value       = resource.mongodbatlas_advanced_cluster.this
  description = "Full cluster configuration for mongodbatlas_advanced_cluster resource"
}
