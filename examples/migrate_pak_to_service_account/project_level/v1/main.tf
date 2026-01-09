############################################################
# v1: Initial State - PAK Resources Only 
############################################################

# Project-level Programmatic API Key
resource "mongodbatlas_project_api_key" "example" {
  description = "Description for the Project API Key"
  project_assignment {
    project_id = var.project_id
    role_names = var.project_roles
  }
}

# IP Access List entry for the API Key
resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_project_api_key.example.api_key_id
  cidr_block = var.cidr_block
  # Alternative: ip_address = "192.168.1.100"
}
