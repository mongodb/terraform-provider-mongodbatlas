resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "111111111111111111111111"
  name         = "mocked-cluster"
  cluster_type = "GEOSHARDED"


  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 5
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "Zone 1"
    }, {
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
    zone_name = "Zone 2"
  }]
}
