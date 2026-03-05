# Example with GCP (Port-Mapped Architecture)
# This example demonstrates the port-mapped architecture.

# Create mongodbatlas_privatelink_endpoint with port-mapped architecture
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id           = var.project_id
  provider_name        = "GCP"
  region               = var.gcp_region
  port_mapping_enabled = true # Enable port-mapped architecture
}

# Create a Google Network
resource "google_compute_network" "default" {
  project                 = var.gcp_project_id
  name                    = var.network_name
  auto_create_subnetworks = false
}

# Create a Google Sub Network
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = var.subnet_name
  ip_cidr_range = var.subnet_ip_cidr_range
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# Create Google Address (1 address for port-mapped architecture)
resource "google_compute_address" "default" {
  project      = google_compute_subnetwork.default.project
  name         = var.endpoint_service_id
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = var.address_ip
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.this]
}

# Create Forwarding Rule (1 rule for port-mapped architecture)
# The service_attachment_names list will contain exactly one service attachment when using the port-mapped architecture.
resource "google_compute_forwarding_rule" "default" {
  target                = mongodbatlas_privatelink_endpoint.this.service_attachment_names[0]
  project               = google_compute_address.default.project
  region                = google_compute_address.default.region
  name                  = google_compute_address.default.name
  ip_address            = google_compute_address.default.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# Create mongodbatlas_privatelink_endpoint_service with port-mapped architecture
# For the port-mapped architecture, endpoint_service_id must match the forwarding rule name 
# and private_endpoint_ip_address the IP address. The endpoints list is no longer used for the port-mapped architecture.
resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id                  = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.this.private_link_id
  provider_name               = "GCP"
  endpoint_service_id         = google_compute_forwarding_rule.default.name
  private_endpoint_ip_address = google_compute_address.default.address
  gcp_project_id              = var.gcp_project_id
}

data "mongodbatlas_privatelink_endpoint" "this" {
  project_id      = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id = mongodbatlas_privatelink_endpoint.this.private_link_id
  provider_name   = "GCP"
  depends_on      = [mongodbatlas_privatelink_endpoint_service.this]
}

data "mongodbatlas_privatelink_endpoints" "this" {
  project_id    = mongodbatlas_privatelink_endpoint.this.project_id
  provider_name = "GCP"
  depends_on    = [mongodbatlas_privatelink_endpoint_service.this]
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count = var.cluster_name == "" ? 0 : 1
  # Use endpoint service as source of project_id to gather cluster data after endpoint changes are applied
  project_id = mongodbatlas_privatelink_endpoint_service.this.project_id
  name       = var.cluster_name

  depends_on = [mongodbatlas_privatelink_endpoint_service.this]
}

locals {
  endpoint_service_id = mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id
  private_endpoints   = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings.private_endpoint : cs]), [])

  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id)
  ]
}

output "connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}

output "privatelink_endpoint" {
  value = data.mongodbatlas_privatelink_endpoint.this
}

output "privatelink_endpoints" {
  value = data.mongodbatlas_privatelink_endpoints.this.results
}
