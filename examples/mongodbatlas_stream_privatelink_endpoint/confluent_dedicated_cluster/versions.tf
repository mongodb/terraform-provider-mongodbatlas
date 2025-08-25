terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
    confluent = {
      source  = "confluentinc/confluent"
      version = "2.12.0"
    }
  }
  required_version = ">= 1.0"
}
