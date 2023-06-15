terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "0.8.0"
    }
  }
  required_version = ">= 0.13"
}
