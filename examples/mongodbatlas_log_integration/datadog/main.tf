resource "mongodbatlas_log_integration" "example" {
  project_id = mongodbatlas_project.project.id
  type       = "DATADOG_LOG_EXPORT"
  log_types  = ["MONGOD"]
  api_key    = var.datadog_api_key
  region     = var.datadog_region
}
