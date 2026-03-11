terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
    null = {
      source = "hashicorp/null"
    }
  }
  required_version = ">= 1.10"
}
