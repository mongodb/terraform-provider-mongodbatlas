############################################################
# v3: Final State - Remove PAK Resources, SA Resources Only 
############################################################

resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "Name for the Service Account"
  description                = "Description for the Service Account"
  roles                      = var.org_roles
  secret_expires_after_hours = var.secret_expires_after_hours
}

resource "mongodbatlas_service_account_project_assignment" "this" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.this.client_id
  roles      = var.project_roles
}

resource "mongodbatlas_service_account_access_list_entry" "this" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.this.client_id
  cidr_block = var.cidr_block
}
