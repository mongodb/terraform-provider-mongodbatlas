# v2: During migration - both resources are active
# The new log_integration uses a distinct prefix path to avoid conflicts

resource "mongodbatlas_project" "project" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}

# Set up cloud provider access in Atlas using the created IAM role
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.project.id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.atlas_role.arn
  }
}

# KEEP: Original push-based log export (will be removed in v3)
resource "mongodbatlas_push_based_log_export" "logs" {
  project_id  = mongodbatlas_project.project.id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "atlas-logs"
}

# NEW: Log integration resource with distinct prefix path
resource "mongodbatlas_log_integration" "logs" {
  project_id  = mongodbatlas_project.project.id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "atlas-logs-new" # Use distinct path during migration
  type        = "S3_LOG_EXPORT"
  log_types   = ["MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"]
}

# Outputs for both resources to compare
output "old_log_prefix" {
  description = "Prefix path for the old push-based log export"
  value       = mongodbatlas_push_based_log_export.logs.prefix_path
}

output "new_log_prefix" {
  description = "Prefix path for the new log integration"
  value       = mongodbatlas_log_integration.logs.prefix_path
}

output "new_log_integration_id" {
  description = "ID of the new log integration"
  value       = mongodbatlas_log_integration.logs.integration_id
}

