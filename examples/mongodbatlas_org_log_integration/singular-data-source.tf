data "mongodbatlas_org_log_integration" "example" {
  org_id         = mongodbatlas_org_log_integration.example.org_id
  integration_id = mongodbatlas_org_log_integration.example.integration_id
}

output "org_log_integration_type" {
  value = data.mongodbatlas_org_log_integration.example.type
}
