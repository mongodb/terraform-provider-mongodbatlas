# Authenticate the Atlas provider using the JWT stored in Secrets Manager by step1.
provider "mongodbatlas" {
  access_token = data.aws_secretsmanager_secret_version.atlas_token.secret_string
  base_url = "https://cloud-dev.mongodb.com/"
}

provider "aws" {
  region = var.aws_region
}
