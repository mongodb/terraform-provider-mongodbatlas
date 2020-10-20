terraform {
  required_providers {
    mongodbatlas = {
      source = "terraform-providers/mongodbatlas"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
  required_version = ">= 0.13"
}
