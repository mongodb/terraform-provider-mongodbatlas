provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
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

  replication_specs { # shard 2 - M20 instance size
    region_configs {
      electable_specs {
        instance_size = "M20"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # shard 3 - M10 instance size
    region_configs {
      electable_specs {
        instance_size = "M10"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_SOUTH_2"
    }
  }

  replication_specs { # shard 4 - M10 instance size
    region_configs {
      electable_specs {
        instance_size = "M10"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_SOUTH_2"
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

resource "mongodbatlas_project" "project" {
  name   = "Asymmetric Sharded Cluster"
  org_id = var.atlas_org_id
}
