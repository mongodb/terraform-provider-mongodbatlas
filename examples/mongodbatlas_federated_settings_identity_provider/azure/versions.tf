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
      source = "mongodb/mongodbatlas"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.2"
    }
  }

  required_version = ">= 1.3.0"
}
