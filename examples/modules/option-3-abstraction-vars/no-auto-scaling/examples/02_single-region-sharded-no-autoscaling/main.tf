# - single region
# - with shards (single zone)

module "single_region_sharded_no_autoscaling" {
  source = "../.."

  project_id             = var.project_id
  name                   = "single-region-sharded-no-autoscaling"
  cluster_type           = "SHARDED"
  mongo_db_major_version = "8.0"

  shards = [
    { # shard 1 (single zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "US_EAST_1"
          instance_size        = "M40" # Independently scaled shard
          electable_node_count = 3
        }
      ]
    },
    { # shard 2 (single zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "US_EAST_1"
          instance_size        = "M30"
          electable_node_count = 3
        }
      ]
    }
  ]

  tags_recommended = { # defined keys are enforced through validations
    department       = "Engineering"
    team_name        = "APIx Integrations"
    application_name = "Telemetry"
    environment      = "prod"
    version          = "1.0"
    email_contact    = "agustin.bettati@mongodb.com"
    criticality      = "Tier 1 with PII"
  }

}
