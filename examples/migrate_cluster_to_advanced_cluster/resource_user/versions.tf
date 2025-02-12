terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.26" # TODO: Update me
    }
  }
  required_version = ">= 1.5" # TODO: Update me
}
