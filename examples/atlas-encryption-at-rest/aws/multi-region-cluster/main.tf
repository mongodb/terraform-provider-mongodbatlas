module "multi-region-cluster" {
  atlas_project_id   = mongodbatlas_encryption_at_rest.test.project_id
  provider_name      = "AWS"
  aws_region_shard_1 = var.aws_region_shard_1
  aws_region_shard_2 = var.aws_region_shard_2
  cluster_name       = var.cluster_name
  instance_size      = "M10"
  source             = "./modules/multi-region-cluster"
}

resource "mongodbatlas_project" "project" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.project.id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_project.project.id
  aws_kms_config {
    enabled                = true
    customer_master_key_id = aws_kms_key.kms_key.id
    region                 = var.atlas_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }
}


