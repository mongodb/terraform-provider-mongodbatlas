terraform {
  required_providers {
    cedar = {
      source  = "common-fate/cedar"
      version = "0.2.0"
    }
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
  }
  required_version = ">= 1.0"
}
