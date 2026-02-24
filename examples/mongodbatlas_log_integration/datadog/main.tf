# Set up log integration with Datadog
resource "mongodbatlas_log_integration" "datadog" {
  project_id  = var.project_id
  type        = "DATADOG_LOG_EXPORT"
  log_types   = ["MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"]
  api_key     = var.datadog_api_key
  region      = var.datadog_region
}