# Authenticate with the short-lived JWT provided by the platform team.
provider "mongodbatlas" {
  access_token = var.access_token
}

provider "aws" {
  region = var.aws_region
}
