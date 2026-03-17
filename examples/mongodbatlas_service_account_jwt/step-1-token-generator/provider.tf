provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
  base_url = "https://cloud-dev.mongodb.com/"
}

provider "aws" {
  region = var.aws_region
}
