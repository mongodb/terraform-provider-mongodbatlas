# v2: Migration Phase - Both Legacy and Port-Based Architectures
# This configuration creates both architectures in parallel for testing

# Legacy endpoint (from v1, required for legacy architecture)
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

# New: Create endpoint with port-based architecture
resource "mongodbatlas_privatelink_endpoint" "test_new" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  port_mapping_enabled     = true
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
}

# Keep existing Google Network (from v1, used for both legacy and new architectures)
resource "google_compute_network" "default" {
  project = var.gcp_project_id
  name    = "my-network"
}

# Keep existing Google Sub Network (from v1, used for both legacy and new architectures)
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = "my-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# Legacy: Google 50 Addresses (from v1, required for legacy architecture)
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

# New: Create Google Address (1 address for new port-based architecture)
# Note: Uses existing network and subnet from v1
resource "google_compute_address" "new" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-port-based-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.100"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test_new]
}

# Legacy: 50 Forwarding rules (from v1, required for legacy architecture)
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

# New: Create Forwarding Rule (1 rule for new port-based architecture)
resource "google_compute_forwarding_rule" "new" {
  target                = mongodbatlas_privatelink_endpoint.test_new.service_attachment_names[0]
  project               = google_compute_address.new.project
  region                = google_compute_address.new.region
  name                  = google_compute_address.new.name
  ip_address            = google_compute_address.new.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# Legacy endpoint service (from v1, required for legacy architecture)
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

# New: Create Endpoint Service with port-based architecture
resource "mongodbatlas_privatelink_endpoint_service" "test_new" {
  project_id                  = mongodbatlas_privatelink_endpoint.test_new.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.test_new.private_link_id
  provider_name               = "GCP"
  endpoint_service_id         = google_compute_forwarding_rule.new.name
  private_endpoint_ip_address = google_compute_address.new.address
  gcp_project_id              = var.gcp_project_id
  delete_on_create_timeout    = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
}

data "mongodbatlas_advanced_cluster" "cluster" {
  count      = var.cluster_name == "" ? 0 : 1
  project_id = mongodbatlas_privatelink_endpoint_service.test_new.project_id
  name       = var.cluster_name
}

locals {
  endpoint_service_id = google_compute_forwarding_rule.new.name
  private_endpoints   = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings : cs.private_endpoint]), [])
  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id)
  ]
}

output "connection_string_legacy" {
  description = "Connection string for legacy endpoint"
  value       = mongodbatlas_privatelink_endpoint_service.test_legacy.id
}

output "connection_string_new" {
  description = "Connection string for new port-based endpoint"
  value       = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
