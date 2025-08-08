resource "mongodbatlas_advanced_cluster" "atlas_cluster" {
  project_id   = var.project_id
  name         = var.atlas_cluster_name
  cluster_type = var.atlas_cluster_type

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = var.provider_instance_size_name
        node_count    = 3
      }
      analytics_specs = {
        instance_size = var.provider_instance_size_name
        node_count    = 1
      }
      provider_name = var.provider_name
      priority      = 7
      region_name   = "US_EAST_1"
    },
    {
      electable_specs = {
        instance_size = var.provider_instance_size_name
        node_count    = 2
      }
      provider_name = var.provider_name
      priority      = 6
      region_name   = "US_EAST_2"
    },
    {
      electable_specs = {
        instance_size = var.provider_instance_size_name
        node_count    = 2
      }
      provider_name = var.provider_name
      priority      = 5
      region_name   = "US_WEST_1"
    }]
  }]
}

resource "mongodbatlas_cluster_outage_simulation" "outage_simulation" {
  cluster_name = mongodbatlas_advanced_cluster.atlas_cluster.name
  project_id   = var.project_id
  outage_filters {
    cloud_provider = var.provider_name
    region_name    = "US_EAST_1"
  }

  outage_filters {
    cloud_provider = var.provider_name
    region_name    = "US_EAST_2"
  }
}
