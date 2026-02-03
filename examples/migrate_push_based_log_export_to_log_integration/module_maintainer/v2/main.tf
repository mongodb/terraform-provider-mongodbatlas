# v2: Updated module supporting migration with feature flag
# - mongodbatlas_log_integration is always created
# - mongodbatlas_push_based_log_export is conditionally created based on skip_push_based_log_export

# Old resource - conditionally created based on feature flag
resource "mongodbatlas_push_based_log_export" "this" {
  count = var.skip_push_based_log_export ? 0 : 1

  project_id  = var.project_id
  bucket_name = var.bucket_name
  iam_role_id = var.iam_role_id
  prefix_path = var.prefix_path
}

# New resource - always created
resource "mongodbatlas_log_integration" "this" {
  project_id  = var.project_id
  bucket_name = var.bucket_name
  iam_role_id = var.iam_role_id
  prefix_path = var.new_prefix_path
  type        = "S3_LOG_EXPORT"
  log_types   = var.log_types
}

