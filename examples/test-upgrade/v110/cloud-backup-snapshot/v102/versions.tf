terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "1.0.2"
    }
  }
  required_version = ">= 0.13"
}
