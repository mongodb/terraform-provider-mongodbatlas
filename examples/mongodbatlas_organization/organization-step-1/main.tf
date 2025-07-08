# Example configuration for creating a new MongoDB Atlas organization
# This configuration includes creation-only attributes that should NOT be used when importing

resource "mongodbatlas_organization" "test" {
  # Creation-only attributes (DO NOT use when importing):
  org_owner_id = var.org_owner_id                      # Required for creation only
  description  = "test API key from Org Creation Test" # Required for creation only
  role_names   = ["ORG_OWNER"]                         # Required for creation only

  # Creation and update attributes (can be used for both creation and import):
  name                       = "testCreateORG"
  multi_factor_auth_required = true
  restrict_employee_access   = true
  api_access_list_required   = false
  security_contact           = var.security_contact
}

output "org_id" {
  value = mongodbatlas_organization.test.org_id
}

output "org_public_key" {
  value = nonsensitive(mongodbatlas_organization.test.public_key)
}

output "org_private_key" {
  value = nonsensitive(mongodbatlas_organization.test.private_key)
}
