project_id             = "664619d870c247237f4b86a6"
cluster_name           = "module-cluster"
cluster_type           = "SHARDED"
instance_size          = "M10"
mongo_db_major_version = "8.0"
provider_name          = "AWS"
disk_size              = 40
tags = {
  env    = "examples"
  module = "cluster_to_advanced_cluster"
}
replication_specs = [
  {
    num_shards = 2
    zone_name  = "Zone 1"
    regions_config = [
      {
        read_only_nodes = 0
        priority        = 7
        region_name     = "US_EAST_1"
        electable_nodes = 3
      },
      {
        read_only_nodes = 0
        priority        = 6
        region_name     = "EU_WEST_1"
        electable_nodes = 2
      },
    ]
  }
]
search_nodes_specs = [
  {
    instance_size = "S20_HIGHCPU_NVME"
    node_count    = 3
  }
]
encryption_at_rest_provider = "AWS"
