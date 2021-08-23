resource "mongodbatlas_project" "test" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_project.test.id

  azure_key_vault_config {
    enabled             = true
    client_id           = var.client_id
    azure_environment   = "AZURE"
    subscription_id     = var.subscription_id
    resource_group_name = var.resource_group_name
    key_vault_name      = var.key_vault_name
    key_identifier      = var.key_identifier
    secret              = var.client_secret
    tenant_id           = var.tenant_id
  }
}

output "project_id" {
  value = mongodbatlas_project.test.id
}
