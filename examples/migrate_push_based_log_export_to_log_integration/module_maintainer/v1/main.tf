# v1: Original module using mongodbatlas_push_based_log_export

resource "mongodbatlas_push_based_log_export" "this" {
  project_id  = var.project_id
  bucket_name = var.bucket_name
  iam_role_id = var.iam_role_id
  prefix_path = var.prefix_path
}

