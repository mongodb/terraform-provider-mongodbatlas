data "mongodbatlas_log_integration" "example" {
  project_id     = mongodbatlas_log_integration.example.project_id
  integration_id = mongodbatlas_log_integration.example.integration_id
}

output "log_integration_type" {
  value = data.mongodbatlas_log_integration.example.type
}
