terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.10"
    }
  }
}

# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}
