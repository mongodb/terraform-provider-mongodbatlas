resource "mongodbatlas_cluster" "geosharded" {
  project_id                   = var.project_id
  name                         = var.cluster_name
  cloud_backup                 = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "5.0"
  cluster_type                 = "GEOSHARDED"
  replication_specs {
    zone_name  = "Zone 1"
    num_shards = 2
    regions_config {
      region_name     = var.atlas_region_east
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
    regions_config {
      region_name     = var.atlas_region_west
      electable_nodes = 2
      priority        = 6
      read_only_nodes = 0
    }
  }

  # Provider settings
  provider_name               = "AWS"
  disk_size_gb                = 80
  provider_instance_size_name = "M30"

  depends_on = [
    mongodbatlas_privatelink_endpoint_service.pe_east_service,
    mongodbatlas_privatelink_endpoint_service.pe_west_service,
    mongodbatlas_private_endpoint_regional_mode.test
  ]
}

