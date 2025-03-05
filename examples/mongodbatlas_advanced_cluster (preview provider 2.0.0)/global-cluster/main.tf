provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = var.cluster_name
  cluster_type = "GEOSHARDED"

  # uncomment next line to use self-managed sharding, see doc for more info
  # global_cluster_self_managed_sharding = true

  backup_enabled = true

  replication_specs = [
    { # shard 1 - zone n1
      zone_name = "zone n1"

      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
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
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
        }
      ]
    },
    { # shard 2 - zone n1

      zone_name = "zone n1"

      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
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
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
        }
      ]
    },
    { # shard 1 - zone n2

      zone_name = "zone n2"

      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        },
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "EUROPE_NORTH"
        }
      ]
    },
    { # shard 2 - zone n2

      zone_name = "zone n2"

      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        },
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "EUROPE_NORTH"
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
  name   = "Global Cluster"
  org_id = var.atlas_org_id
}
