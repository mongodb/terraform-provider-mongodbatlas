# versions.tf
terraform {
  required_version = ">= 1.0.0"
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = ">= 1.10.0"
    }
  }
}
