terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
  }
}
