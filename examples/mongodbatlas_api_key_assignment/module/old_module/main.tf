resource "mongodbatlas_project_api_key" "this" {
  description = "Legacy module-managed API Key"

  project_assignment {
    project_id = var.project_id
    role_names = var.role_names
  }
}
