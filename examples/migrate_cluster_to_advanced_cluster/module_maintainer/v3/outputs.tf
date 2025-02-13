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
  value       = local.use_new_replication_specs ? [] : data.mongodbatlas_cluster.this[0].replication_specs
  description = "Replication Specs for cluster, will be empty if var.replication_specs_new is set"
}

output "mongodbatlas_cluster" {
  value       = local.use_new_replication_specs ? null : data.mongodbatlas_cluster.this[0]
  description = "Full cluster configuration for mongodbatlas_cluster resource, will be null if var.replication_specs_new is set"
}

output "mongodbatlas_advanced_cluster" {
  value       = data.mongodbatlas_cluster.this
  description = "Full cluster configuration for mongodbatlas_advanced_cluster resource"
}
