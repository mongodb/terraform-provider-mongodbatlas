# Step 2: Retrieve the JWT from AWS Secrets Manager and use it to create
# an Atlas project, demonstrating token-based provider authentication.

data "aws_secretsmanager_secret_version" "atlas_token" {
  secret_id = var.aws_secret_id
}

resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.org_id
}
