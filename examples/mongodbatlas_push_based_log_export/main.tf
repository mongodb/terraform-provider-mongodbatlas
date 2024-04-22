resource "mongodbatlas_project" "project-tf" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}

# Set up cloud provider access in Atlas using the created IAM role
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.project-tf.id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.project-tf.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

# Set up push-based log export with authorized IAM role
resource "mongodbatlas_push_based_log_export" "test" {
  project_id  = mongodbatlas_project.project-tf.id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "push-based-log-test"
}

data "mongodbatlas_push_based_log_export" "test" {
  project_id = mongodbatlas_push_based_log_export.test.project_id
}

output "test" {
   value = data.mongodbatlas_push_based_log_export.test.prefix_path
}
