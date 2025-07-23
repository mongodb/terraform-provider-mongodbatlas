
terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.26"
    }
  }
  required_version = ">= 1.8"
}

