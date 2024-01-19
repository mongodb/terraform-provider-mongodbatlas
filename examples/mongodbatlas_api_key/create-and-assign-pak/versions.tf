terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.13.2"
    }
  }
  required_version = ">= 1.0"
}
