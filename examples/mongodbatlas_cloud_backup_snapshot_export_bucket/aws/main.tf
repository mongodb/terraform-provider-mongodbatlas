resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}


resource "aws_s3_bucket" "test_bucket" {
  bucket = "mongo-test-bucket-1"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}

resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id = var.project_id

  iam_role_id    = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  bucket_name    = aws_s3_bucket.test_bucket.bucket
  cloud_provider = "AWS"
}

