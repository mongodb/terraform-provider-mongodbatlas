terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.34"
    }
  }
  required_version = ">= 1.0"
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
  count                    = var.add_schedule ? 1 : 0
  project_id               = var.project_id
  cluster_name             = mongodbatlas_advanced_cluster.this.name
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
    frequencies = ["HOURLY",
      "DAILY",
      "WEEKLY",
      "MONTHLY",
      "YEARLY",
    "ON_DEMAND"]
    region_name        = "US_EAST_2"
    zone_id            = mongodbatlas_advanced_cluster.this.replication_specs[0].zone_id
    should_copy_oplogs = false
  }
}