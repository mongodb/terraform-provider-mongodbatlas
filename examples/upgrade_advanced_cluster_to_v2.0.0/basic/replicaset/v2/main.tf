resource "mongodbatlas_advanced_cluster" "this" {
  project_id             = var.project_id
  name                   = var.cluster_name
  cluster_type           = "REPLICASET"
  retain_backups_enabled = true

  replication_specs = [{  # replication_specs is now a list of objects instead of blocks
    region_configs = [{   # region_configs is now a list of objects instead of blocks
      electable_specs = { # electable_specs is now an attribute instead of a block
        instance_size = var.instance_size
        node_count    = var.node_count_electable
        disk_size_gb  = var.disk_size_gb # disk_size_gb moved from root-level into inner specs
      }
      analytics_specs = { # analytics_specs is now an attribute instead of a block
        instance_size = var.instance_size
        node_count    = 1
        disk_size_gb  = var.disk_size_gb # disk_size_gb moved from root-level into inner specs
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
  }]
}
