# - single region
# - with shards (single zone)

resource "mongodbatlas_advanced_cluster" "single_region_sharded_no_autoscaling" {
  project_id             = var.project_id
  name                   = "single-region-sharded-no-autoscaling"
  cluster_type           = "SHARDED"
  mongo_db_major_version = "8.0"
  replication_specs = [
    { # shard 1 (single zone)
      region_configs = [
        {
          electable_specs = {
            instance_size = "M40" # Independently scaled shard
            node_count    = 3
          }
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"
        }
      ]
    },
    { # shard 2 (single zone)
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]
}
