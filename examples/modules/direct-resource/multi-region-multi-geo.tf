# - multiple regions (different geographies)

resource "mongodbatlas_advanced_cluster" "atlas-cluster-multiregion-multigeo" {
  project_id = var.project_id
  name = "MultiRegionMultiGeoCluster"
  cluster_type = "REPLICASET"
  mongo_db_major_version = "8.0"
  replication_specs {
    region_configs {
      provider_name = "AWS"
      region_name   = "US_EAST_1" # North America
      priority      = 7
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      auto_scaling {
        disk_gb_enabled = true
        compute_enabled = true
        compute_max_instance_size = "M60"
        compute_min_instance_size = "M30"
      }
    }
    region_configs {
      provider_name = "AWS"
      region_name   = "EU_WEST_1" # Europe
      priority      = 6
      electable_specs {
        instance_size = "M30"
        node_count    = 2
      }
      auto_scaling {
        disk_gb_enabled = true
        compute_enabled = true
        compute_max_instance_size = "M60"
        compute_min_instance_size = "M30"
      }
    }
  }
}