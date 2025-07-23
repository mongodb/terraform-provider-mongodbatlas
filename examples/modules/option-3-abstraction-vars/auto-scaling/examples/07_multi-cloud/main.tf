# multiple providers, multi-region single geo
# multi-cloud can also apply to multi-region multi geo, with sharding.

module "multi_cloud" {
  source = "../.."

  project_id             = var.project_id
  name                   = "multi-cloud"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"

  region_configs = [
    {
      provider_name        = "AZURE"
      region_name          = "US_WEST_2"
      electable_node_count = 2
      priority             = 7
    },
    {
      provider_name        = "AWS"
      region_name          = "US_EAST_2"
      electable_node_count = 1
      read_only_node_count = 2
      priority             = 6
    },
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