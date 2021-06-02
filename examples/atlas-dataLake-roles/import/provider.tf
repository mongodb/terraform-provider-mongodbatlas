provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
  base_url    = var.base_url
}
provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.aws_region
}
