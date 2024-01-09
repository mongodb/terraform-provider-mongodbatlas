terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.13"
    }
  }
  required_version = ">= 1.0"
}