# Example with GCP with Port-Based (1 endpoint)
# This example demonstrates the new PSC port-based architecture which requires only 1 endpoint.
# The new architecture is enabled by setting port_mapping_enabled = true on the endpoint resource.
# This simplifies setup and management compared to the legacy architecture which requires multiple endpoints.
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  port_mapping_enabled     = true # Enable new PSC port-based architecture
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
}

# Create a Google Network
resource "google_compute_network" "default" {
  project = var.gcp_project_id
  name    = "my-network"
}

# Create a Google Sub Network
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = "my-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# Create Google Address (1 address for new PSC port-based architecture)
resource "google_compute_address" "default" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-psc-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.1"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test]
}

# Create Forwarding Rule (1 rule for new PSC port-based architecture)
# The service_attachment_names list will contain exactly one service attachment when using the new architecture.
resource "google_compute_forwarding_rule" "default" {
  target                = mongodbatlas_privatelink_endpoint.test.service_attachment_names[0]
  project               = google_compute_address.default.project
  region                = google_compute_address.default.region
  name                  = google_compute_address.default.name
  ip_address            = google_compute_address.default.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# Create MongoDB Atlas Private Endpoint Service
# For the new port-based architecture, endpoint_service_id must match the forwarding rule name 
# and private_endpoint_ip_address the IP address. The endpoints list is no longer used for the new architecture.
resource "mongodbatlas_privatelink_endpoint_service" "test" {
  project_id                = mongodbatlas_privatelink_endpoint.test.project_id
  private_link_id           = mongodbatlas_privatelink_endpoint.test.private_link_id
  provider_name             = "GCP"
  endpoint_service_id       = google_compute_forwarding_rule.default.name
  private_endpoint_ip_address = google_compute_address.default.address
  gcp_project_id           = var.gcp_project_id
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }

  depends_on = [google_compute_forwarding_rule.default]
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count = var.cluster_name == "" ? 0 : 1
  # Use endpoint service as source of project_id to gather cluster data after endpoint changes are applied
  project_id = mongodbatlas_privatelink_endpoint_service.test.project_id
  name       = var.cluster_name
}

locals {
  endpoint_service_id = google_compute_forwarding_rule.default.name
  private_endpoints   = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings : cs.private_endpoint]), [])
  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id)
  ]
}
output "connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
