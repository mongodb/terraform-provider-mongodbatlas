# - multiple regions (different geographies) 
# - sharded zones

module "multi_geo_zone_sharded" {
  source = "../.."

  project_id             = var.project_id
  name                   = "multi-geo-zone-sharded"
  cluster_type           = "GEOSHARDED"
  mongo_db_major_version = "8.0"

  shards = [
    {
      zone_name = "US" # shard 1 (US zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "US_EAST_1"
          electable_node_count = 3
        }
      ]
    },
    {
      zone_name = "EU" # shard 2 (EU zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "EU_WEST_1"
          electable_node_count = 3
        }
      ]
    }
  ]

  auto_scaling = {
    compute_max_instance_size = "M60"
    compute_min_instance_size = "M30"
  }
  analytics_auto_scaling = {
    compute_max_instance_size = "M30"
    compute_min_instance_size = "M10"
  }
}