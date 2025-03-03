resource "mongodbatlas_encryption_at_rest" "ear" {
  project_id = var.atlas_project_id

  aws_kms_config {
    require_private_networking = true

    enabled                = true
    customer_master_key_id = var.aws_kms_key_id
    region                 = var.atlas_aws_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }
}

# Creates private endpoint
resource "mongodbatlas_encryption_at_rest_private_endpoint" "endpoint" {
  project_id     = mongodbatlas_encryption_at_rest.ear.project_id
  cloud_provider = "AWS"
  region_name    = var.atlas_aws_region
}
