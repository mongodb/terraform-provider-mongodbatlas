provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}
provider "google" {
  credentials = file("service-account.json")
  project     = var.gcp_project_id
  region      = var.gcp_region # us-central1
}
