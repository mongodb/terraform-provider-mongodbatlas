provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

provider "aws" {
  alias  = "s3_region"
  region = var.region
}
