terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
    }
  }
  required_version = ">= 1.0"
}
