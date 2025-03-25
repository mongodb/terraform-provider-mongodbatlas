resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "111111111111111111111111"
  name         = "mocked-cluster"
  cluster_type = "REPLICASET"

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        disk_size_gb  = 10
        instance_size = "M10"
        node_count    = 5
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
  }]
  timeouts = {
    create = "6000s"
  }
}
