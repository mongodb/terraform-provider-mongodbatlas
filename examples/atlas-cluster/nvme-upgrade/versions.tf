terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.9"
    }
  }
  required_version = ">= 1.0"
}
