# Authenticate the Atlas provider using the JWT stored in Secrets Manager by step 1.
provider "mongodbatlas" {
  access_token = data.aws_secretsmanager_secret_version.atlas_token.secret_string
}

provider "aws" {
  region = var.aws_region
}
