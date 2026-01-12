############################################################
# v2: Intermediate State - Add Service Account Resources Alongside Existing PAK Resources 
############################################################

resource "mongodbatlas_project_service_account" "this" {
  project_id                 = var.project_id
  name                       = var.service_account_name
  description                = "Description for the Project Service Account"
  roles                      = var.project_roles
  secret_expires_after_hours = var.secret_expires_after_hours
}

resource "mongodbatlas_project_service_account_access_list_entry" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  cidr_block = var.cidr_block
  # Alternative: ip_address = "192.168.1.100"
}

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
}
