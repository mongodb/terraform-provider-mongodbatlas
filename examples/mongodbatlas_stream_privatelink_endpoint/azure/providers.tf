provider "azurerm" {
  features {}
  # assumes Azure CLI login ('az login') or other standard auth
}

provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}