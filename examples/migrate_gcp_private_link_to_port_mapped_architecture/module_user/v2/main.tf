provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

module "private_link" {
  source = "../../module_maintainer/v2"

  project_id                      = var.project_id
  gcp_project_id                  = var.gcp_project_id
  gcp_region                      = var.gcp_region
  cluster_name                    = var.cluster_name
  legacy_endpoint_count           = var.legacy_endpoint_count
  legacy_endpoint_service_id      = var.legacy_endpoint_service_id
  legacy_address_name_prefix      = var.legacy_address_name_prefix
  legacy_address_base_ip          = var.legacy_address_base_ip
  port_mapped_endpoint_service_id = var.port_mapped_endpoint_service_id
  port_mapped_address_ip          = var.port_mapped_address_ip
  network_name                    = var.network_name
  subnet_name                     = var.subnet_name
  subnet_ip_cidr_range            = var.subnet_ip_cidr_range
}

output "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture"
  value       = module.private_link.legacy_endpoint_service_id
}

output "port_mapped_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture"
  value       = module.private_link.port_mapped_endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID"
  value       = module.private_link.private_link_id
}

output "legacy_connection_string" {
  description = "Connection string for legacy endpoint"
  value       = module.private_link.legacy_connection_string
}

output "port_mapped_connection_string" {
  description = "Connection string for port-mapped endpoint"
  value       = module.private_link.port_mapped_connection_string
}
