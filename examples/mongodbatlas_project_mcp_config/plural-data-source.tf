data "mongodbatlas_project_mcp_configs" "this" {
  project_id = var.project_id
}

output "mcp_configs_results" {
  value = data.mongodbatlas_project_mcp_configs.this.results
}
