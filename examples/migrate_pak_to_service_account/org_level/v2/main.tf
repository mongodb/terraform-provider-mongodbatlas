############################################################
# v2: Intermediate State - Add Service Account Resources Alongside Existing PAK Resources
############################################################

# Service Account (new)
resource "mongodbatlas_service_account" "example" {
  org_id                     = var.org_id
  name                       = var.service_account_name
  description                = "Description for the Service Account"
  roles                      = var.org_roles
  secret_expires_after_hours = var.secret_expires_after_hours
}

# Service Account Project Assignment (new)
resource "mongodbatlas_service_account_project_assignment" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.example.client_id
  roles      = var.project_roles
}

# Service Account Access List Entry (new)
resource "mongodbatlas_service_account_access_list_entry" "example" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.example.client_id
  cidr_block = var.cidr_block
  # Alternative: ip_address = "192.168.1.100"
}

# Keep existing PAK resources (for now)
resource "mongodbatlas_api_key" "example" {
  org_id      = var.org_id
  description = "Description for the Organization API Key"
  role_names  = var.org_roles
}

resource "mongodbatlas_api_key_project_assignment" "example" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  roles      = var.project_roles
}

resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  cidr_block = var.cidr_block
}
