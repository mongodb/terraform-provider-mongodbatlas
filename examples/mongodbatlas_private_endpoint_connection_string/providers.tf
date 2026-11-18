# Credentials for both providers come from environment variables:
# MongoDB Atlas: MONGODB_ATLAS_CLIENT_ID and MONGODB_ATLAS_CLIENT_SECRET (Service Accounts)
# AWS: standard AWS credential chain (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ...)
provider "mongodbatlas" {}

provider "aws" {
  region = var.aws_region
}
