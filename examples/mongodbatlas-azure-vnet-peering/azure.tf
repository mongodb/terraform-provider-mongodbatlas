provider "azurerm" {
  # The "feature" block is required for AzureRM provider 2.x. 
  # If you're using version 1.x, the "features" block is not allowed.
  version = "~>2.0"
  features {}
}
data "azurerm_client_config" "current" {
}

output "subscription_id" {
  value = data.azurerm_client_config.current.subscription_id
}
output "tenant_id" {
  value = data.azurerm_client_config.current.tenant_id
}
resource "azurerm_resource_group" "resourcegroup" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_virtual_network" "vnet" {
  name                = var.vnet_name
  resource_group_name = azurerm_resource_group.resourcegroup.name
  address_space       = [var.address_space]
  location            = var.location
}
resource "azuread_service_principal" "sp" {
  application_id               = var.application_id
  app_role_assignment_required = true
}
resource "azurerm_role_definition" "rd" {
  name  = "my-custom-role-definition"
  scope = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.resource_group_name}/providers/Microsoft.Network/virtualNetworks/${var.vnet_name}"

  permissions {
    actions = ["Microsoft.Network/virtualNetworks/virtualNetworkPeerings/read",
      "Microsoft.Network/virtualNetworks/virtualNetworkPeerings/write",
      "Microsoft.Network/virtualNetworks/virtualNetworkPeerings/delete",
    "Microsoft.Network/virtualNetworks/peer/action"]
    not_actions = []
  }

  assignable_scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.resource_group_name}/providers/Microsoft.Network/virtualNetworks/${var.vnet_name}",
  ]
}
resource "azurerm_role_assignment" "ra" {
  scope              = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.resource_group_name}/providers/Microsoft.Network/virtualNetworks/${var.vnet_name}"
  role_definition_id = azurerm_role_definition.rd.role_definition_resource_id
  principal_id       = azuread_service_principal.sp.id
}
