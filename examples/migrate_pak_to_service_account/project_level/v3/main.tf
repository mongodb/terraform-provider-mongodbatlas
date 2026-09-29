############################################################
# v3: Final State - Remove PAK Resources, SA Resources Only 
############################################################

resource "mongodbatlas_project_service_account" "this" {
  project_id             = var.project_id
  name                   = var.project_service_account_name
  description            = "Description for the Project Service Account"
  roles                  = var.project_roles
  without_initial_secret = true
}

# Create the secret explicitly so this configuration owns it and can rotate it later.
resource "mongodbatlas_project_service_account_secret" "this" {
  project_id                 = var.project_id
  client_id                  = mongodbatlas_project_service_account.this.client_id
  secret_expires_after_hours = var.secret_expires_after_hours
}

resource "mongodbatlas_project_service_account_access_list_entry" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  cidr_block = var.cidr_block
}
