# v3: After migration - only log_integration remains
# The old push_based_log_export has been removed

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

# Final log integration configuration
# You can now use the original prefix path if desired
resource "mongodbatlas_log_integration" "logs" {
  project_id  = mongodbatlas_project.project.id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "atlas-logs" # Can revert to original path after old resource is removed
  type        = "S3_LOG_EXPORT"
  log_types   = ["MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"]
}

output "log_prefix" {
  description = "Prefix path for the log integration"
  value       = mongodbatlas_log_integration.logs.prefix_path
}

output "log_integration_id" {
  description = "ID of the log integration"
  value       = mongodbatlas_log_integration.logs.integration_id
}

