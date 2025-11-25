# Old module usage
module "project_api_key" {
  source     = "./old_module"
  project_id = var.project_id
  role_names = var.role_names
}

# New module usage
module "api_key_assignment" {
  source     = "./new_module"
  org_id     = var.org_id
  project_id = var.project_id
  role_names = var.role_names
}
