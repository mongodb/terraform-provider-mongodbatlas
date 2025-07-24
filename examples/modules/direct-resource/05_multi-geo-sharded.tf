# - multiple regions (different geographies) 
# - with shards (single zone)

resource "mongodbatlas_advanced_cluster" "multi_geo_sharded" {
  project_id             = var.project_id
  name                   = "multi-geo-sharded"
  cluster_type           = "SHARDED"
  mongo_db_major_version = "8.0"
  replication_specs = [
    { # shard 1 (single zone)
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1" # North America
          priority      = 7
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          auto_scaling = {
            disk_gb_enabled           = true
            compute_enabled           = true
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        },
        {
          provider_name = "AWS"
          region_name   = "EU_WEST_1" # Europe
          priority      = 6
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          auto_scaling = {
            disk_gb_enabled           = true
            compute_enabled           = true
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        }
      ]
    },
    { # shard 2 (single zone)
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1" # North America
          priority      = 7
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          auto_scaling = {
            disk_gb_enabled           = true
            compute_enabled           = true
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        },
        {
          provider_name = "AWS"
          region_name   = "EU_WEST_1" # Europe
          priority      = 6
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          auto_scaling = {
            disk_gb_enabled           = true
            compute_enabled           = true
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        }
      ]
    }
  ]

  tags = {
    department       = "Engineering"
    team_name        = "APIx Integrations"
    application_name = "Telemetry"
    environment      = "prod"
    version          = "1.0"
    email_contact    = "agustin.bettati@mongodb.com"
    criticality      = "Tier 1 with PII"
  }

  lifecycle {
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      replication_specs[0].region_configs[1].electable_specs.instance_size,
      replication_specs[0].region_configs[0].read_only_specs.instance_size,
      replication_specs[0].region_configs[1].read_only_specs.instance_size
    ]
  }
}
