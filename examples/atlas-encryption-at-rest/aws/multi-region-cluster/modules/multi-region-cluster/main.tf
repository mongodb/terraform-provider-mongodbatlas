resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = var.project_id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  backup_enabled = true

  replication_specs {
    num_shards = 3 # 3-shard Multi-Cloud Cluster

    region_configs { # shard n1 
      electable_specs {
        instance_size = var.instance_size
        node_count    = 3
      }
      analytics_specs {
        instance_size = var.instance_size
        node_count    = 1
      }
      provider_name = var.provider
      priority      = 7
      region_name   = var.aws_region_shard_1
    }

    region_configs { # shard n2
      electable_specs {
        instance_size = var.instance_size
        node_count    = 2
      }
      analytics_specs {
        instance_size = var.instance_size
        node_count    = 1
      }
      provider_name = var.provider
      priority      = 6
      region_name   = var.aws_region_shard_2
    }

    region_configs { # shard n3
      electable_specs {
        instance_size = var.instance_size
        node_count    = 2
      }
      analytics_specs {
        instance_size = var.instance_size
        node_count    = 1
      }
      provider_name = var.provider
      priority      = 0
      region_name   = var.aws_region_shard_3
    }
  }

  advanced_configuration {
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }
}
