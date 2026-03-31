resource "azurerm_resource_group" "this" {
  name     = var.azure_resource_group_name
  location = var.azure_location
}

resource "azurerm_virtual_network" "this" {
  name                = var.vnet_name
  address_space       = var.vnet_cidr
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
}

resource "azurerm_subnet" "this" {
  name                 = var.subnet_name
  resource_group_name  = azurerm_resource_group.this.name
  virtual_network_name = azurerm_virtual_network.this.name
  address_prefixes     = var.subnet_cidr

  private_endpoint_network_policies = "Disabled"
}
