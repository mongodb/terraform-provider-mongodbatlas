data "mongodbatlas_serverless_instance" "aws_private_connection" {
  project_id = mongodbatlas_serverless_instance.aws_private_connection.project_id
  name       = mongodbatlas_serverless_instance.aws_private_connection.name

  depends_on = [mongodbatlas_privatelink_endpoint_service_serverless.pe_east_service]
}

locals {
  private_endpoints = coalesce(data.mongodbatlas_serverless_instance.aws_private_connection.connection_strings_private_endpoint_srv, [])
}

output "connection_strings" {
  value = local.private_endpoints
}