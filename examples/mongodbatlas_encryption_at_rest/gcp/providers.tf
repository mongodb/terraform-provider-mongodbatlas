provider "google" {
  project = var.gcp_project_id
}

provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}
