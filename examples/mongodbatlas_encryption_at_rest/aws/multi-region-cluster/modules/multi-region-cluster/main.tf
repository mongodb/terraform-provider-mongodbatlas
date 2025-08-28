resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id                  = var.atlas_project_id
  name                        = var.cluster_name
  cluster_type                = "SHARDED"
  backup_enabled              = true
  encryption_at_rest_provider = var.provider_name

  replication_specs = [{
    num_shards = 2 # 2-shard Multi-Region Cluster

    region_configs = [
      { # shard n1 
        electable_specs = {
          instance_size = var.instance_size
          node_count    = 3
        }
        analytics_specs = {
          instance_size = var.instance_size
          node_count    = 1
        }
        provider_name = var.provider_name
        priority      = 7
        region_name   = var.aws_region_shard_1
      },
      { # shard n2
        electable_specs = {
          instance_size = var.instance_size
          node_count    = 2
        }
        analytics_specs = {
          instance_size = var.instance_size
          node_count    = 1
        }
        provider_name = var.provider_name
        priority      = 6
        region_name   = var.aws_region_shard_2
      }
    ]
  }]

  advanced_configuration = {
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }
}
