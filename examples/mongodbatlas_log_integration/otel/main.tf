# Set up log integration with OTel endpoint
resource "mongodbatlas_log_integration" "otel" {
  project_id  = var.project_id
  type        = "OTEL_LOG_EXPORT"
  log_types   = ["MONGOD"]
  otel_endpoint = "https://otelexample.com:1234/v1/logs"
  otel_supplied_headers = [{
    name = "header-0"
    value = "header-0-val"
  }, {
    name = "header-1"
    value = "header-1-val1"
  }]
}