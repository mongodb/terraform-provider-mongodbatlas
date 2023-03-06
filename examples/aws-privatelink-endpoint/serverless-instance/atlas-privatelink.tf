resource "mongodbatlas_privatelink_endpoint_serverless" "pe_east" {
  project_id    = mongodbatlas_serverless_instance.aws_private_connection.project_id
  instance_name = mongodbatlas_serverless_instance.aws_private_connection.name
  provider_name = "AWS"
}

resource "mongodbatlas_privatelink_endpoint_service_serverless" "pe_east_service" {
  project_id                 = mongodbatlas_privatelink_endpoint_serverless.pe_east.project_id
  instance_name              = mongodbatlas_privatelink_endpoint_serverless.pe_east.instance_name
  endpoint_id                = mongodbatlas_privatelink_endpoint_serverless.pe_east.endpoint_id
  provider_name              = mongodbatlas_privatelink_endpoint_serverless.pe_east.provider_name
  cloud_provider_endpoint_id = aws_vpc_endpoint.vpce_east.id
  comment                    = "New serverless endpoint"
}