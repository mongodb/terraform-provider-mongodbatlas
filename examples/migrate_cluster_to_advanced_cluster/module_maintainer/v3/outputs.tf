output "mongodb_connection_strings" {
  value       = mongodbatlas_advanced_cluster.this.connection_strings
  description = "This is the MongoDB Atlas connection strings. Note, these do not show the connection mechanism of the database details"
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
  value       = data.mongodbatlas_cluster.this.replication_specs
  description = "Replication Specs for cluster"
}

output "mongodbatlas_cluster" {
  value       = data.mongodbatlas_cluster.this
  description = "Full cluster configuration for mongodbatlas_cluster resource"
}
