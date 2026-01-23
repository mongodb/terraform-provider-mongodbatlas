# v3: Final State - Port-Based Architecture Only
# This configuration uses only the new port-based architecture

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

# Create Google Address (1 address for new GCP port-based architecture)
resource "google_compute_address" "new" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-port-based-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.100"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test_new]
}

# Create Forwarding Rule (1 rule for new GCP port-based architecture)
resource "google_compute_forwarding_rule" "new" {
  target                = mongodbatlas_privatelink_endpoint.test_new.service_attachment_names[0]
  project               = google_compute_address.new.project
  region                = google_compute_address.new.region
  name                  = google_compute_address.new.name
  ip_address            = google_compute_address.new.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

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

  depends_on = [google_compute_forwarding_rule.new]
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

output "connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
