# v1: Original configuration with mongodbatlas_push_based_log_export
# This represents the starting point before migration

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

# Original push-based log export configuration
resource "mongodbatlas_push_based_log_export" "logs" {
  project_id  = mongodbatlas_project.project.id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "atlas-logs"
}

output "log_prefix" {
  value = mongodbatlas_push_based_log_export.logs.prefix_path
}

output "log_state" {
  value = mongodbatlas_push_based_log_export.logs.state
}

