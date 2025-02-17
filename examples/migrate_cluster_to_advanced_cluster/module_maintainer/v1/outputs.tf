output "mongodb_connection_strings" {
  value       = mongodbatlas_cluster.this.connection_strings
  description = "Collection of Uniform Resource Locators that point to the MongoDB database."
}

output "cluster_name" {
  value       = mongodbatlas_cluster.this.name
  description = "MongoDB Atlas cluster name"
}

output "project_id" {
  value       = mongodbatlas_cluster.this.project_id
  description = "MongoDB Atlas project id"
}

output "replication_specs" {
  value       = mongodbatlas_cluster.this.replication_specs
  description = "Replication Specs for cluster"
}

output "mongodbatlas_cluster" {
  value       = mongodbatlas_cluster.this
  description = "Full cluster configuration for mongodbatlas_cluster resource"
}
