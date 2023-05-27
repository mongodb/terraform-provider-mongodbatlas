resource "mongodbatlas_organization" "test" {
  org_owner_id = var.org_owner_id
  name         = "testCreateORG"
  description  = "test API key from Org Creation Test"
  role_names   = ["ORG_OWNER"]
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