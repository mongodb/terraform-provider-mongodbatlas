# Outputs for the new log integration (always available)
output "log_integration_id" {
  description = "The ID of the log integration"
  value       = mongodbatlas_log_integration.this.integration_id
}

output "new_prefix_path" {
  description = "The prefix path for the new log integration"
  value       = mongodbatlas_log_integration.this.prefix_path
}

# Outputs for the old push-based log export (only available when not skipped)
output "old_prefix_path" {
  description = "The prefix path for the old push-based log export (null if skipped)"
  value       = var.skip_push_based_log_export ? null : mongodbatlas_push_based_log_export.this[0].prefix_path
}

output "old_state" {
  description = "The state of the old push-based log export (null if skipped)"
  value       = var.skip_push_based_log_export ? null : mongodbatlas_push_based_log_export.this[0].state
}

