terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    mongodbatlas = {
      source = "terraform-providers/mongodbatlas"
    }
  }
  required_version = ">= 0.13"
}
