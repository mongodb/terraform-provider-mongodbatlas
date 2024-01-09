resource "mongodbatlas_privatelink_endpoint" "pe_east" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = "us-east-1"
}

resource "mongodbatlas_privatelink_endpoint_service" "pe_east_service" {
  project_id          = mongodbatlas_privatelink_endpoint.pe_east.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.pe_east.id
  endpoint_service_id = aws_vpc_endpoint.vpce_east.id
  provider_name       = "AWS"
}
