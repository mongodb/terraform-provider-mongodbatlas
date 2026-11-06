# Lists all workspace private endpoints registered in a project.
data "mongodbatlas_stream_workspace_private_endpoints" "example" {
  project_id = mongodbatlas_stream_workspace_private_endpoint.example.project_id
  depends_on = [mongodbatlas_stream_workspace_private_endpoint.example]
}

output "endpoint_ids" {
  value = [for e in data.mongodbatlas_stream_workspace_private_endpoints.example.results : e.endpoint_id]
}
