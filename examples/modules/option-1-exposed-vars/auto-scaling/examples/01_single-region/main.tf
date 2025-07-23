module "single_region" {
  source                 = "../.."
  project_id             = var.project_id
  name                   = "single-region"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"
  replication_specs = [
    {
      region_configs = [
        {
          auto_scaling = {
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
          electable_specs = {
            node_count = 3
          }
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]
}
