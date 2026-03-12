resource "mongodbatlas_log_integration" "example" {
  project_id            = mongodbatlas_project.project.id
  type                  = "OTEL_LOG_EXPORT"
  log_types             = ["MONGOD"]
  otel_endpoint         = var.otel_endpoint
  otel_supplied_headers = var.otel_supplied_headers
}
