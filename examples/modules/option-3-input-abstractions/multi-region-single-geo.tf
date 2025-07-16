# 5-Node 2-Region Architecture

# - multiple regions (same geography)
# no sharding

module "cluster_abstraction" {
  source = "./cluster-abstraction"

  project_id = var.project_id
  name       = "multi-region-single-geo-no-sharding"
  cluster_type = "REPLICASET"
  mongo_db_major_version = "8.0"

  replica_set_regions = [
    {
      provider_name  = "AWS"
      region_name    = "US_EAST_1"
      instance_size  = "M30"
      electable_node_count = 2
      priority       = 7
    },
    {
      provider_name  = "AWS"
      region_name    = "US_EAST_2"
      instance_size  = "M30"
      electable_node_count = 1
      read_only_node_count = 2
      priority       = 6
    },
  ]

  auto_scaling = {
    disk_gb_enabled           = true
    compute_enabled           = true
    compute_max_instance_size = "M60"
    compute_min_instance_size = "M30"
  }
}