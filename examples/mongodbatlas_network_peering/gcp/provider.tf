provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}
provider "google" {
  # credentials = file("service-account.json")
  project = var.gcpprojectid
  region  = var.gcp_region
  # zone="us-central-1c"
}
