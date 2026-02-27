resource "mongodbatlas_privatelink_endpoint" "pe_east" {
  project_id               = var.project_id
  provider_name            = "AWS"
  region                   = "us-east-1"
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
}

resource "mongodbatlas_privatelink_endpoint_service" "pe_east_service" {
  project_id               = mongodbatlas_privatelink_endpoint.pe_east.project_id
  private_link_id          = mongodbatlas_privatelink_endpoint.pe_east.id
  endpoint_service_id      = aws_vpc_endpoint.vpce_east.id
  provider_name            = "AWS"
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
}

data "mongodbatlas_privatelink_endpoint" "pe_east" {
  project_id      = mongodbatlas_privatelink_endpoint.pe_east.project_id
  private_link_id = mongodbatlas_privatelink_endpoint.pe_east.private_link_id
  provider_name   = "AWS"
  depends_on      = [mongodbatlas_privatelink_endpoint_service.pe_east_service]
}

data "mongodbatlas_privatelink_endpoints" "endpoints" {
  project_id    = mongodbatlas_privatelink_endpoint.pe_east.project_id
  provider_name = "AWS"
  depends_on    = [mongodbatlas_privatelink_endpoint_service.pe_east_service]
}
