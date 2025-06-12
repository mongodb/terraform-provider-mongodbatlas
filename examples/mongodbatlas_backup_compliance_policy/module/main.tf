data "mongodbatlas_atlas_user" "this" {
  user_id = var.user_id
}

resource "mongodbatlas_backup_compliance_policy" "this" {
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

module "cluster_with_schedule" {
  source = "./modules/cluster_with_schedule"

  project_id    = var.project_id
  instance_size = var.instance_size
  cluster_name  = var.cluster_name
  add_schedule  = true # change to false in Step 2
}

# Step 2: For removing the `mongodbatlas_cloud_backup_schedule` resource
# Rename the resource to avoid the `Removed Resource still exists error`
# moved {
#   from = module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this[0] # must be deleted with the `add_schedule` variable set to false
#   to   = mongodbatlas_cloud_backup_schedule.to_be_deleted                              # any resource name that doesn't exist works!
# }

# removed {
#   from = mongodbatlas_cloud_backup_schedule.to_be_deleted # any resource name that doesn't exist works!

#   lifecycle {
#     destroy = false
#   }
# }
