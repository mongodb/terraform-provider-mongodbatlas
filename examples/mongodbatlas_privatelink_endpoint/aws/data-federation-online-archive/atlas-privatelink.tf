resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
  project_id                 = var.project_id
  endpoint_id                = aws_vpc_endpoint.vpce_east.id
  provider_name              = "AWS"
  comment                    = "Terraform Example Comment"
  region                     = "US_EAST_1"
  customer_endpoint_dns_name = aws_vpc_endpoint.vpce_east.dns_entry[0].dns_name
}
