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

  project_id              = var.project_id
  gcp_project_id          = var.gcp_project_id
  gcp_region              = var.gcp_region
  cluster_name            = var.cluster_name
  new_endpoint_service_id = var.new_endpoint_service_id
}

output "new_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (forwarding rule name)"
  value       = module.private_link.new_endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID"
  value       = module.private_link.private_link_id
}

output "connection_string_new" {
  description = "Connection string for port-mapped endpoint"
  value       = module.private_link.connection_string_new
}
