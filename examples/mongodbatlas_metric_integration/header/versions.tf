terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
    datadog = {
      source  = "DataDog/datadog"
      version = ">= 3.0"
    }
  }
  required_version = ">= 1.10"
}
