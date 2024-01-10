provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}
provider "azuread" {
  tenant_id = var.azure_tenant_id
}
