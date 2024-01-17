resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id   = var.project_id
  name         = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs {
    zone_name  = "Zone 1"
    num_shards = 1

    region_configs {
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      priority      = 7

      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}