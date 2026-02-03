provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

module "private_link" {
  source = "../../module_maintainer/v3"

  project_id                      = var.project_id
  gcp_project_id                  = var.gcp_project_id
  gcp_region                      = var.gcp_region
  cluster_name                    = var.cluster_name
  port_mapped_endpoint_service_id = var.port_mapped_endpoint_service_id
  port_mapped_address_ip          = var.port_mapped_address_ip
  network_name                    = var.network_name
  subnet_name                     = var.subnet_name
  subnet_ip_cidr_range            = var.subnet_ip_cidr_range
}

output "port_mapped_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (forwarding rule name)"
  value       = module.private_link.port_mapped_endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID"
  value       = module.private_link.private_link_id
}

output "port_mapped_connection_string" {
  description = "Connection string for port-mapped endpoint"
  value       = module.private_link.port_mapped_connection_string
}
