module "multi_geo_sharded" {
  source                 = "../.."
  project_id             = var.project_id
  name                   = "multi-geo-sharded"
  cluster_type           = "SHARDED"
  mongo_db_major_version = "8.0"
  replication_specs = [
    { # shard 1 (single zone)
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
        },
        {
          provider_name = "AWS"
          region_name   = "EU_WEST_1" # Europe
          priority      = 6
          electable_specs = {
            node_count = 2
          }
          auto_scaling = {
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        }
      ]
    },
    { # shard 2 (single zone)
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
        },
        {
          provider_name = "AWS"
          region_name   = "EU_WEST_1" # Europe
          priority      = 6
          electable_specs = {
            node_count = 2
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
