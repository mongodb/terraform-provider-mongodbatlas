resource "azurerm_resource_group" "rg" {
  name     = var.azure_resource_group
  location = var.azure_region
}

resource "azurerm_virtual_network" "vnet" {
  name                = var.vnet_name
  address_space       = var.vnet_address_space
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_subnet" "subnet" {
  name                 = var.subnet_name
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = var.subnet_address_prefix
}

resource "azurerm_eventhub_namespace" "eventhub_ns" {
  name = var.eventhub_namespace_name
  location = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  sku = "Standard" # Minimum SKU for Private Link
  capacity = 1
}

resource "azurerm_eventhub" "eventhub" {
  name                = var.eventhub_name
  namespace_name = azurerm_eventhub_namespace.eventhub_ns.name
  resource_group_name = azurerm_resource_group.rg.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_private_dns_zone" "dns_zone" {
  name                = "privatelink.servicebus.windows.net" # should always be "privatelink.servicebus.windows.net"
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "dns_zone_link" {
  name                  = "${var.vnet_name}-dns-link"
  resource_group_name   = azurerm_resource_group.rg.name
  private_dns_zone_name = azurerm_private_dns_zone.dns_zone.name
  virtual_network_id    = azurerm_virtual_network.vnet.id
}

resource "azurerm_private_endpoint" "eventhub_endpoint" {
 name = "pe-${var.eventhub_namespace_name}"
    location = azurerm_resource_group.rg.location
    resource_group_name = azurerm_resource_group.rg.name
    subnet_id = azurerm_subnet.subnet.id

    private_service_connection {
        name = "psc-${var.eventhub_namespace_name}"
        is_manual_connection = false
        private_connection_resource_id = azurerm_eventhub_namespace.eventhub_ns.id
        subresource_names = ["namespace"]
    }

    private_dns_zone_group {
        name = "default-dns-group"
        private_dns_zone_ids = [azurerm_private_dns_zone.dns_zone.id]
    }

    depends_on = [azurerm_private_dns_zone_virtual_network_link.dns_zone_link]
}

data "azurerm_client_config" "current" {}
