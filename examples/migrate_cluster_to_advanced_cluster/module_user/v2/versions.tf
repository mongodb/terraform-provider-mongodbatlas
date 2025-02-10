terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.26" # todo: PLACEHOLDER_TPF_RELEASE
    }
  }
  required_version = ">= 1.0" # todo: minimum moved block supported version
}
