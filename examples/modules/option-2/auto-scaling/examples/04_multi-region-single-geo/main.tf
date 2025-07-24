# 5-Node 2-Region Architecture

# - multiple regions (same geography)
# no sharding

# using abstraction variable
module "abstraction_variables_multi_region_single_geo_no_sharding" {
  source = "../.."

  project_id             = var.project_id
  name                   = "multi-region-single-geo-no-sharding"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"

  region_configs = [
    {
      provider_name        = "AWS"
      region_name          = "US_EAST_1"
      electable_node_count = 2
      priority             = 7
    },
    {
      provider_name        = "AWS"
      region_name          = "US_EAST_2"
      electable_node_count = 1
      read_only_node_count = 2
      priority             = 6
    },
  ]

  auto_scaling = {
    compute_max_instance_size = "M60"
    compute_min_instance_size = "M30"
  }

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

# using direct variables (replication_specs)
module "direct_variables_multi_region_single_geo_no_sharding" {
  source                 = "../.."
  project_id             = var.project_id
  name                   = "multi-region-single-geo"
  cluster_type           = "REPLICASET"
  mongo_db_major_version = "8.0"
  replication_specs = [
    {
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          priority      = 7
          electable_specs = {
            node_count = 2
          }
          auto_scaling = {
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        },
        {
          provider_name = "AWS"
          region_name   = "US_EAST_2"
          priority      = 6
          electable_specs = {
            node_count = 1
          }
          read_only_specs = {
            node_count = 2
          }
          auto_scaling = {
            compute_max_instance_size = "M60"
            compute_min_instance_size = "M30"
          }
        }
      ]
    }
  ]

  tags = { # has flexibility to ignore defining certain tags
    department       = "Engineering"
    team_name        = "APIx Integrations"
    application_name = "Telemetry"
    environment      = "prod"
    version          = "1.0"
    email_contact    = "agustin.bettati@mongodb.com"
    criticality      = "Tier 1 with PII"
  }
}
