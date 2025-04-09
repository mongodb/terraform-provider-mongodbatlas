resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "111111111111111111111111"
  name         = "mocked-cluster"
  cluster_type = "REPLICASET"

  replication_specs = [{
    region_configs = [
      {
      auto_scaling = {
        compute_enabled            = false
        compute_scale_down_enabled = false
        disk_gb_enabled            = true
      }
      electable_specs = {
        disk_size_gb  = 10
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    },
    {
      auto_scaling = {
        compute_enabled            = false
        compute_scale_down_enabled = false
        disk_gb_enabled            = true
      }
      electable_specs = {
        disk_size_gb  = 10
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_2"
    }
    ]
  }]
  timeouts = {
    create = "6000s"
  }
}
