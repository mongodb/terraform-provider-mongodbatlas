module "cluster_skip_python" {
  source = "../.."


  replication_specs = [{
    region_configs = [{
      provider_name = "AWS",
      region_name   = "EU_WEST_1",
      priority      = 7,
      electable_specs = {
        node_count    = 3
        instance_size = "M10"
        disk_size_gb  = 10
      }
    }]
  }]

  project_id   = "your-project-id"
  name         = "created-from-resource-module"
  cluster_type = "REPLICASET"


}