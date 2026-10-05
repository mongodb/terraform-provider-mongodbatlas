data "mongodbatlas_org_log_integrations" "example" {
  org_id     = mongodbatlas_org_log_integration.example.org_id
  depends_on = [mongodbatlas_org_log_integration.example]
}

output "org_log_integration_ids" {
  value = [for r in data.mongodbatlas_org_log_integrations.example.results : r.integration_id]
}
