resource "mongodbatlas_private_endpoint_regional_mode" "test" {
  project_id = var.project_id
  enabled    = true
}

resource "mongodbatlas_privatelink_endpoint" "pe_east" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = var.aws_region_east
}

resource "mongodbatlas_privatelink_endpoint" "pe_west" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = var.aws_region_west
}

resource "mongodbatlas_privatelink_endpoint_service" "pe_west_service" {
  project_id          = mongodbatlas_privatelink_endpoint.pe_west.project_id
  endpoint_service_id = aws_vpc_endpoint.vpce_west.id
  private_link_id     = mongodbatlas_privatelink_endpoint.pe_west.id
  provider_name       = "AWS"
}

resource "mongodbatlas_privatelink_endpoint_service" "pe_east_service" {
  project_id          = mongodbatlas_privatelink_endpoint.pe_east.project_id
  endpoint_service_id = aws_vpc_endpoint.vpce_east.id
  private_link_id     = mongodbatlas_privatelink_endpoint.pe_east.id
  provider_name       = "AWS"
}

