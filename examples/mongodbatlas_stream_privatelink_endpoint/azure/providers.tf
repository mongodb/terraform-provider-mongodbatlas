provider "azurerm" {
  features {}
  # assumes Azure CLI login ('az login') or other standard auth
}

provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}