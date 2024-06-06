terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>3.106.1"
    }
    cloudinit = {
      source  = "hashicorp/cloudinit"
      version = "2.3.4"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = ">=1.17.0"
    }
  }
}
provider "mongodbatlas" {}

provider "azurerm" {
  features {}
}
