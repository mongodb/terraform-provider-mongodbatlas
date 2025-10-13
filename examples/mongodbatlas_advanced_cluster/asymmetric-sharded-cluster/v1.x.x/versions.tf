terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "<=1.41.0"
    }
  }
  required_version = ">= 1.0"
}
