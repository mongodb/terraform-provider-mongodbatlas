# v1: Initial State - Legacy Architecture Only
# This configuration uses the legacy GCP architecture with dedicated resources per Atlas node

resource "mongodbatlas_privatelink_endpoint" "test_legacy" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
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

# Create Google 50 Addresses (required for legacy architecture)
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

# Create 50 Forwarding rules (required for legacy architecture)
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

resource "mongodbatlas_privatelink_endpoint_service" "test_legacy" {
  project_id               = mongodbatlas_privatelink_endpoint.test_legacy.project_id
  private_link_id          = mongodbatlas_privatelink_endpoint.test_legacy.private_link_id
  provider_name            = "GCP"
  endpoint_service_id      = "legacy-endpoint-group"
  gcp_project_id           = var.gcp_project_id
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
  dynamic "endpoints" {
    for_each = google_compute_address.legacy

    content {
      ip_address    = endpoints.value["address"]
      endpoint_name = google_compute_forwarding_rule.legacy[endpoints.key].name
    }
  }

  depends_on = [google_compute_forwarding_rule.legacy]
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count      = var.cluster_name == "" ? 0 : 1
  project_id = mongodbatlas_privatelink_endpoint_service.test_legacy.project_id
  name       = var.cluster_name
}

locals {
  endpoint_service_id = "legacy-endpoint-group"
  private_endpoints   = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings : cs.private_endpoint]), [])
  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id)
  ]
}

output "connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
