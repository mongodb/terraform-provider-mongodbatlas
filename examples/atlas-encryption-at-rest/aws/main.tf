resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = var.cloud_provider_access_name
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.project_id

  aws_kms_config {
    enabled                = true
    customer_master_key_id = var.customer_master_key
    region                 = var.atlas_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }

}

