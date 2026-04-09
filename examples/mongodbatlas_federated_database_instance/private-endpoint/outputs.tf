output "private_endpoint_hostnames" {
  description = "Private endpoint hostnames assigned to the Federated Database Instance"
  value       = mongodbatlas_federated_database_instance.this.private_endpoint_hostnames
}
