resource "mongodbatlas_mcp_config" "this" {
  org_id          = var.org_id
  mcp_config_name = "example-mcp-config"
  roles           = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_mcp_config_secret" "this" {
  org_id                     = var.org_id
  mcp_config_id              = mongodbatlas_mcp_config.this.mcp_config_id
  secret_expires_after_hours = 720 # 30 days
}

output "secret_id" {
  value = mongodbatlas_mcp_config_secret.this.secret_id
}

output "secret" {
  description = "The plain-text secret value. Returned only in the create response and cannot be retrieved again."
  sensitive   = true
  value       = mongodbatlas_mcp_config_secret.this.secret
}
