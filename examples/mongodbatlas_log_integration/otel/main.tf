# Set up log integration with OTel endpoint
resource "mongodbatlas_log_integration" "otel" {
  project_id  = var.project_id
  type        = "OTEL_LOG_EXPORT"
  log_types   = ["MONGOD"]
  otel_endpoint = var.otel_endpoint
  otel_supplied_headers = var.otel_supplied_headers
}