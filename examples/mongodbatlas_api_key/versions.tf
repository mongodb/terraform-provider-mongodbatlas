terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.37.0"
    }
  }
  required_version = ">= 1.0"
}