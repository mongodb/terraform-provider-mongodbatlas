# - single region
# - no sharding
# - with analytics nodes defined

module "single_region_analytics" {
  source = "../.."

  project_id             = var.project_id
  name                   = "single-region-analytics"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"

  region_configs = [
    {
      provider_name        = "AWS"
      region_name          = "US_EAST_1"
      electable_node_count = 3
      analytics_specs = {
        node_count    = 1
      }
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