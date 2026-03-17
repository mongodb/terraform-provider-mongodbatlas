# Phase 2 (App Ongoing): Read SA credentials from Secrets Manager and provision
# infrastructure. This phase runs independently of the bootstrap JWT.

data "aws_secretsmanager_secret_version" "sa_creds" {
  secret_id = var.aws_secret_id
}

resource "mongodbatlas_flex_cluster" "this" {
  project_id = var.project_id
  name       = var.cluster_name
  provider_settings = {
    backing_provider_name = "AWS"
    region_name           = "US_EAST_1"
  }
}
