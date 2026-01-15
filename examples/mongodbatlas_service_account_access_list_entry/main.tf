resource "mongodbatlas_service_account" "this" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

# Add IP Access List Entry to Service Account using CIDR Block
resource "mongodbatlas_service_account_access_list_entry" "cidr" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.this.client_id
  cidr_block = "1.2.3.4/32"
}

# Add IP Access List Entry to Service Account using IP Address
resource "mongodbatlas_service_account_access_list_entry" "ip" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.this.client_id
  ip_address = "2.3.4.5"
}

# Data source to read a single service account access list entry
data "mongodbatlas_service_account_access_list_entry" "test" {
  org_id     = mongodbatlas_service_account_access_list_entry.cidr.org_id
  client_id  = mongodbatlas_service_account_access_list_entry.cidr.client_id
  cidr_block = mongodbatlas_service_account_access_list_entry.cidr.cidr_block
}

# Data source to read all service account access list entries
data "mongodbatlas_service_account_access_list_entries" "test" {
  org_id    = mongodbatlas_service_account.this.org_id
  client_id = mongodbatlas_service_account.this.client_id
}

output "access_list_entry_cidr_block" {
  value = data.mongodbatlas_service_account_access_list_entry.test.cidr_block
}

output "all_access_list_entries" {
  value = data.mongodbatlas_service_account_access_list_entries.test.results
}
