# v1: Legacy Architecture Module
# This module creates GCP private link resources using the legacy architecture

# Create mongodbatlas_privatelink_endpoint with legacy architecture
resource "mongodbatlas_privatelink_endpoint" "legacy" {
  project_id    = var.project_id
  provider_name = "GCP"
  region        = var.gcp_region
  # port_mapping_enabled is not set (defaults to false for legacy architecture)
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

# Create Google Addresses (required for legacy architecture)
resource "google_compute_address" "legacy" {
  count        = var.legacy_endpoint_count
  project      = google_compute_subnetwork.default.project
  name         = "${var.legacy_address_name_prefix}${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "${var.legacy_address_base_ip}.${count.index}"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.legacy]
}

# Create Forwarding Rules (required for legacy architecture)
resource "google_compute_forwarding_rule" "legacy" {
  count                 = var.legacy_endpoint_count
  target                = mongodbatlas_privatelink_endpoint.legacy.service_attachment_names[count.index]
  project               = google_compute_address.legacy[count.index].project
  region                = google_compute_address.legacy[count.index].region
  name                  = google_compute_address.legacy[count.index].name
  ip_address            = google_compute_address.legacy[count.index].id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# Create mongodbatlas_privatelink_endpoint_service with legacy architecture
resource "mongodbatlas_privatelink_endpoint_service" "legacy" {
  project_id          = mongodbatlas_privatelink_endpoint.legacy.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.legacy.private_link_id
  provider_name       = "GCP"
  endpoint_service_id = var.legacy_endpoint_service_id
  gcp_project_id      = var.gcp_project_id

  # Legacy architecture requires the endpoints list
  dynamic "endpoints" {
    for_each = google_compute_address.legacy

    content {
      ip_address    = endpoints.value["address"]
      endpoint_name = google_compute_forwarding_rule.legacy[endpoints.key].name
    }
  }
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count      = var.cluster_name == "" ? 0 : 1
  project_id = var.project_id
  name       = var.cluster_name

  depends_on = [mongodbatlas_privatelink_endpoint_service.legacy]
}

locals {
  legacy_endpoint_service_id = mongodbatlas_privatelink_endpoint_service.legacy.endpoint_service_id
  private_endpoints          = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings.private_endpoint : cs]), [])

  legacy_connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.legacy_endpoint_service_id)
  ]
}
