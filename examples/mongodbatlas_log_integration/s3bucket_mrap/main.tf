# Set up cloud provider access in Atlas
resource "mongodbatlas_cloud_provider_access_setup" "setup" {
  project_id    = var.project_id
  provider_name = "AWS"
}

# Authorize the IAM role for Atlas access
resource "mongodbatlas_cloud_provider_access_authorization" "auth" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.atlas_role.arn
  }
}

# Set up log integration with MRAP
resource "mongodbatlas_log_integration" "mrap_s3" {
  project_id  = var.project_id
  type        = "S3_LOG_EXPORT"
  bucket_name = aws_s3control_multi_region_access_point.atlas_logs.arn # MRAP alias
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth.role_id
  prefix_path = var.prefix_path
  log_types   = var.log_types
}

output "mrap_alias" {
  description = "The MRAP alias used in bucket_name"
  value       = aws_s3control_multi_region_access_point.atlas_logs.alias
}

output "mrap_arn" {
  description = "The MRAP ARN"
  value       = aws_s3control_multi_region_access_point.atlas_logs.arn
}

output "integration_id" {
  description = "The ID of the log integration"
  value       = mongodbatlas_log_integration.mrap_s3.integration_id
}

output "iam_role_arn" {
  description = "ARN of the IAM role used for the integration"
  value       = aws_iam_role.atlas_role.arn
}

output "cloud_provider_role_id" {
  description = "Atlas Cloud Provider Access Role ID"
  value       = mongodbatlas_cloud_provider_access_authorization.auth.role_id
}

output "backing_buckets" {
  description = "The S3 buckets backing the MRAP"
  value = {
    us_east_1 = aws_s3_bucket.mrap_bucket_us_east_1.bucket
    us_west_2 = aws_s3_bucket.mrap_bucket_us_west_2.bucket
  }
}
