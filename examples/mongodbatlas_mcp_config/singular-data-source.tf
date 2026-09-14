data "mongodbatlas_mcp_config" "this" {
  org_id        = var.org_id
  mcp_config_id = mongodbatlas_mcp_config.this.mcp_config_id
}

output "mcp_config_name" {
  value = data.mongodbatlas_mcp_config.this.mcp_config_name
}
