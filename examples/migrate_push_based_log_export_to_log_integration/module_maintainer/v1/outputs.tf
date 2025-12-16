output "prefix_path" {
  description = "The prefix path where logs are stored"
  value       = mongodbatlas_push_based_log_export.this.prefix_path
}

output "state" {
  description = "The state of the push-based log export"
  value       = mongodbatlas_push_based_log_export.this.state
}

