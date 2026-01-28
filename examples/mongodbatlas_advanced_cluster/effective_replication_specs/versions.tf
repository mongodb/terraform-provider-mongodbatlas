terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = ">= 3.0.0"
    }
  }
  required_version = ">= 1.0"
}
