module "cluster" {
  source                 = "../.."
  project_id             = var.project_id
  name                   = "single-region-analytics"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"
  replication_specs = [
    {
      region_configs = [
        {
          auto_scaling = {
            disk_gb_enabled           = true
            compute_enabled           = true
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
          analytics_auto_scaling = {
            disk_gb_enabled           = true
            compute_enabled           = true
            compute_max_instance_size = "M30"
            compute_min_instance_size = "M10"
          }
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          analytics_specs = {
            instance_size = "M10"
            node_count    = 1
          }
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]
}