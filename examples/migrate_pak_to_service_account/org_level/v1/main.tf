############################################################
# v1: Initial State - PAK Resources Only 
############################################################

# Organization-level Programmatic API Key
resource "mongodbatlas_api_key" "example" {
  org_id      = var.org_id
  description = "Description for the Organization API Key"
  role_names  = var.org_roles
}

# Project assignment for the API Key
resource "mongodbatlas_api_key_project_assignment" "example" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  roles      = var.project_roles
}

# IP Access List entry for the API Key
resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  cidr_block = var.cidr_block
  # Alternative: ip_address = "192.168.1.100"
}

