provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}
provider "azuread" {
  tenant_id = var.azure_tenant_id
}
