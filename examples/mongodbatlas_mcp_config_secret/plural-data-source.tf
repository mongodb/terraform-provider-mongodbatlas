data "mongodbatlas_mcp_config_secrets" "this" {
  org_id        = var.org_id
  mcp_config_id = mongodbatlas_mcp_config.this.mcp_config_id
  depends_on    = [mongodbatlas_mcp_config_secret.this]
}

output "mcp_config_secrets_results" {
  value = data.mongodbatlas_mcp_config_secrets.this.results
}
