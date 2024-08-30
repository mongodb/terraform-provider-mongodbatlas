resource "mongodbatlas_encryption_at_rest" "ear" {
  project_id = var.atlas_project_id

  azure_key_vault_config {
    require_private_networking = true

    enabled           = true
    azure_environment = "AZURE"

    tenant_id       = var.azure_tenant_id
    subscription_id = var.azure_subscription_id
    client_id       = var.azure_client_id
    secret          = var.azure_client_secret

    resource_group_name = var.azure_resource_group_name
    key_vault_name      = var.azure_key_vault_name
    key_identifier      = var.azure_key_identifier
  }
}

# Creates private endpoint
resource "mongodbatlas_encryption_at_rest_private_endpoint" "endpoint" {
  project_id     = mongodbatlas_encryption_at_rest.ear.project_id
  cloud_provider = "AZURE"
  region_name    = var.azure_region_name
}

locals {
  key_vault_resource_id = "/subscriptions/${var.azure_subscription_id}/resourceGroups/${var.azure_resource_group_name}/providers/Microsoft.KeyVault/vaults/${var.azure_key_vault_name}"
}

# Approves private endpoint connection from Azure Key Vault
resource "azapi_update_resource" "approval" {
  type      = "Microsoft.KeyVault/Vaults/PrivateEndpointConnections@2023-07-01"
  name      = mongodbatlas_encryption_at_rest_private_endpoint.endpoint.private_endpoint_connection_name
  parent_id = local.key_vault_resource_id

  body = jsonencode({
    properties = {
      privateLinkServiceConnectionState = {
        description = "Approved via Terraform"
        status      = "Approved"
      }
    }
  })
}
