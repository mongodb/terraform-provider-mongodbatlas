resource "mongodbatlas_cluster" "atlas_cluster" {
  project_id   = var.project_id
  name         = var.atlas_cluster_name
  cluster_type = var.atlas_cluster_type

  provider_name               = var.provider_name
  provider_instance_size_name = var.provider_instance_size_name

  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
    regions_config {
      region_name     = "US_EAST_2"
      electable_nodes = 2
      priority        = 6
      read_only_nodes = 0
    }
    regions_config {
      region_name     = "US_WEST_1"
      electable_nodes = 2
      priority        = 5
      read_only_nodes = 2
    }
  }
}

resource "mongodbatlas_cluster_outage_simulation" "outage_simulation" {
  cluster_name = var.atlas_cluster_name
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
