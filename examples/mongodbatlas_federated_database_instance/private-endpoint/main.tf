resource "aws_vpc" "this" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_subnet" "this" {
  vpc_id            = aws_vpc.this.id
  cidr_block        = var.subnet_cidr
  availability_zone = var.availability_zone
}

data "aws_security_group" "default" {
  name   = "default"
  vpc_id = aws_vpc.this.id
}

resource "aws_vpc_endpoint" "this" {
  vpc_id             = aws_vpc.this.id
  service_name       = var.vpce_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.this.id]
  security_group_ids = [data.aws_security_group.default.id]
}

resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "this" {
  project_id                 = var.project_id
  endpoint_id                = aws_vpc_endpoint.this.id
  provider_name              = "AWS"
  region                     = var.atlas_region
  customer_endpoint_dns_name = aws_vpc_endpoint.this.dns_entry[0].dns_name
}

resource "mongodbatlas_federated_database_instance" "this" {
  project_id = var.project_id
  name       = var.federated_instance_name

  data_process_region {
    cloud_provider = "AWS"
    region         = "VIRGINIA_USA"
  }

  depends_on = [mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.this]
}
