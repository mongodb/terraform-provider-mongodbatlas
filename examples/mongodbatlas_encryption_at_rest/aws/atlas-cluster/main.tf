resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.atlas_project_id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.atlas_project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  aws_kms_config {
    enabled                = true
    customer_master_key_id = aws_kms_key.kms_key.id
    region                 = var.atlas_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }
}

resource "mongodbatlas_cluster" "cluster" {
  project_id                  = var.atlas_project_id
  name                        = "MyCluster"
  cluster_type                = "REPLICASET"
  provider_name               = "AWS"
  encryption_at_rest_provider = "AWS"
  backing_provider_name       = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
}
