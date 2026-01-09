############################################################
# v3: Final State - Remove PAK Resources, SA Resources Only 
############################################################

resource "mongodbatlas_project_service_account" "example" {
  project_id                 = var.project_id
  name                       = var.service_account_name
  description                = "Description for the Project Service Account"
  roles                      = var.project_roles
  secret_expires_after_hours = var.secret_expires_after_hours
}

resource "mongodbatlas_project_service_account_access_list_entry" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.example.client_id
  cidr_block = var.cidr_block
}
