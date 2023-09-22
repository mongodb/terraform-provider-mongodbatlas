terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "1.12.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.17.0"
    }
  }
  required_version = ">= 0.13"
}