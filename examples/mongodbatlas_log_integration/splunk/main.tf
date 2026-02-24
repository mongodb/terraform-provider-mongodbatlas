# Set up log integration with authorized IAM role
resource "mongodbatlas_log_integration" "splunk" {
  project_id  = var.project_id

  type        = "SPLUNK_LOG_EXPORT"
  log_types   = ["MONGOD_AUDIT"] # TODO take as variable

  hec_token = "test-hec-token" # TODO take as variable
  hec_url = "https://test-hec-url.com:1234" # TODO take as variable
}