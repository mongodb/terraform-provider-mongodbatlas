resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = var.project_id
  provider_name = "GCP"
  region        = var.gcp_region
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

# Create Google 50 Addresses
resource "google_compute_address" "default" {
  count        = 50
  project      = google_compute_subnetwork.default.project
  name         = "tf-test${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.${count.index}"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test]
}

# Create 50 Forwarding rules
resource "google_compute_forwarding_rule" "default" {
  count                 = 50
  target                = mongodbatlas_privatelink_endpoint.test.service_attachment_names[count.index]
  project               = google_compute_address.default[count.index].project
  region                = google_compute_address.default[count.index].region
  name                  = google_compute_address.default[count.index].name
  ip_address            = google_compute_address.default[count.index].id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

resource "mongodbatlas_privatelink_endpoint_service" "test" {
  project_id          = mongodbatlas_privatelink_endpoint.test.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.test.private_link_id
  provider_name       = "GCP"
  endpoint_service_id = google_compute_network.default.name
  gcp_project_id      = var.gcp_project_id

  dynamic "endpoints" {
    for_each = mongodbatlas_privatelink_endpoint.test.service_attachment_names

    content {
      ip_address    = google_compute_address.default[endpoints.key].address
      endpoint_name = google_compute_forwarding_rule.default[endpoints.key].name
    }
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
  endpoint_service_id = google_compute_network.default.name
  private_endpoints   = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings : cs.private_endpoint]), [])
  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id)
  ]
}
output "connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
