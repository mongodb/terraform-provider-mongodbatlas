# Phase 1 (App Bootstrap): Create a project-scoped Service Account, store its
# credentials in AWS Secrets Manager, and trigger the app-ongoing phase.

resource "mongodbatlas_project_service_account" "this" {
  project_id                 = var.project_id
  name                       = var.sa_name
  description                = "App team Service Account for ongoing Atlas operations."
  roles                      = ["GROUP_OWNER"]
  secret_expires_after_hours = 2160
}

resource "mongodbatlas_service_account_secret" "this" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_project_service_account.this.client_id
  secret_expires_after_hours = 2160
}

resource "aws_secretsmanager_secret" "sa_creds" {
  name = var.secret_name
}

resource "aws_secretsmanager_secret_version" "sa_creds" {
  secret_id = aws_secretsmanager_secret.sa_creds.id
  secret_string_wo = jsonencode({
    client_id     = mongodbatlas_project_service_account.this.client_id
    client_secret = mongodbatlas_service_account_secret.this.secret
  })
  secret_string_wo_version = var.token_version
}

resource "terraform_data" "trigger_ongoing" {
  provisioner "local-exec" {
    working_dir = "${path.module}/../phase-2-app-ongoing"
    command     = "terraform init -input=false && terraform apply -auto-approve"
    environment = {
      TF_VAR_aws_secret_id = aws_secretsmanager_secret.sa_creds.arn
      TF_VAR_project_id    = var.project_id
      TF_VAR_aws_region    = var.aws_region
    }
  }
}
