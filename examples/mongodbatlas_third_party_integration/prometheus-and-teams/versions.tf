terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }

    template = {
      source  = "hashicorp/template"
      version = "2.2.0"
    }
  }
  required_version = ">= 1.0"
}
