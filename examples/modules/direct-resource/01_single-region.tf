# - single region
# - no sharding

resource "mongodbatlas_advanced_cluster" "single_region" {
  project_id             = var.project_id
  name                   = "single-region"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"
  replication_specs = [
    {
      region_configs = [
        {
          auto_scaling = {
            disk_gb_enabled           = true
            compute_enabled           = true
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          priority      = 7
          provider_name = "AWS"
          region_name   = "US_EAST_1"
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
      replication_specs[0].region_configs[0].electable_specs.instance_size
    ]
  }
}

