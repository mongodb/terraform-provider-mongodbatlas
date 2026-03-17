# Phase 0 (Platform): Create a project for the app team, generate a short-lived
# JWT, and trigger the app-bootstrap phase automatically.

resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.org_id
}

ephemeral "mongodbatlas_service_account_jwt" "token" {}

resource "terraform_data" "trigger_bootstrap" {
  provisioner "local-exec" {
    working_dir = "${path.module}/../step-2-app-bootstrap"
    command     = "terraform init -input=false && terraform apply -auto-approve"
    environment = {
      TF_VAR_access_token = ephemeral.mongodbatlas_service_account_jwt.token.access_token
      TF_VAR_project_id   = mongodbatlas_project.this.id
      TF_VAR_org_id       = var.org_id
    }
  }
}
