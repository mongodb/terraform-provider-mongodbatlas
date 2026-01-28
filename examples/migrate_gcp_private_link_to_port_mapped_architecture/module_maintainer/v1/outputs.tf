output "endpoint_service_id" {
  description = "Endpoint service ID"
  value       = mongodbatlas_privatelink_endpoint_service.legacy.endpoint_service_id
}

output "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture"
  value       = mongodbatlas_privatelink_endpoint_service.legacy.endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID"
  value       = mongodbatlas_privatelink_endpoint.legacy.private_link_id
}

output "project_id" {
  description = "MongoDB Atlas project ID"
  value       = var.project_id
}

output "mongodbatlas_privatelink_endpoint_service" {
  description = "Full mongodbatlas_privatelink_endpoint_service resource"
  value       = mongodbatlas_privatelink_endpoint_service.legacy
}

output "mongodbatlas_privatelink_endpoint" {
  description = "Full mongodbatlas_privatelink_endpoint resource"
  value       = mongodbatlas_privatelink_endpoint.legacy
}

output "connection_string_legacy" {
  description = "Connection string for legacy endpoint"
  value       = length(local.connection_strings_legacy) > 0 ? local.connection_strings_legacy[0] : ""
}
