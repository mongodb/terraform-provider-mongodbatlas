terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
    }
  }
  required_version = ">= 4.0.0"
}
