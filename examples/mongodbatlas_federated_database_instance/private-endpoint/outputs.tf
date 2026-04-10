output "private_endpoint_hostnames" {
  description = "Private endpoint hostnames assigned to the Federated Database Instance"
  value       = mongodbatlas_federated_database_instance.this.private_endpoint_hostnames
}
output "mongodbatlas_data_federation_instance_name" {
  description = "Name of the MongoDB Atlas Federated Database Instance"
  value       = mongodbatlas_federated_database_instance.this.name
}

output "vpc_endpoint_id" {
  description = "ID of the AWS VPC endpoint"
  value       = aws_vpc_endpoint.this.id
}
