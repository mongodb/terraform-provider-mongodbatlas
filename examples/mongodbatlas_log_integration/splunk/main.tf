# Set up log integration with authorized IAM role
resource "mongodbatlas_log_integration" "splunk" {
  project_id  = var.project_id

  type        = "SPLUNK_LOG_EXPORT"
  log_types   = var.log_types

  hec_token = var.hec_token
  hec_url   = var.hec_url
}