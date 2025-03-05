# resource "mongodbatlas_cluster" "this" {
#   project_id                  = var.project_id
#   name                        = var.cluster_name
#   cluster_type                = "REPLICASET"
#   provider_name               = "AWS"
#   provider_instance_size_name = var.instance_size
#   mongo_db_major_version      = var.mongo_db_major_version
#   disk_size_gb                = 30

#   advanced_configuration {
#     javascript_enabled = true
#   }
#   tags {
#     key   = "ManagedBy"
#     value = "Terraform"
#   }
#   tags {
#     key   = "Example"
#     value = "examples-migrate_cluster_to_advanced_cluster-basic"
#   }
#   replication_specs {
#     num_shards = 1
#     regions_config {
#       region_name     = "US_WEST_1"
#       electable_nodes = 2
#       priority        = 7
#     }
#     regions_config {
#       region_name     = "US_EAST_1"
#       electable_nodes = 3
#       analytics_nodes = 1
#       read_only_nodes = 1
#       priority        = 6
#     }
#   }
# }
