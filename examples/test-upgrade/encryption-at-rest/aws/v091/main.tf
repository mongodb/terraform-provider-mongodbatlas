resource "mongodbatlas_project" "test" {
  name   = var.project_name
  org_id = var.org_id
}
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.test.id
  provider_name = var.cloud_provider_access_name
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.test.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws = {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_project.test.id

  aws_kms = {
    enabled                = true
    customer_master_key_id = var.customer_master_key
    region                 = var.atlas_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }
}

output "project_id" {
  value = mongodbatlas_project.test.id
}

output "role_name" {
  value = aws_iam_role.test_role.name
}
output "role_policy_name" {
  value = aws_iam_role_policy.test_policy.name
}
