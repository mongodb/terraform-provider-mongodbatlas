resource "mongodbatlas_privatelink_endpoint_serverless" "atlaspl" {
  project_id    = var.atlasprojectid
  provider_name = "AWS"
  instance_name = mongodbatlas_serverless_instance.cluster_atlas.name
}

resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id             = aws_vpc.primary.id
  service_name       = mongodbatlas_privatelink_endpoint_serverless.atlaspl.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.primary-az1.id, aws_subnet.primary-az2.id]
  security_group_ids = [aws_security_group.primary_default.id]
}

resource "mongodbatlas_privatelink_endpoint_service_serverless" "atlaseplink" {
  project_id                 = mongodbatlas_privatelink_endpoint_serverless.atlaspl.project_id
  instance_name              = mongodbatlas_serverless_instance.cluster_atlas.name
  endpoint_id                = mongodbatlas_privatelink_endpoint_serverless.atlaspl.endpoint_id
  cloud_provider_endpoint_id = aws_vpc_endpoint.ptfe_service.id
  provider_name              = "AWS"
  comment                    = "test"

}
