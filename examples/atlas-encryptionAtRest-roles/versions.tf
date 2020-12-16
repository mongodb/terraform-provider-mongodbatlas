terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"

  provider_installation {
    filesystem_mirror {
      path    = "../"
      include = ["terraform-providers/mongodbatlas"]
    }
  }
}
