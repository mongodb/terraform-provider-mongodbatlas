terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.11"
    }
  }
  required_version = ">= 0.13"
}
