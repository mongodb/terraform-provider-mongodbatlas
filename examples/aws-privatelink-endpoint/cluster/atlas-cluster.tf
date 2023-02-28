resource "mongodbatlas_cluster" "aws_private_connection" {
  project_id                   = var.project_id
  name                         = var.cluster_name
  cloud_backup                 = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "4.2"
  cluster_type                 = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  # Provider settings
  provider_name               = "AWS"
  disk_size_gb                = 10
  provider_instance_size_name = "M10"

  depends_on = [mongodbatlas_privatelink_endpoint_service.pe_east_service]
}