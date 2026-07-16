# Reads a single failover connection by its ID.
data "mongodbatlas_stream_connection_failover" "example" {
  project_id             = var.project_id
  workspace_name         = mongodbatlas_stream_workspace.example.workspace_name
  connection_name        = mongodbatlas_stream_connection_failover.example.connection_name
  failover_connection_id = mongodbatlas_stream_connection_failover.example.failover_connection_id
}

output "failover_bootstrap_servers" {
  value = data.mongodbatlas_stream_connection_failover.example.bootstrap_servers
}
