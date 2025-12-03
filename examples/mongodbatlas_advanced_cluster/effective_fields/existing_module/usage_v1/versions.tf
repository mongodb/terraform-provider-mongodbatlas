terraform {
  required_version = ">= 1.0"
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = ">= 2.0.0"
    }
  }
}

provider "mongodbatlas" {
  # Credentials should be provided via environment variables:
  # - MONGODB_ATLAS_CLIENT_ID
  # - MONGODB_ATLAS_CLIENT_SECRET
}
