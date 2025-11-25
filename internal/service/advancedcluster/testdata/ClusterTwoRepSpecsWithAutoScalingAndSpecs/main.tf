resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "111111111111111111111111"
  name         = "mocked-cluster"
  cluster_type = "GEOSHARDED"


  replication_specs = [{
    region_configs = [{
      analytics_auto_scaling = {
        compute_enabled            = true
        compute_max_instance_size  = "M30"
        compute_min_instance_size  = "M10"
        compute_scale_down_enabled = true
        disk_gb_enabled            = true
      }
      auto_scaling = {
        compute_enabled            = true
        compute_max_instance_size  = "M30"
        compute_min_instance_size  = "M10"
        compute_scale_down_enabled = true
        disk_gb_enabled            = true
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 5
      }
      priority      = 7
      provider_name = "AWS"
      read_only_specs = {
        instance_size = "M10"
        node_count    = 2
      }
      region_name = "US_EAST_1"
    }]
    zone_name = "Zone 1"
    }, {
    region_configs = [{
      analytics_auto_scaling = {
        compute_enabled            = true
        compute_max_instance_size  = "M30"
        compute_min_instance_size  = "M10"
        compute_scale_down_enabled = true
        disk_gb_enabled            = true
      }
      analytics_specs = {
        instance_size = "M10"
        node_count    = 4
      }
      auto_scaling = {
        compute_enabled            = true
        compute_max_instance_size  = "M30"
        compute_min_instance_size  = "M10"
        compute_scale_down_enabled = true
        disk_gb_enabled            = true
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      read_only_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      region_name = "US_WEST_2"
    }]
    zone_name = "Zone 2"
  }]
}
