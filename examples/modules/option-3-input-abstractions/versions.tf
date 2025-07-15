terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.0"
    }
  }
  required_version = ">= 1.0"
}

provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}