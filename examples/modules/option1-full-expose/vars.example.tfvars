
backup_enabled = true
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

project_id   = "664619d870c247237f4b86a6"
name         = "created-from-option1-module"
cluster_type = "REPLICASET"

