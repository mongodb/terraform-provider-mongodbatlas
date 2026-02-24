# Set up log integration with authorized IAM role
resource "mongodbatlas_log_integration" "datadog" {
  project_id  = var.project_id
  type        = "DATADOG_LOG_EXPORT"
  log_types   = ["MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"]
  api_key = "test-dd-api-key2"
  region = "US1"
}