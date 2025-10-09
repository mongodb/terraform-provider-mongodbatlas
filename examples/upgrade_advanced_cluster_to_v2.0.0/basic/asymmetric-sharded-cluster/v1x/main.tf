resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = var.project_id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  backup_enabled = true

  replication_specs { # shard 1 - M30 instance size
    region_configs {
      electable_specs {
        instance_size = "M30"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # shard 2 - M30 instance size
    region_configs {
      electable_specs {
        instance_size = "M30"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # shard 3 - M40 instance size
    region_configs {
      electable_specs {
        instance_size = "M40"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # shard 4 - M40 instance size
    region_configs {
      electable_specs {
        instance_size = "M40"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  advanced_configuration {
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }

  tags {
    key   = "environment"
    value = "dev"
  }
}
