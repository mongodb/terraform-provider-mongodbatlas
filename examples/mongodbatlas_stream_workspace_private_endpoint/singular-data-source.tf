# Reads a single registered workspace private endpoint by its ID.
data "mongodbatlas_stream_workspace_private_endpoint" "example" {
  project_id = var.project_id
  id         = var.vpc_endpoint_id
}

output "endpoint_status" {
  value = data.mongodbatlas_stream_workspace_private_endpoint.example.status
}
