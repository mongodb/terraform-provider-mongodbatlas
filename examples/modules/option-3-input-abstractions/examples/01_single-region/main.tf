# Example usage of the cluster-abstraction module for a single-region cluster

module "single_region" {
  source = "../.."

  project_id             = var.project_id
  name                   = "single-region"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"

  region_configs = [
    {
      provider_name        = "AWS"
      region_name          = "US_EAST_1"
      instance_size        = "M30"
      electable_node_count = 3
    }
  ]

  auto_scaling = {
    disk_gb_enabled           = true
    compute_enabled           = true
    compute_max_instance_size = "M30"
    compute_min_instance_size = "M60"
  }
}
