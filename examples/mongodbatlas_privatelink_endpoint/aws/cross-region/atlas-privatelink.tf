# Create a single AWS PrivateLink endpoint service with cross-region support.
# Instead of creating one endpoint service per region (see cluster-geosharded example),
# this creates one endpoint service that accepts connections from both regions.
resource "mongodbatlas_privatelink_endpoint" "pe_east" {
  project_id               = var.project_id
  provider_name            = "AWS"
  region                   = var.aws_region_east
  supported_remote_regions = [var.aws_region_west]
}

# Connect from the primary region (us-east-1).
resource "mongodbatlas_privatelink_endpoint_service" "pe_east_service" {
  project_id          = mongodbatlas_privatelink_endpoint.pe_east.project_id
  endpoint_service_id = aws_vpc_endpoint.vpce_east.id
  private_link_id     = mongodbatlas_privatelink_endpoint.pe_east.id
  provider_name       = "AWS"
}

# Connect from the remote region (us-west-2) using the same endpoint service.
resource "mongodbatlas_privatelink_endpoint_service" "pe_west_service" {
  project_id          = mongodbatlas_privatelink_endpoint.pe_east.project_id
  endpoint_service_id = aws_vpc_endpoint.vpce_west.id
  private_link_id     = mongodbatlas_privatelink_endpoint.pe_east.id
  provider_name       = "AWS"
}
