provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}
provider "azuread" {
  tenant_id = var.azure_tenant_id
}
provider "azurerm" {
  subscription_id = var.subscription_id
  client_id       = var.client_id
  client_secret   = var.client_secret
  tenant_id       = var.tenant_id
  features {
  }
}