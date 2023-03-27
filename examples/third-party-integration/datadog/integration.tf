resource "mongodbatlas_third_party_integration" "test-datadog" {
  project_id = var.project_id
  type       = "DATADOG"
  api_key    = var.datadog_api_key
  region     = var.datadog_region
}
