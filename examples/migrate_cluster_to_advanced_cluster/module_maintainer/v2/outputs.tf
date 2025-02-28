output "cluster_name" {
  value       = mongodbatlas_advanced_cluster.this.name
  description = "MongoDB Atlas cluster name"
}

output "project_id" {
  value       = mongodbatlas_advanced_cluster.this.project_id
  description = "MongoDB Atlas project id"
}

# these outputs use the new data source: data.mongodbatlas_cluster.this to avoid plan changes during move
output "mongodb_connection_strings" {
  # aws_private_link and aws_private_link_srv are extras in the data source
  # they are instead exposed in mongodbatlas_cluster.this.private_endpoint
  value = [
    for conn_obj in data.mongodbatlas_cluster.this.connection_strings : {
      for key, value in conn_obj : key => value if !startswith(key, "aws_private_link")
    }
  ]
  description = "Collection of Uniform Resource Locators that point to the MongoDB database."
}

output "replication_specs" {
  value       = data.mongodbatlas_cluster.this.replication_specs
  description = "Replication Specs for cluster"
}

output "mongodbatlas_cluster" {
  value       = data.mongodbatlas_cluster.this
  description = "Full cluster configuration for mongodbatlas_cluster resource"
}
