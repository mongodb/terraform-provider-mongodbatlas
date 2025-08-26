provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

provider "aws" {
  alias  = "s3_region"
  region = var.region
}
