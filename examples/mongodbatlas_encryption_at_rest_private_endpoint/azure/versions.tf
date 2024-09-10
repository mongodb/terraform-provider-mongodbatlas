terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.18"
    }

    azapi = {
      source  = "Azure/azapi"
      version = "~> 1.15"
    }
  }
  required_version = ">= 1.0"
}
