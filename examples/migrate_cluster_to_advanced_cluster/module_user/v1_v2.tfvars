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
        region_name     = "US_EAST_1"
        electable_nodes = 3
        priority        = 7
        read_only_nodes = 0
      }
    ]
  }
]