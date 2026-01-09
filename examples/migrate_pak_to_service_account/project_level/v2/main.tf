############################################################
# v2: Intermediate State - Add Service Account Resources Alongside Existing PAK Resources 
############################################################

# Project Service Account (new)
resource "mongodbatlas_project_service_account" "example" {
  project_id                 = var.project_id
  name                       = var.service_account_name
  description                = "Description for the Project Service Account"
  roles                      = var.project_roles
  secret_expires_after_hours = var.secret_expires_after_hours
}

# Project Service Account Access List Entry (new)
resource "mongodbatlas_project_service_account_access_list_entry" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.example.client_id
  cidr_block = var.cidr_block
  # Alternative: ip_address = "192.168.1.100"
}

# Keep existing PAK resources (for now)
resource "mongodbatlas_project_api_key" "example" {
  description = "Description for the Project API Key"
  project_assignment {
    project_id = var.project_id
    role_names = var.project_roles
  }
}

resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_project_api_key.example.api_key_id
  cidr_block = var.cidr_block
}
