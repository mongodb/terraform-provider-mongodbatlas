module "multi_region_single_geo_no_sharding" {
  source                 = "../.."
  project_id             = var.project_id
  name                   = "multi-region-single-geo"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"
  replication_specs = [
    {
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          priority      = 7
          electable_specs = {
            node_count = 2
          }
          auto_scaling = {
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        },
        {
          provider_name = "AWS"
          region_name   = "US_EAST_2"
          priority      = 6
          electable_specs = {
            node_count = 1
          }
          read_only_specs = {
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
