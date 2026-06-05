data "mongodbatlas_log_integrations" "example" {
  project_id = mongodbatlas_log_integration.example.project_id
  depends_on = [mongodbatlas_log_integration.example]
}

output "log_integration_ids" {
  value = [for r in data.mongodbatlas_log_integrations.example.results : r.integration_id]
}
