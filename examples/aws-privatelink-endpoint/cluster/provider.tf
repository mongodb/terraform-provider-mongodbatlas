provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}
provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = "us-east-1"
}
