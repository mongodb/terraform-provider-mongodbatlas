provider "azurerm" {
  subscription_id = var.subscription_id
  client_id       = var.client_id
  client_secret   = var.client_secret
  tenant_id       = var.tenant_id
  features {
  }
}

data "azurerm_resource_group" "this" {
  name = var.resource_group_name
}

resource "azurerm_virtual_network" "this" {
  name                = "acceptanceTestVirtualNetwork1"
  address_space       = ["10.0.0.0/16"]
  location            = data.azurerm_resource_group.this.location
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet" "this" {
  name                                          = "testsubnet"
  resource_group_name                           = var.resource_group_name
  virtual_network_name                          = azurerm_virtual_network.this.name
  address_prefixes                              = ["10.0.1.0/24"]
  private_link_service_network_policies_enabled = true
  private_endpoint_network_policies_enabled     = true
}

resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id               = var.project_id
  provider_name            = "AZURE"
  region                   = "eastus2"
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
}

resource "azurerm_private_endpoint" "this" {
  name                = "endpoint-test"
  location            = data.azurerm_resource_group.this.location
  resource_group_name = var.resource_group_name
  subnet_id           = azurerm_subnet.this.id
  private_service_connection {
    name                           = mongodbatlas_privatelink_endpoint.this.private_link_service_name
    private_connection_resource_id = mongodbatlas_privatelink_endpoint.this.private_link_service_resource_id
    is_manual_connection           = true
    request_message                = "Azure Private Link test"
  }

}

resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id                  = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.this.private_link_id
  endpoint_service_id         = azurerm_private_endpoint.this.id
  private_endpoint_ip_address = azurerm_private_endpoint.this.private_service_connection[0].private_ip_address
  provider_name               = "AZURE"
  delete_on_create_timeout    = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count      = var.cluster_name == "" ? 0 : 1
  project_id = var.project_id
  name       = var.cluster_name
  depends_on = [mongodbatlas_privatelink_endpoint_service.this]
}

locals {
  private_endpoints = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings : cs.private_endpoint]), [])
  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], azurerm_private_endpoint.this.id)
  ]
}

output "connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
