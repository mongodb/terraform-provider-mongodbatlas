provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}
provider "google" {
  # Credentials commented out to pass lint
  # Add your own service-account & path to run example
  # credentials = file("service-account.json")
  project = var.gcp_project_id
  region  = var.gcp_region # us-central1
}
