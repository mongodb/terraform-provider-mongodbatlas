data "mongodbatlas_atlas_user" "this" {
  user_id = var.user_id
}

resource "mongodbatlas_backup_compliance_policy" "backup_policy" {
  project_id                 = var.project_id
  authorized_email           = data.mongodbatlas_atlas_user.this.email_address
  authorized_user_first_name = data.mongodbatlas_atlas_user.this.first_name
  authorized_user_last_name  = data.mongodbatlas_atlas_user.this.last_name
  copy_protection_enabled    = false
  pit_enabled                = false
  encryption_at_rest_enabled = false

  restore_window_days = 7
  on_demand_policy_item {
    frequency_interval = 0
    retention_unit     = "days"
    retention_value    = 1
  }
  policy_item_daily {
    frequency_interval = 0
    retention_unit     = "days"
    retention_value    = 1
  }
}

resource "mongodbatlas_advanced_cluster" "this" {
  project_id             = var.project_id
  name                   = var.cluster_name
  cluster_type           = "REPLICASET"
  backup_enabled         = true
  retain_backups_enabled = true

  replication_specs = [
    {
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          priority      = 7
          electable_specs = {
            node_count    = 3
            instance_size = var.instance_size
            disk_size_gb  = 10
          }
        },
      ]
    }
  ]
}

resource "mongodbatlas_cloud_backup_schedule" "this" {
  project_id               = var.project_id
  cluster_name             = resource.mongodbatlas_advanced_cluster.this.name
  reference_hour_of_day    = 19
  reference_minute_of_hour = 15
  restore_window_days      = 7
  policy_item_daily {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 1
  }
  copy_settings {
    cloud_provider = "AWS"
    frequencies = [
      "HOURLY",
      "DAILY",
      "WEEKLY",
      "MONTHLY",
      "YEARLY",
      "ON_DEMAND",
    ]
    region_name        = "US_EAST_2"
    zone_id            = mongodbatlas_advanced_cluster.this.replication_specs[0].zone_id
    should_copy_oplogs = false
  }
}
