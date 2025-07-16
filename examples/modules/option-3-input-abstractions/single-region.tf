# Example usage of the cluster-abstraction module for a single-region cluster

module "single_region" {
  source = "./cluster-abstraction"

  project_id = var.project_id
  name       = "single-region"
  cluster_type = "REPLICASET"
  mongo_db_major_version = "8.0"

  replica_set_regions = [
    {
      provider_name  = "AWS"
      region_name    = "US_EAST_1"
      instance_size  = "M10"
      electable_node_count = 3
      read_only_node_count = 2
    }
  ]

  auto_scaling = {
    disk_gb_enabled           = true
    compute_enabled           = true
    compute_max_instance_size = "M20"
    compute_min_instance_size = "M10"
  }
}

output "cluster_id" {
  value = module.single_region.cluster_id
}

output "connection_strings" {
  value = module.single_region.connection_strings
}
