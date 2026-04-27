resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = var.project_id
  name         = var.cluster_name
  cluster_type = "REPLICASET"
  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AZURE"
      region_name   = "US_EAST_2"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

resource "mongodbatlas_stream_privatelink_endpoint" "this" {
  project_id    = var.project_id
  vendor        = "AZURE_BLOB_STORAGE"
  provider_name = "AZURE"
  region        = var.atlas_region
  # dns_domain follows the format `{storageAccount}.blob.core.windows.net`
  dns_domain = "${var.storage_account_name}.blob.core.windows.net"
  # service_endpoint_id follows the format `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Storage/storageAccounts/{storageAccount}`
  service_endpoint_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.azure_resource_group}/providers/Microsoft.Storage/storageAccounts/${var.storage_account_name}"
  depends_on          = [mongodbatlas_advanced_cluster.cluster, azurerm_private_endpoint.blob_endpoint]
}

output "privatelink_endpoint_id" {
  value = mongodbatlas_stream_privatelink_endpoint.this.id
}
