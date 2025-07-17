# - single region
# - with shards (single zone)

module "single_region_sharded" {
  source = "./cluster-abstraction"

  project_id             = var.project_id
  name                   = "single-region-sharded"
  cluster_type           = "SHARDED"
  mongo_db_major_version = "8.0"

  shards = [
    { # shard 1 (single zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "US_EAST_1"
          instance_size        = "M40" # Independently scaled shard
          electable_node_count = 3
        }
      ]
    },
    { # shard 2 (single zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "US_EAST_1"
          instance_size        = "M30"
          electable_node_count = 3
        }
      ]
    }
  ]

  auto_scaling = {
    disk_gb_enabled           = true
    compute_enabled           = true
    compute_max_instance_size = "M60"
    compute_min_instance_size = "M30"
  }
}