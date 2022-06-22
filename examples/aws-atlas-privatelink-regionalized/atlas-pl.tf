resource "mongodbatlas_private_endpoint_regional_mode" "test" {
  project_id = var.atlasprojectid
  enabled    = true
}

resource "mongodbatlas_privatelink_endpoint" "atlaspl_east" {
  project_id    = var.atlasprojectid
  provider_name = "AWS"
  region        = var.aws_region_east
}

resource "mongodbatlas_privatelink_endpoint" "atlaspl_west" {
  project_id    = var.atlasprojectid
  provider_name = "AWS"
  region        = var.aws_region_west
}

resource "aws_vpc_endpoint" "ptfe_service_west" {
  provider           = aws.west
  vpc_id             = aws_vpc.west.id
  service_name       = mongodbatlas_privatelink_endpoint.atlaspl_west.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.west.id]
  security_group_ids = [aws_security_group.west.id]
}

resource "aws_vpc_endpoint" "ptfe_service_east" {
  vpc_id             = aws_vpc.primary.id
  service_name       = mongodbatlas_privatelink_endpoint.atlaspl_east.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.primary.id]
  security_group_ids = [aws_security_group.primary_default.id]
}

resource "mongodbatlas_privatelink_endpoint_service" "atlaseplink_west" {
  project_id          = mongodbatlas_privatelink_endpoint.atlaspl_west.project_id
  endpoint_service_id = aws_vpc_endpoint.ptfe_service_west.id
  private_link_id     = mongodbatlas_privatelink_endpoint.atlaspl_west.id
  provider_name       = "AWS"
}

resource "mongodbatlas_privatelink_endpoint_service" "atlaseplink_east" {
  project_id          = mongodbatlas_privatelink_endpoint.atlaspl_east.project_id
  endpoint_service_id = aws_vpc_endpoint.ptfe_service_east.id
  private_link_id     = mongodbatlas_privatelink_endpoint.atlaspl_east.id
  provider_name       = "AWS"
}
