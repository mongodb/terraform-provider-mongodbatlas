
output "cluster_id" {
  description = "Unique 24-hexadecimal digit string that identifies the cluster."
  value       = mongodbatlas_advanced_cluster.this.cluster_id
}


output "config_server_type" {
  description = "Describes a sharded cluster's config server type."
  value       = mongodbatlas_advanced_cluster.this.config_server_type
}


output "connection_strings" {
  description = "Collection of Uniform Resource Locators that point to the MongoDB database."
  value       = mongodbatlas_advanced_cluster.this.connection_strings
}


output "create_date" {
  description = "Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC."
  value       = mongodbatlas_advanced_cluster.this.create_date
}


output "mongo_db_version" {
  description = "Version of MongoDB that the cluster runs."
  value       = mongodbatlas_advanced_cluster.this.mongo_db_version
}


output "state_name" {
  description = "Human-readable label that indicates the current operating condition of this cluster."
  value       = mongodbatlas_advanced_cluster.this.state_name
}
