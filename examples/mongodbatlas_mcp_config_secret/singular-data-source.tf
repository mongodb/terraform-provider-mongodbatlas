data "mongodbatlas_mcp_config_secret" "this" {
  org_id        = var.org_id
  mcp_config_id = mongodbatlas_mcp_config.this.mcp_config_id
  secret_id     = mongodbatlas_mcp_config_secret.this.secret_id
}

output "secret_expires_at" {
  value = data.mongodbatlas_mcp_config_secret.this.expires_at
}
