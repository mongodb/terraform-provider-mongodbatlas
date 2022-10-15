terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.4.6"
    }
  }
  required_version = ">= 0.13"
}
