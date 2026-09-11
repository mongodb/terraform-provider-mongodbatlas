data "mongodbatlas_project_mcp_config" "this" {
  project_id    = var.project_id
  mcp_config_id = mongodbatlas_project_mcp_config.this.mcp_config_id
}

output "mcp_config_name" {
  value = data.mongodbatlas_project_mcp_config.this.mcp_config_name
}
