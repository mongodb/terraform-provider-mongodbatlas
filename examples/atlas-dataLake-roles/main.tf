resource "mongodbatlas_project" "test" {
  name   = var.project_name
  org_id = var.org_id
}


resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.test.id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.test.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

resource "mongodbatlas_data_lake" "test" {
  project_id         = mongodbatlas_project.test.id
  name               = var.data_lake_name
  aws_role_id        = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  aws_test_s3_bucket = var.test_s3_bucket
  data_process_region = {
    cloud_provider = "AWS"
    region         = var.data_lake_region
  }
}


output "project_id" {
  value = mongodbatlas_project.test.id
}
output "role_id" {
  value = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
}
output "role_name" {
  value = aws_iam_role.test_role.name
}
output "policy_name" {
  value = aws_iam_role_policy.test_policy.name
}
output "data_lake_name" {
  value = mongodbatlas_data_lake.test.name
}
output "s3_bucket" {
  value = mongodbatlas_data_lake.test.aws_test_s3_bucket
}
