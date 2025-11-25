# Resource outputs
output "stream_workspace_id" {
  description = "The stream workspace ID"
  value       = mongodbatlas_stream_workspace.example.id
}

output "stream_workspace_hostnames" {
  description = "The stream workspace hostnames"
  value       = mongodbatlas_stream_workspace.example.hostnames
}

# Data source outputs
output "workspace_from_data_source" {
  description = "Stream workspace details from data source"
  value = {
    id             = data.mongodbatlas_stream_workspace.example.id
    workspace_name = data.mongodbatlas_stream_workspace.example.workspace_name
    hostnames      = data.mongodbatlas_stream_workspace.example.hostnames
  }
}

output "total_workspaces_count" {
  description = "Total number of stream workspaces in the project"
  value       = length(data.mongodbatlas_stream_workspaces.example.results)
}
