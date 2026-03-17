data "mongodbatlas_log_integration" "example" {
  project_id     = mongodbatlas_log_integration.example.project_id
  integration_id = mongodbatlas_log_integration.example.integration_id
}

output "log_integration_storage_container_name" {
  value = data.mongodbatlas_log_integration.example.storage_container_name
}
