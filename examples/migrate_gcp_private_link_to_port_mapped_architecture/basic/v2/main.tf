# v2: Migration Phase - Both Legacy and Port-Mapped Architectures
# This configuration creates both architectures in parallel for testing

# from v1, legacy architecture
resource "mongodbatlas_privatelink_endpoint" "legacy" {
  project_id    = var.project_id
  provider_name = "GCP"
  region        = var.gcp_region
}

# New: Create mongodbatlas_privatelink_endpoint with port-mapped architecture
resource "mongodbatlas_privatelink_endpoint" "port_mapped" {
  project_id           = var.project_id
  provider_name        = "GCP"
  region               = var.gcp_region
  port_mapping_enabled = true
}

# from v1, also used for the port-mapped architecture
resource "google_compute_network" "default" {
  project = var.gcp_project_id
  name    = var.network_name
}

# from v1, also used for the port-mapped architecture
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = var.subnet_name
  ip_cidr_range = var.subnet_ip_cidr_range
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# from v1, legacy architecture
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

# New: Create Google Address (1 address for port-mapped architecture)
# Note: Uses existing network and subnet from v1
resource "google_compute_address" "port_mapped" {
  project      = google_compute_subnetwork.default.project
  name         = var.port_mapped_endpoint_service_id
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = var.port_mapped_address_ip
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.port_mapped]
}

# from v1, legacy architecture
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

# New: Create Forwarding Rule (1 rule for port-mapped architecture)
resource "google_compute_forwarding_rule" "port_mapped" {
  target                = mongodbatlas_privatelink_endpoint.port_mapped.service_attachment_names[0]
  project               = google_compute_address.port_mapped.project
  region                = google_compute_address.port_mapped.region
  name                  = google_compute_address.port_mapped.name
  ip_address            = google_compute_address.port_mapped.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# from v1, legacy architecture
resource "mongodbatlas_privatelink_endpoint_service" "legacy" {
  project_id          = mongodbatlas_privatelink_endpoint.legacy.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.legacy.private_link_id
  provider_name       = "GCP"
  endpoint_service_id = var.legacy_endpoint_service_id
  gcp_project_id      = var.gcp_project_id
  dynamic "endpoints" {
    for_each = google_compute_address.legacy

    content {
      ip_address    = endpoints.value["address"]
      endpoint_name = google_compute_forwarding_rule.legacy[endpoints.key].name
    }
  }
}

# New: Create mongodbatlas_privatelink_endpoint_service with port-mapped architecture
resource "mongodbatlas_privatelink_endpoint_service" "port_mapped" {
  project_id                  = mongodbatlas_privatelink_endpoint.port_mapped.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.port_mapped.private_link_id
  provider_name               = "GCP"
  endpoint_service_id         = google_compute_forwarding_rule.port_mapped.name
  private_endpoint_ip_address = google_compute_address.port_mapped.address
  gcp_project_id              = var.gcp_project_id
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count      = var.cluster_name == "" ? 0 : 1
  project_id = mongodbatlas_privatelink_endpoint_service.port_mapped.project_id
  name       = var.cluster_name

  depends_on = [
    mongodbatlas_privatelink_endpoint_service.legacy,
    mongodbatlas_privatelink_endpoint_service.port_mapped
  ]
}

locals {
  port_mapped_endpoint_service_id = mongodbatlas_privatelink_endpoint_service.port_mapped.endpoint_service_id
  legacy_endpoint_service_id      = mongodbatlas_privatelink_endpoint_service.legacy.endpoint_service_id
  private_endpoints               = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings.private_endpoint : cs]), [])

  port_mapped_connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.port_mapped_endpoint_service_id)
  ]
  legacy_connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.legacy_endpoint_service_id)
  ]
}

output "legacy_connection_string" {
  description = "Connection string for legacy endpoint"
  value       = length(local.legacy_connection_strings) > 0 ? local.legacy_connection_strings[0] : ""
}

output "port_mapped_connection_string" {
  description = "Connection string for port-mapped endpoint"
  value       = length(local.port_mapped_connection_strings) > 0 ? local.port_mapped_connection_strings[0] : ""
}
