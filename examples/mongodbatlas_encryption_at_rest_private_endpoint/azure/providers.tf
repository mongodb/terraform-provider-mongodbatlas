provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

provider "azapi" {
  tenant_id       = var.azure_tenant_id
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
}
