# - multiple regions (different geographies) 
# - with shards (single zone)

module "multi_geo_sharded" {
  source = "../.."

  project_id             = var.project_id
  name                   = "multi-geo-sharded"
  cluster_type           = "SHARDED"
  mongo_db_major_version = "8.0"

  shards = [
    { # shard 1 (single zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "US_EAST_1" # North America
          electable_node_count = 3
          priority             = 7
        },
        {
          provider_name        = "AWS"
          region_name          = "EU_WEST_1" # Europe
          electable_node_count = 2
          priority             = 6
        }
      ]
    },
    { # shard 2 (single zone)
      region_configs = [
        {
          provider_name        = "AWS"
          region_name          = "US_EAST_1" # North America
          electable_node_count = 3
          priority             = 7
        },
        {
          provider_name        = "AWS"
          region_name          = "EU_WEST_1" # Europe
          electable_node_count = 2
          priority             = 6
        }
      ]
    }
  ]

  auto_scaling = {
    compute_max_instance_size = "M60"
    compute_min_instance_size = "M30"
  }
  analytics_auto_scaling = {
    compute_max_instance_size = "M30"
    compute_min_instance_size = "M10"
  }

  tags = { # defined keys are enforced through validations
    department    = "Engineering"
    team_name     = "APIx Integrations"
    application_name = "Telemetry"
    environment   = "prod"
    version       = "1.0"
    email_contact = "agustin.bettati@mongodb.com"
    criticality   = "Tier 1 with PII"
  }
}