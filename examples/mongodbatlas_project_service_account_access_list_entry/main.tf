resource "mongodbatlas_project_service_account" "this" {
  project_id                 = var.project_id
  name                       = "example-project-service-account"
  description                = "Example Project Service Account"
  roles                      = ["GROUP_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

# Add IP Access List Entry to Project Service Account using CIDR Block
resource "mongodbatlas_project_service_account_access_list_entry" "cidr" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  cidr_block = "1.2.3.4/32"
}

# Add IP Access List Entry to Project Service Account using IP Address
resource "mongodbatlas_project_service_account_access_list_entry" "ip" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id
  ip_address = "2.3.4.5"
}

# Data source to read a single Project Service Account Access List entry
data "mongodbatlas_project_service_account_access_list_entry" "this" {
  project_id = mongodbatlas_project_service_account_access_list_entry.cidr.project_id
  client_id  = mongodbatlas_project_service_account_access_list_entry.cidr.client_id
  cidr_block = mongodbatlas_project_service_account_access_list_entry.cidr.cidr_block
}

output "access_list_entry_cidr_block" {
  value = data.mongodbatlas_project_service_account_access_list_entry.this.cidr_block
}

# Data source to read all Project Service Account Access List entries
data "mongodbatlas_project_service_account_access_list_entries" "this" {
  project_id = mongodbatlas_project_service_account.this.project_id
  client_id  = mongodbatlas_project_service_account.this.client_id

  depends_on = [
    mongodbatlas_project_service_account_access_list_entry.cidr,
    mongodbatlas_project_service_account_access_list_entry.ip
  ]
}

output "all_access_list_entries" {
  value = data.mongodbatlas_project_service_account_access_list_entries.this.results
}
