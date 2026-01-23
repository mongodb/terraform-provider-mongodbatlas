############################################################
# v1: Initial State - PAK Resources Only 
############################################################

resource "mongodbatlas_api_key" "this" {
  org_id      = var.org_id
  description = "Description for the Organization API Key"
  role_names  = var.org_roles
}

resource "mongodbatlas_api_key_project_assignment" "this" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  roles      = var.project_roles
}

resource "mongodbatlas_access_list_api_key" "this" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  cidr_block = var.cidr_block
  # Alternative: ip_address = "192.168.1.100"
}
