# Set up cloud provider access in Atlas for AWS
resource "mongodbatlas_cloud_provider_access_setup" "setup" {
  project_id    = mongodbatlas_project.project.id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth" {
  project_id = mongodbatlas_project.project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.atlas_role.arn
  }
}

resource "mongodbatlas_log_integration" "example" {
  project_id  = mongodbatlas_project.project.id
  type        = "S3_LOG_EXPORT"
  log_types   = ["MONGOD_AUDIT"]
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth.role_id
  prefix_path = "atlas-logs"
}
