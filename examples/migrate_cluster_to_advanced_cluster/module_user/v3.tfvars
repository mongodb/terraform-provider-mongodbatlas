project_id             = "664619d870c247237f4b86a6"
cluster_name           = "module-cluster"
cluster_type           = "SHARDED"
mongo_db_major_version = "8.0"

tags = {
  env    = "examples"
  module = "cluster_to_advanced_cluster"
}
# removed variables: replication_specs, provider_name, instance_size, disk_size
replication_specs_new = [
  { # shard 1
    region_configs = [{
      electable_specs = {
        disk_size_gb  = 50    # must be the same for all replication specs
        instance_size = "M30" # changed
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      read_only_specs = { # new
        disk_size_gb  = 50
        instance_size = "M30"
        node_count    = 1
      }
      region_name = "US_EAST_1"
    }]
    zone_name = "Zone 1"
  },
  { # shard 2
    region_configs = [{
      electable_specs = {
        disk_size_gb  = 50
        instance_size = "M50" # increased from M30
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      read_only_specs = { # new
        disk_size_gb  = 50
        instance_size = "M50"
        node_count    = 1
      }
      region_name = "US_EAST_1"
    }]
    zone_name = "Zone 1"
  },
]

