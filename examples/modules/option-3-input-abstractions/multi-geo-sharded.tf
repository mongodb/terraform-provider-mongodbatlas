# - multiple regions (different geographies) 
# - with shards (single zone)

module "multi_geo_sharded" {
  source = "./cluster-abstraction"

  project_id = var.project_id
  name       = "multi_geo_sharded"
  cluster_type = "SHARDED"
  mongo_db_major_version = "8.0"

  shards = [ 
    { # shard 1 (single zone)
      region_configs = [
        {
          provider_name  = "AWS"
          region_name    = "US_EAST_1" # North America
          instance_size  = "M30"
          electable_node_count = 3
          priority = 7
        },
        {
          provider_name  = "AWS"
          region_name    = "EU_WEST_1" # Europe
          instance_size  = "M30"
          electable_node_count = 2
          priority = 6
        }
      ]
    },
    { # shard 2 (single zone)
      region_configs = [
        {
          provider_name  = "AWS"
          region_name    = "US_EAST_1" # North America
          instance_size  = "M30"
          electable_node_count = 3
          priority = 7
        },
        {
          provider_name  = "AWS"
          region_name    = "EU_WEST_1" # Europe
          instance_size  = "M30"
          electable_node_count = 2
          priority = 6
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