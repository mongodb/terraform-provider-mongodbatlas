output "endpoint_service_id" {
  description = "Endpoint service ID (forwarding rule name for port-mapped architecture)"
  value       = mongodbatlas_privatelink_endpoint_service.port_mapped.endpoint_service_id
}

output "port_mapped_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (forwarding rule name)"
  value       = mongodbatlas_privatelink_endpoint_service.port_mapped.endpoint_service_id
}

output "private_link_id" {
  description = "Private link ID"
  value       = mongodbatlas_privatelink_endpoint.port_mapped.private_link_id
}

output "project_id" {
  description = "MongoDB Atlas project ID"
  value       = var.project_id
}

output "mongodbatlas_privatelink_endpoint_service" {
  description = "Full mongodbatlas_privatelink_endpoint_service resource"
  value       = mongodbatlas_privatelink_endpoint_service.port_mapped
}

output "mongodbatlas_privatelink_endpoint" {
  description = "Full mongodbatlas_privatelink_endpoint resource"
  value       = mongodbatlas_privatelink_endpoint.port_mapped
}

output "port_mapped_connection_string" {
  description = "Connection string for port-mapped endpoint"
  value       = length(local.port_mapped_connection_strings) > 0 ? local.port_mapped_connection_strings[0] : ""
}
