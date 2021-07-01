# Container example provided but not always required, 
# see network_container documentation for details. 
resource "mongodbatlas_network_container" "test" {
  project_id       = var.project_id
  atlas_cidr_block = "10.8.0.0/18"
  provider_name    = "GCP"
}

# Create the peering connection request
resource "mongodbatlas_network_peering" "test" {
  project_id     = var.project_id
  container_id   = mongodbatlas_network_container.test.container_id
  provider_name  = "GCP"
  gcp_project_id = var.gcpprojectid
  network_name   = "default"
}

# the following assumes a GCP provider is configured
data "google_compute_network" "default" {
  name = "default"
}

# Create the GCP peer
resource "google_compute_network_peering" "peering" {
  name         = "peering-gcp-terraform-test"
  network      = data.google_compute_network.default.self_link
  peer_network = "https://www.googleapis.com/compute/v1/projects/${mongodbatlas_network_peering.test.atlas_gcp_project_id}/global/networks/${mongodbatlas_network_peering.test.atlas_vpc_name}"
}
resource "mongodbatlas_project_ip_access_list" "test" {
  project_id = var.project_id
  cidr_block = var.gcp_cidr
  comment    = "cidr block for GCP VPC Whitelist"
}
