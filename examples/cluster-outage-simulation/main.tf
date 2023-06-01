resource "mongodbatlas_advanced_cluster" "atlas_cluster" {
  project_id   = var.project_id
  name         = var.atlas_cluster_name
  cluster_type = var.atlas_cluster_type

  replication_specs {
    region_configs {
      provider_name = var.provider_name
      electable_specs {
        instance_size = var.provider_instance_size_name
      }
      region_name = "US_EAST_1"
      priority    = 7
    }

    region_configs {
      provider_name = var.provider_name
      electable_specs {
        instance_size = var.provider_instance_size_name
      }
      region_name = "US_EAST_2"
      priority    = 6
    }

    region_configs {
      provider_name = var.provider_name
      electable_specs {
        instance_size = var.provider_instance_size_name
      }
      region_name = "US_WEST_1"
      priority    = 5
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
