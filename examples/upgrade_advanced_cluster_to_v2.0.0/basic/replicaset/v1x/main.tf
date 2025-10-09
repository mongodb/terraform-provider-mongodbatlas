resource "mongodbatlas_advanced_cluster" "this" {
  project_id   = var.project_id
  name         = var.cluster_name
  cluster_type = "REPLICASET"

  # v1.x legacy schema (blocks)
  retain_backups_enabled = true
  disk_size_gb           = var.disk_size_gb

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = var.instance_size
        node_count    = var.node_count_electable
      }
      analytics_specs {
        instance_size = var.instance_size
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_WEST_2"
    }
  }
}
