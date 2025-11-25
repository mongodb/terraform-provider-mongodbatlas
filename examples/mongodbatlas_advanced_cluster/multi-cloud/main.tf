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
    { # shard 1
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          analytics_specs = {
            instance_size = "M30"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        },
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          analytics_specs = {
            instance_size = "M30"
            node_count    = 1
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
        }
      ]
    },
    { # shard 2
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          analytics_specs = {
            instance_size = "M30"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        },
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          analytics_specs = {
            instance_size = "M30"
            node_count    = 1
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
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
    "environment" = "dev"
  }
}

resource "mongodbatlas_project" "project" {
  name   = "Multi-Cloud Cluster"
  org_id = var.atlas_org_id
}
