resource "mongodbatlas_cluster" "cluster-atlas" {
  project_id                   = var.atlasprojectid
  name                         = "cluster-atlas"
  cloud_backup                 = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "4.2"
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
  }

  replication_specs {
    zone_name  = "Zone 2"
    num_shards = 2
    regions_config {
      region_name     = var.atlas_region_west
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  # Provider settings
  provider_name               = "AWS"
  disk_size_gb                = 80
  provider_instance_size_name = "M30"
}

data "mongodbatlas_cluster" "cluster-atlas" {
  project_id = var.atlasprojectid
  name       = mongodbatlas_cluster.cluster-atlas.name
  depends_on = [
    mongodbatlas_privatelink_endpoint_service.atlaseplink_east,
    mongodbatlas_privatelink_endpoint_service.atlaseplink_west,
    mongodbatlas_private_endpoint_regional_mode.test
  ]
}

output "atlasclusterstring" {
  value = data.mongodbatlas_cluster.cluster-atlas.connection_strings
}
