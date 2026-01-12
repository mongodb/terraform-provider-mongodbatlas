############################################################
# v1: Initial State - PAK Resources Only 
############################################################

resource "mongodbatlas_project_api_key" "this" {
  description = "Description for the Project API Key"
  project_assignment {
    project_id = var.project_id
    role_names = var.project_roles
  }
}

resource "mongodbatlas_access_list_api_key" "this" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_project_api_key.this.api_key_id
  cidr_block = var.cidr_block
  # Alternative: ip_address = "192.168.1.100"
}
