provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  backup_enabled = true

  replication_specs = [
    { # shard 1 - M30 instance size
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            disk_iops     = 3000
            node_count    = 3
            disk_size_gb  = 10
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    },
    { # shard 2 - M30 instance size

      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            disk_iops     = 3000
            node_count    = 3
            disk_size_gb  = 10
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    },
    { # shard 3 - M40 instance size

      region_configs = [
        {
          electable_specs = {
            instance_size = "M40"
            disk_iops     = 3000
            node_count    = 3
            disk_size_gb  = 10
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    },
    { # shard 4 - M40 instance size

      region_configs = [
        {
          electable_specs = {
            instance_size = "M40"
            disk_iops     = 3000
            node_count    = 3
            disk_size_gb  = 10
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    }
  ]

  advanced_configuration = {
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }

  tags = {
    environment = "dev"
  }
}

resource "mongodbatlas_project" "project" {
  name   = "Asymmetric Sharded Cluster"
  org_id = var.atlas_org_id
}
