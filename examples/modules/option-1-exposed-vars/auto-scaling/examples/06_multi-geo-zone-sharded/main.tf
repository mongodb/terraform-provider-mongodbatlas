module "multi_geo_zone_sharded" {
  source                 = "../.."
  project_id             = var.project_id
  name                   = "multi-geo-zone-sharded"
  cluster_type           = "GEOSHARDED"
  mongo_db_major_version = "8.0"
  replication_specs = [
    { # shard 1 (US zone)
      zone_name = "US"
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1" # North America
          priority      = 7
          electable_specs = {
            node_count = 3
          }
          auto_scaling = {
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        }
      ]
    },
    { # shard 2 (EU zone)
      zone_name = "EU"
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "EU_WEST_1" # Europe
          priority      = 7
          electable_specs = {
            node_count = 3
          }
          auto_scaling = {
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        }
      ]
    }
  ]
}
