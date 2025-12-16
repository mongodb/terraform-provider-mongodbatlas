terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
    }
    aws = {
      source  = "hashicorp/aws"
    }
  }
  required_version = ">= 1.0"
}

