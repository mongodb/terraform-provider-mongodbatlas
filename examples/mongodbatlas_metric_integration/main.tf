resource "datadog_api_key" "atlas_metrics" {
  name = "mongodb-atlas-metrics"
}

resource "mongodbatlas_metric_integration" "example" {
  project_id              = mongodbatlas_project.project.id
  integration_type        = "OTEL"
  provider_type           = "CUSTOM"
  aggregation_temporality = "DELTA"
  endpoint                = var.datadog_endpoint
  metric_selection        = ["ATLAS_STREAM_PROCESSING"]

  headers = [
    {
      name  = "dd-api-key"
      value = datadog_api_key.atlas_metrics.key
    }
  ]
}

data "mongodbatlas_metric_integration" "example" {
  project_id            = mongodbatlas_metric_integration.example.project_id
  metric_integration_id = mongodbatlas_metric_integration.example.metric_integration_id
}

output "metric_integration_type" {
  value = data.mongodbatlas_metric_integration.example.integration_type
}

data "mongodbatlas_metric_integrations" "example" {
  project_id = mongodbatlas_metric_integration.example.project_id
  depends_on = [mongodbatlas_metric_integration.example]
}

output "metric_integration_ids" {
  value = [for r in data.mongodbatlas_metric_integrations.example.results : r.metric_integration_id]
}
