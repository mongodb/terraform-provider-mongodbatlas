provider "mongodbatlas" {
  public_key  = var.atlas_public_key
  private_key = var.atlas_private_key
  base_url    = "https://cloud-dev.mongodb.com/"
}

provider "aws" {
  region     = var.aws_region
  access_key = var.access_key
  secret_key = var.secret_key
}
