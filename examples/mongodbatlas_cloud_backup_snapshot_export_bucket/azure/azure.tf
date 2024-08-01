resource "azuread_service_principal" "mongo" {
  client_id                    = var.azure_atlas_app_id
  app_role_assignment_required = false
}

# Define the resource group
resource "azurerm_resource_group" "test_resource_group" {
  name     = "mongo-test-resource-group"
  location = var.azure_resource_group_location
}

resource "azurerm_storage_account" "test_storage_account" {
  name                     = var.storage_account_name
  resource_group_name      = azurerm_resource_group.test_resource_group.name
  location                 = azurerm_resource_group.test_resource_group.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test_storage_container" {
  name                  = "mongo-test-storage-container"
  storage_account_name  = azurerm_storage_account.test_storage_account.name
  container_access_type = "private"
}

resource "azurerm_role_assignment" "test_role_assignment" {
  principal_id   = azuread_service_principal.mongo.id
  role_definition_name = "Storage Blob Data Contributor"
  scope          = azurerm_storage_account.test_storage_account.id
}
