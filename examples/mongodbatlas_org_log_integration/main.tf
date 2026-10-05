resource "mongodbatlas_org_log_integration" "example" {
  org_id                = var.atlas_org_id
  type                  = "OTEL_LOG_EXPORT"
  log_types             = ["EVENTS"]
  otel_endpoint         = var.otel_endpoint
  otel_supplied_headers = var.otel_supplied_headers
}
