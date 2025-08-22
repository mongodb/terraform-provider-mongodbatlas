terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }

    azapi = {
      source  = "Azure/azapi"
      version = "~> 1.15"
    }
  }
  required_version = ">= 1.0"
}
