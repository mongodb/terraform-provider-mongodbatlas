resource "mongodbatlas_advanced_cluster" "geosharded" {
  project_id     = var.project_id
  name           = var.cluster_name
  cluster_type   = "GEOSHARDED"
  backup_enabled = true

  replication_specs = [
    { # Shard 1
      zone_name = "Zone 1"

      region_configs = [{
        electable_specs = {
          instance_size = "M30"
          node_count    = 3
        }
        provider_name = "AWS"
        priority      = 7
        region_name   = var.atlas_region_east
        },
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AWS"
          priority      = 6
          region_name   = var.atlas_region_west
      }]
    },
    { # Shard 2
      zone_name = "Zone 1"

      region_configs = [{
        electable_specs = {
          instance_size = "M30"
          node_count    = 3
        }
        provider_name = "AWS"
        priority      = 7
        region_name   = var.atlas_region_east
        },
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AWS"
          priority      = 6
          region_name   = var.atlas_region_west
      }]
    }
  ]

  depends_on = [
    mongodbatlas_privatelink_endpoint_service.pe_east_service,
    mongodbatlas_privatelink_endpoint_service.pe_west_service,
    mongodbatlas_private_endpoint_regional_mode.test
  ]
}
