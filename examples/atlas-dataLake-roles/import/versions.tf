terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
      version = "0.1.0-dev"
    }
  }
  required_version = ">= 0.15"
}
