# Registers an existing AWS interface endpoint from your VPC for
# inbound private connectivity to Atlas Stream Processing in the project.
resource "mongodbatlas_stream_workspace_private_endpoint" "example" {
  project_id        = var.project_id
  endpoint_id       = var.vpc_endpoint_id
  region            = "US_EAST_1"
  customer_dns_name = var.vpc_endpoint_dns_name
}
