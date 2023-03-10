provider "azurerm" {
  subscription_id = var.subscription_id
  client_id       = var.client_id
  client_secret   = var.client_secret
  tenant_id       = var.tenant_id
  features {
  }
}

data "azurerm_resource_group" "test" {
  name = var.resource_group_name
}

resource "azurerm_virtual_network" "test" {
  name                = "acceptanceTestVirtualNetwork1"
  address_space       = ["10.0.0.0/16"]
  location            = data.azurerm_resource_group.test.location
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet" "test" {
  name                                           = "testsubnet"
  resource_group_name                            = var.resource_group_name
  virtual_network_name                           = azurerm_virtual_network.test.name
  address_prefixes                               = ["10.0.1.0/24"]
  enforce_private_link_service_network_policies  = true
  enforce_private_link_endpoint_network_policies = true
}

resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
  project_id    = var.project_id
  instance_name = mongodbatlas_serverless_instance.test.name
  provider_name = "AZURE"
}

resource "azurerm_private_endpoint" "test" {
  name                = "endpoint-test"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = var.resource_group_name
  subnet_id           = azurerm_subnet.test.id
  private_service_connection {
    name                           = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_service_name
    private_connection_resource_id = mongodbatlas_privatelink_endpoint_serverless.test.private_link_service_resource_id
    is_manual_connection           = true
    request_message                = "Azure Private Link test"
  }

}

resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
  project_id                  = mongodbatlas_privatelink_endpoint_serverless.test.project_id
  instance_name               = mongodbatlas_serverless_instance.test.name
  endpoint_id                 = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
  cloud_provider_endpoint_id  = azurerm_private_endpoint.test.id
  private_endpoint_ip_address = azurerm_private_endpoint.test.private_service_connection[0].private_ip_address
  provider_name               = "AZURE"
  comment                     = "test"
}

resource "mongodbatlas_serverless_instance" "test" {
  project_id                              = var.project_id
  name                                    = var.cluster_name
  provider_settings_backing_provider_name = "AZURE"
  provider_settings_provider_name         = "SERVERLESS"
  provider_settings_region_name           = "US_EAST_2"
  continuous_backup_enabled               = true
}

output "private_endpoints" {
  value = mongodbatlas_serverless_instance.test.connection_strings_private_endpoint_srv[0]
}