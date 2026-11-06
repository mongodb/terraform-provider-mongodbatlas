# Registers an existing AWS interface endpoint from your VPC for
# inbound private connectivity to Atlas Stream Processing in the project.
resource "mongodbatlas_stream_workspace_private_endpoint" "example" {
  project_id        = var.project_id
  endpoint_id       = var.vpc_endpoint_id
  region            = "US_EAST_1"
  customer_dns_name = var.vpc_endpoint_dns_name
}

# Reads a single registered workspace private endpoint by its ID.
data "mongodbatlas_stream_workspace_private_endpoint" "example" {
  project_id  = mongodbatlas_stream_workspace_private_endpoint.example.project_id
  endpoint_id = mongodbatlas_stream_workspace_private_endpoint.example.endpoint_id
  depends_on  = [mongodbatlas_stream_workspace_private_endpoint.example]
}

# Lists all workspace private endpoints registered in a project.
data "mongodbatlas_stream_workspace_private_endpoints" "example" {
  project_id = mongodbatlas_stream_workspace_private_endpoint.example.project_id
  depends_on = [mongodbatlas_stream_workspace_private_endpoint.example]
}

output "endpoint_status" {
  value = data.mongodbatlas_stream_workspace_private_endpoint.example.status
}

output "endpoint_ids" {
  value = [for e in data.mongodbatlas_stream_workspace_private_endpoints.example.results : e.endpoint_id]
}
