# Azure resource group for log storage
resource "azurerm_resource_group" "log_rg" {
  name     = var.azure_resource_group_name
  location = var.azure_resource_group_location
}

# Azure storage account
resource "azurerm_storage_account" "log_storage" {
  name                     = var.azure_storage_account_name
  resource_group_name      = azurerm_resource_group.log_rg.name
  location                 = azurerm_resource_group.log_rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Grant the Atlas service principal permission to the storage account
resource "azurerm_role_assignment" "atlas_storage_access" {
  scope                = azurerm_storage_account.log_storage.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = var.azure_service_principal_id
}

# Azure storage container
resource "azurerm_storage_container" "log_container" {
  name               = var.azure_storage_container_name
  storage_account_id = azurerm_storage_account.log_storage.id
}
