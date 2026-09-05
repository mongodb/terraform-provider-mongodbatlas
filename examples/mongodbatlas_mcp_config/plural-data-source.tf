data "mongodbatlas_mcp_configs" "this" {
  org_id = var.org_id
}

output "mcp_configs_results" {
  value = data.mongodbatlas_mcp_configs.this.results
}
