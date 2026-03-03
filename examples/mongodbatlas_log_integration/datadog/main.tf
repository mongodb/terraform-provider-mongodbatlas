resource "mongodbatlas_project" "project" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}

# Set up log integration to export logs to Datadog
resource "mongodbatlas_log_integration" "example" {
  project_id = mongodbatlas_project.project.id
  type       = "DATADOG_LOG_EXPORT"
  log_types  = ["MONGOD"]
  api_key    = var.datadog_api_key
  region     = var.datadog_region
}

data "mongodbatlas_log_integration" "example" {
  project_id     = mongodbatlas_log_integration.example.project_id
  integration_id = mongodbatlas_log_integration.example.integration_id
}

data "mongodbatlas_log_integrations" "example" {
  project_id = mongodbatlas_log_integration.example.project_id
  depends_on = [mongodbatlas_log_integration.example]
}

output "log_integration_id" {
  value = data.mongodbatlas_log_integration.example.integration_id
}

output "log_integration_ids" {
  value = [for r in data.mongodbatlas_log_integrations.example.results : r.integration_id]
}
