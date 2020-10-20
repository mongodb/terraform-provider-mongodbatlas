resource "mongodbatlas_private_endpoint" "atlaspl" {
  project_id    = var.atlasprojectid
  provider_name = "AWS"
  region        = var.aws_region
}

resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id             = aws_vpc.primary.id
  service_name       = mongodbatlas_private_endpoint.atlaspl.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.primary-az1.id, aws_subnet.primary-az2.id]
  security_group_ids = [aws_security_group.primary_default.id]
}

resource "mongodbatlas_private_endpoint_interface_link" "atlaseplink" {
  project_id            = mongodbatlas_private_endpoint.atlaspl.project_id
  private_link_id       = mongodbatlas_private_endpoint.atlaspl.private_link_id
  interface_endpoint_id = aws_vpc_endpoint.ptfe_service.id
}
