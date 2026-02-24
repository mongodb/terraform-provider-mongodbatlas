terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
    google = {
      source  = "hashicorp/google"
      required_version = ">= 7.0.0"
    }
  }
}
