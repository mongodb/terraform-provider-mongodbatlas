terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "0.7-dev"
    }
  }
  required_version = ">= 0.13"
}
