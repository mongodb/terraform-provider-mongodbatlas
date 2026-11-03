# Lists all workspace private endpoints registered in a project.
data "mongodbatlas_stream_workspace_private_endpoints" "example" {
  project_id = var.project_id
}

output "endpoint_ids" {
  value = [for e in data.mongodbatlas_stream_workspace_private_endpoints.example.results : e.id]
}
