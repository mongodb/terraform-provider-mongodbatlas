terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "1.22.0"
    }
    confluent = {
      source  = "confluentinc/confluent"
      version = "2.11.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "5.0.0"
    }
  }
  required_version = ">= 1.0"
}
