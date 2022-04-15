terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
       version = "1.3.2"
    }
  }
  required_version = ">= 0.13"
}
