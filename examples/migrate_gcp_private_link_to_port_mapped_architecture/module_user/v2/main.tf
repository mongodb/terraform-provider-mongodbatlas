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

  project_id                 = var.project_id
  gcp_project_id             = var.gcp_project_id
  gcp_region                 = var.gcp_region
  cluster_name               = var.cluster_name
  legacy_endpoint_count      = var.legacy_endpoint_count
  legacy_endpoint_service_id = var.legacy_endpoint_service_id
  new_endpoint_service_id    = var.new_endpoint_service_id
}

output "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture"
  value       = module.private_link.legacy_endpoint_service_id
}

output "new_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture"
  value       = module.private_link.new_endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID"
  value       = module.private_link.private_link_id
}

output "connection_string_legacy" {
  description = "Connection string for legacy endpoint"
  value       = module.private_link.connection_string_legacy
}

output "connection_string_new" {
  description = "Connection string for port-mapped endpoint"
  value       = module.private_link.connection_string_new
}
