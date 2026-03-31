resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = var.project_id
  provider_name = "AZURE"
  region        = var.atlas_region
}

resource "azurerm_private_endpoint" "this" {
  name                = "pe-atlas-datafederation-onlinearchive"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  subnet_id           = azurerm_subnet.this.id

  private_service_connection {
    name                           = mongodbatlas_privatelink_endpoint.this.private_link_service_name
    private_connection_resource_id = mongodbatlas_privatelink_endpoint.this.private_link_service_resource_id
    is_manual_connection           = true
    request_message                = "Terraform example for Atlas Data Federation private endpoint"
  }
}

resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "this" {
  project_id                    = var.project_id
  endpoint_id                   = azurerm_private_endpoint.this.id
  provider_name                 = "AZURE"
  customer_endpoint_ip_address  = azurerm_private_endpoint.this.private_service_connection[0].private_ip_address
}

data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "single" {
  project_id  = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.this.project_id
  endpoint_id = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.this.endpoint_id

  depends_on = [mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.this]
}

data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archives" "list" {
  project_id = var.project_id

  depends_on = [mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.this]
}
