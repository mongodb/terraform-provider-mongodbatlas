resource "mongodbatlas_log_integration" "example" {
  project_id = mongodbatlas_project.project.id
  type       = "SPLUNK_LOG_EXPORT"
  log_types  = ["MONGOD"]
  hec_token  = var.splunk_hec_token
  hec_url    = var.splunk_hec_url
}
