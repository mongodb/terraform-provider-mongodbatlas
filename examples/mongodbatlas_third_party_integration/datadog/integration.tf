resource "mongodbatlas_third_party_integration" "test-datadog" {
  project_id = var.project_id
  type       = "DATADOG"
  api_key    = var.datadog_api_key
  region     = var.datadog_region

  send_collection_latency_metrics  = var.send_collection_latency_metrics
  send_database_metrics            = var.send_database_metrics
  send_user_provided_resource_tags = var.send_user_provided_resource_tags
}
