# Configure the MongoDB Atlas Provider and connect via a key
provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

# Create the mongodb atlas Azure cluster
resource "mongodbatlas_advanced_cluster" "azure-cluster" {
  project_id     = var.project_id
  name           = var.name
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AZURE"
      region_name   = var.provider_region_name
      electable_specs = {
        instance_size = var.provider_instance_size_name
        node_count    = 3
      }
    }]
  }]
}

# Create the peering connection request
resource "mongodbatlas_network_peering" "test" {
  project_id            = var.project_id
  container_id          = one(values(mongodbatlas_advanced_cluster.azure-cluster.replication_specs[0].container_id))
  provider_name         = "AZURE"
  azure_directory_id    = data.azurerm_client_config.current.tenant_id
  azure_subscription_id = data.azurerm_client_config.current.subscription_id
  resource_group_name   = var.resource_group_name
  vnet_name             = var.vnet_name
  atlas_cidr_block      = var.atlas_cidr_block
}
