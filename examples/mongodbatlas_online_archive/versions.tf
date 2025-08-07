terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.0" // TODO: CLOUDP-335982: Update to 2.0.0
    }
  }
  required_version = ">= 1.0"
}
