terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.21.3"
    }
  }
  required_version = ">= 1.0"
}
