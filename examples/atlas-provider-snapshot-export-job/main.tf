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

resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = var.project_id
  name         = "MyCluster"
  disk_size_gb = 1

  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = var.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = "myDescription"
  retention_in_days = 1
}

resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id = var.project_id

  iam_role_id    = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  bucket_name    = aws_s3_bucket.test_bucket.bucket
  cloud_provider = "AWS"
}

resource "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
  project_id       = var.project_id
  cluster_name     = mongodbatlas_cluster.my_cluster.name
  snapshot_id      = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
  export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id


  custom_data {
    key   = "exported by"
    value = "myName"
  }
}
