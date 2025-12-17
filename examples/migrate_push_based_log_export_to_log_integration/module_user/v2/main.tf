# v2: Module user - Migration complete, old resource removed

# Using the updated module (v2) with the feature flag set to true
# This removes the old push_based_log_export resource
module "log_export" {
  source = "../../module_maintainer/v2"

  project_id      = var.project_id
  bucket_name     = var.bucket_name
  iam_role_id     = var.iam_role_id
  prefix_path     = "atlas-logs" # No longer used (old resource removed)
  new_prefix_path = "atlas-logs" # Can use original path now
  log_types       = ["MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"]

  # Remove old resource - migration complete
  skip_push_based_log_export = true
}

output "prefix_path" {
  description = "Prefix path for the log integration"
  value       = module.log_export.new_prefix_path
}

output "log_integration_id" {
  description = "ID of the log integration"
  value       = module.log_export.log_integration_id
}

