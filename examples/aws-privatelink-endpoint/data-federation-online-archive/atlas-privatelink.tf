resource "mongodbatlas_privatelink_endpoint" "pe_east" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = "us-east-1"
}

resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
  project_id    = var.project_id
  endpoint_id   = aws_vpc_endpoint.vpce_east.id
  provider_name = "AWS"
  type          = "DATA_LAKE"
  comment       = "Terraform Acceptance Test"
}
