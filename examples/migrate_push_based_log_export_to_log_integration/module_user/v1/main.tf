# v1: Module user - Initial upgrade with both resources active

# Using the updated module (v2) with the feature flag set to false
# This keeps both resources active during the validation period
module "log_export" {
  source = "../../module_maintainer/v2"

  project_id      = var.project_id
  bucket_name     = var.bucket_name
  iam_role_id     = var.iam_role_id
  prefix_path     = "atlas-logs"     # Old resource path
  new_prefix_path = "atlas-logs-new" # New resource path (distinct during migration)
  log_types       = ["MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"]

  # Keep old resource active during migration
  skip_push_based_log_export = false
}

output "old_prefix_path" {
  description = "Prefix path for the old push-based log export"
  value       = module.log_export.old_prefix_path
}

output "new_prefix_path" {
  description = "Prefix path for the new log integration"
  value       = module.log_export.new_prefix_path
}

output "log_integration_id" {
  description = "ID of the new log integration"
  value       = module.log_export.log_integration_id
}

