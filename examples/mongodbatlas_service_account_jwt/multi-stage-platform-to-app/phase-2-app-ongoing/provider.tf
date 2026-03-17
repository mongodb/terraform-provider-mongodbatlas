# Authenticate with the app team's own SA credentials from Secrets Manager.
locals {
  sa_creds = jsondecode(data.aws_secretsmanager_secret_version.sa_creds.secret_string)
}

provider "mongodbatlas" {
  client_id     = local.sa_creds.client_id
  client_secret = local.sa_creds.client_secret
}

provider "aws" {
  region = var.aws_region
}
