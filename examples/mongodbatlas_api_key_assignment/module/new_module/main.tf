resource "mongodbatlas_api_key" "this" {
  org_id      = var.org_id
  description = "Module-managed API Key"
  role_names  = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_api_key_project_assignment" "this" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  roles      = var.role_names
}
