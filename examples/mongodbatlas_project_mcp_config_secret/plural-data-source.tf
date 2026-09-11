data "mongodbatlas_project_mcp_config_secrets" "this" {
  project_id    = var.project_id
  mcp_config_id = mongodbatlas_project_mcp_config.this.mcp_config_id
  depends_on    = [mongodbatlas_project_mcp_config_secret.this]
}

output "mcp_config_secrets_results" {
  value = data.mongodbatlas_project_mcp_config_secrets.this.results
}
