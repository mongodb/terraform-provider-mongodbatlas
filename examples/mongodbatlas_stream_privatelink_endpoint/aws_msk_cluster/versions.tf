terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.81"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.30"
    }
  }
  required_version = ">= 1.0"
}
