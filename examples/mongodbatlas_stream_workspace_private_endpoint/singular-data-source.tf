# Reads a single registered workspace private endpoint by its ID.
data "mongodbatlas_stream_workspace_private_endpoint" "example" {
  project_id  = mongodbatlas_stream_workspace_private_endpoint.example.project_id
  endpoint_id = mongodbatlas_stream_workspace_private_endpoint.example.endpoint_id
  depends_on  = [mongodbatlas_stream_workspace_private_endpoint.example]
}

output "endpoint_status" {
  value = data.mongodbatlas_stream_workspace_private_endpoint.example.status
}
