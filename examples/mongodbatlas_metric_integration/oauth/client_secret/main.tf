# Export Atlas metrics to an OTLP-compatible endpoint using OAuth 2.0 client-secret authentication.
resource "mongodbatlas_project" "this" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}

resource "mongodbatlas_metric_integration" "this" {
  project_id              = mongodbatlas_project.this.id
  integration_type        = "OTEL"
  provider_type           = "CUSTOM"
  auth_type               = "OAUTH2"
  aggregation_temporality = "DELTA"
  endpoint                = var.otel_endpoint
  metric_selection        = var.metric_selection

  oauth = {
    client_auth_method = "CLIENT_SECRET"
    token_endpoint     = var.token_endpoint
    client_id          = var.client_id
    client_secret      = var.client_secret
    scopes             = var.oauth_scopes
  }
}

data "mongodbatlas_metric_integration" "this" {
  project_id            = mongodbatlas_metric_integration.this.project_id
  metric_integration_id = mongodbatlas_metric_integration.this.metric_integration_id
}

data "mongodbatlas_metric_integrations" "this" {
  project_id = mongodbatlas_metric_integration.this.project_id
  depends_on = [mongodbatlas_metric_integration.this]
}

output "metric_integration_type" {
  description = "Type of the metric integration."
  value       = data.mongodbatlas_metric_integration.this.integration_type
}

output "metric_integration_ids" {
  description = "IDs of the metric integrations in the project."
  value       = [for r in data.mongodbatlas_metric_integrations.this.results : r.metric_integration_id]
}
