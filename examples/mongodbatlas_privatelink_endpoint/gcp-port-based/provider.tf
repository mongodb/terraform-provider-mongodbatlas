provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}
provider "google" {
  # Credentials commented out to pass lint
  # Add your own service-account & path to run example
  # credentials = file("service-account.json")
  project = var.gcp_project_id
  region  = var.gcp_region # us-central1
}
