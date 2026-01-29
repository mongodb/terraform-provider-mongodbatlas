output "endpoint_service_id" {
  description = "Endpoint service ID (returns port-mapped architecture value; legacy value available via legacy_endpoint_service_id output)"
  value       = mongodbatlas_privatelink_endpoint_service.port_mapped.endpoint_service_id
}

output "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture"
  value       = mongodbatlas_privatelink_endpoint_service.legacy.endpoint_service_id
}

output "port_mapped_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture"
  value       = mongodbatlas_privatelink_endpoint_service.port_mapped.endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID (returns port-mapped architecture value; legacy value available via mongodbatlas_privatelink_endpoint_legacy output)"
  value       = mongodbatlas_privatelink_endpoint.port_mapped.private_link_id
}

output "project_id" {
  description = "MongoDB Atlas project ID"
  value       = var.project_id
}

output "mongodbatlas_privatelink_endpoint_service_legacy" {
  description = "Full mongodbatlas_privatelink_endpoint_service resource for legacy architecture"
  value       = mongodbatlas_privatelink_endpoint_service.legacy
}

output "mongodbatlas_privatelink_endpoint_service_port_mapped" {
  description = "Full mongodbatlas_privatelink_endpoint_service resource for port-mapped architecture"
  value       = mongodbatlas_privatelink_endpoint_service.port_mapped
}

output "mongodbatlas_privatelink_endpoint_legacy" {
  description = "Full mongodbatlas_privatelink_endpoint resource for legacy architecture"
  value       = mongodbatlas_privatelink_endpoint.legacy
}

output "mongodbatlas_privatelink_endpoint_port_mapped" {
  description = "Full mongodbatlas_privatelink_endpoint resource for port-mapped architecture"
  value       = mongodbatlas_privatelink_endpoint.port_mapped
}

output "legacy_connection_string" {
  description = "Connection string for legacy endpoint"
  value       = length(local.legacy_connection_strings) > 0 ? local.legacy_connection_strings[0] : ""
}

output "port_mapped_connection_string" {
  description = "Connection string for port-mapped endpoint"
  value       = length(local.port_mapped_connection_strings) > 0 ? local.port_mapped_connection_strings[0] : ""
}
