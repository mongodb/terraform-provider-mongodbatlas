provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

provider "gcp" {
  region     = var.gcp_region
  access_key = var.access_key
  secret_key = var.secret_key
}
