# v2: Migration Phase - Both Legacy and Port-Mapped Architectures
# This configuration creates both architectures in parallel for testing

# from v1, legacy architecture
resource "mongodbatlas_privatelink_endpoint" "test_legacy" {
  project_id    = var.project_id
  provider_name = "GCP"
  region        = var.gcp_region
}

# New: Create mongodbatlas_privatelink_endpoint with port-mapped architecture
resource "mongodbatlas_privatelink_endpoint" "test_new" {
  project_id           = var.project_id
  provider_name        = "GCP"
  region               = var.gcp_region
  port_mapping_enabled = true
}

# from v1, also used for the port-mapped architecture
resource "google_compute_network" "default" {
  project = var.gcp_project_id
  name    = "my-network"
}

# from v1, also used for the port-mapped architecture
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = "my-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# from v1, legacy architecture
resource "google_compute_address" "legacy" {
  count        = 50
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-legacy${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.${count.index}"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test_legacy]
}

# New: Create Google Address (1 address for port-mapped architecture)
# Note: Uses existing network and subnet from v1
resource "google_compute_address" "new" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-port-mapped-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.100"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test_new]
}

# from v1, legacy architecture
resource "google_compute_forwarding_rule" "legacy" {
  count                 = 50
  target                = mongodbatlas_privatelink_endpoint.test_legacy.service_attachment_names[count.index]
  project               = google_compute_address.legacy[count.index].project
  region                = google_compute_address.legacy[count.index].region
  name                  = google_compute_address.legacy[count.index].name
  ip_address            = google_compute_address.legacy[count.index].id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# New: Create Forwarding Rule (1 rule for port-mapped architecture)
resource "google_compute_forwarding_rule" "new" {
  target                = mongodbatlas_privatelink_endpoint.test_new.service_attachment_names[0]
  project               = google_compute_address.new.project
  region                = google_compute_address.new.region
  name                  = google_compute_address.new.name
  ip_address            = google_compute_address.new.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# from v1, legacy architecture
resource "mongodbatlas_privatelink_endpoint_service" "test_legacy" {
  project_id          = mongodbatlas_privatelink_endpoint.test_legacy.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.test_legacy.private_link_id
  provider_name       = "GCP"
  endpoint_service_id = "legacy-endpoint-group"
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
resource "mongodbatlas_privatelink_endpoint_service" "test_new" {
  project_id                  = mongodbatlas_privatelink_endpoint.test_new.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.test_new.private_link_id
  provider_name               = "GCP"
  endpoint_service_id         = google_compute_forwarding_rule.new.name
  private_endpoint_ip_address = google_compute_address.new.address
  gcp_project_id              = var.gcp_project_id
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count      = var.cluster_name == "" ? 0 : 1
  project_id = mongodbatlas_privatelink_endpoint_service.test_new.project_id
  name       = var.cluster_name
}

locals {
  endpoint_service_id_new    = mongodbatlas_privatelink_endpoint_service.test_new.endpoint_service_id
  endpoint_service_id_legacy = mongodbatlas_privatelink_endpoint_service.test_legacy.endpoint_service_id
  private_endpoints          = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings.private_endpoint : cs]), [])

  connection_strings_new = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id_new)
  ]
  connection_strings_legacy = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id_legacy)
  ]
}

output "connection_string_legacy" {
  description = "Connection string for legacy endpoint"
  value       = length(local.connection_strings_legacy) > 0 ? local.connection_strings_legacy[0] : ""
}

output "connection_string_new" {
  description = "Connection string for port-mapped endpoint"
  value       = length(local.connection_strings_new) > 0 ? local.connection_strings_new[0] : ""
}
