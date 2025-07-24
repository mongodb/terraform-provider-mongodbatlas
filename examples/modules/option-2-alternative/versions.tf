
terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.26"
    }
    external = {
      source  = "hashicorp/external"
      version = "~>2.0"
    }
  }
  required_version = ">= 1.8"
}
