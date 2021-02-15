terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
      //version = "0.7.0"
    }
  }
  required_version = ">= 0.13"
}
