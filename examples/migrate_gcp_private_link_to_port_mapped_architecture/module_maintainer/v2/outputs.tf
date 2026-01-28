output "endpoint_service_id" {
  description = "Endpoint service ID (returns port-mapped architecture value; legacy value available via legacy_endpoint_service_id output)"
  value       = mongodbatlas_privatelink_endpoint_service.new.endpoint_service_id
}

output "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture"
  value       = mongodbatlas_privatelink_endpoint_service.legacy.endpoint_service_id
}

output "new_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture"
  value       = mongodbatlas_privatelink_endpoint_service.new.endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID (returns port-mapped architecture value; legacy value available via mongodbatlas_privatelink_endpoint_legacy output)"
  value       = mongodbatlas_privatelink_endpoint.new.private_link_id
}

output "project_id" {
  description = "MongoDB Atlas project ID"
  value       = var.project_id
}

output "mongodbatlas_privatelink_endpoint_service_legacy" {
  description = "Full mongodbatlas_privatelink_endpoint_service resource for legacy architecture"
  value       = mongodbatlas_privatelink_endpoint_service.legacy
}

output "mongodbatlas_privatelink_endpoint_service_new" {
  description = "Full mongodbatlas_privatelink_endpoint_service resource for port-mapped architecture"
  value       = mongodbatlas_privatelink_endpoint_service.new
}

output "mongodbatlas_privatelink_endpoint_legacy" {
  description = "Full mongodbatlas_privatelink_endpoint resource for legacy architecture"
  value       = mongodbatlas_privatelink_endpoint.legacy
}

output "mongodbatlas_privatelink_endpoint_new" {
  description = "Full mongodbatlas_privatelink_endpoint resource for port-mapped architecture"
  value       = mongodbatlas_privatelink_endpoint.new
}

output "connection_string_legacy" {
  description = "Connection string for legacy endpoint"
  value       = length(local.connection_strings_legacy) > 0 ? local.connection_strings_legacy[0] : ""
}

output "connection_string_new" {
  description = "Connection string for port-mapped endpoint"
  value       = length(local.connection_strings_new) > 0 ? local.connection_strings_new[0] : ""
}
