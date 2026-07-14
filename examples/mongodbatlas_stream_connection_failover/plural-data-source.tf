# Lists all failover connections configured for a given (primary) stream connection.
data "mongodbatlas_stream_connection_failovers" "example" {
  project_id      = var.project_id
  workspace_name  = mongodbatlas_stream_workspace.example.workspace_name
  connection_name = mongodbatlas_stream_connection_failover.example.connection_name
}

output "failover_regions" {
  value = [for r in data.mongodbatlas_stream_connection_failovers.example.results : r.region]
}
