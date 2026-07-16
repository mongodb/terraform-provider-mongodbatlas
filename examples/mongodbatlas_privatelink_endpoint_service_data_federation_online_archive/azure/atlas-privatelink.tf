resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "this" {
  project_id                   = var.project_id
  endpoint_id                  = azurerm_private_endpoint.this.id
  provider_name                = "AZURE"
  region                       = "US_EAST_2"
  customer_endpoint_ip_address = azurerm_private_endpoint.this.private_service_connection[0].private_ip_address
  comment                      = "Terraform Example Comment"
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
