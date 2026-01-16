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

# Data source to read a single Access List entry for the Service Account
data "mongodbatlas_service_account_access_list_entry" "this" {
  org_id     = mongodbatlas_service_account_access_list_entry.cidr.org_id
  client_id  = mongodbatlas_service_account_access_list_entry.cidr.client_id
  cidr_block = mongodbatlas_service_account_access_list_entry.cidr.cidr_block
}

output "access_list_entry_cidr_block" {
  value = data.mongodbatlas_service_account_access_list_entry.this.cidr_block
}

# Data source to read all Access List entries for the Service Account
data "mongodbatlas_service_account_access_list_entries" "this" {
  org_id    = mongodbatlas_service_account.this.org_id
  client_id = mongodbatlas_service_account.this.client_id

  depends_on = [
    mongodbatlas_service_account_access_list_entry.cidr,
    mongodbatlas_service_account_access_list_entry.ip
  ]
}

output "all_access_list_entries" {
  value = data.mongodbatlas_service_account_access_list_entries.this.results
}
