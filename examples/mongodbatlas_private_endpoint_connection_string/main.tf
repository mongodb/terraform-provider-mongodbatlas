# Creates an AWS private endpoint for an Atlas project, links it to Atlas, and creates a customer-managed private endpoint connection string for an existing cluster.

resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = var.atlas_region_name
}

resource "aws_vpc_endpoint" "this" {
  vpc_id             = var.vpc_id
  service_name       = mongodbatlas_privatelink_endpoint.this.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = var.subnet_ids
  security_group_ids = var.security_group_ids
}

resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.this.private_link_id
  endpoint_service_id = aws_vpc_endpoint.this.id
  provider_name       = "AWS"
}

resource "mongodbatlas_private_endpoint_connection_string" "this" {
  project_id             = var.project_id
  cluster_name           = var.cluster_name
  name                   = var.connection_string_name
  optimized_mode_enabled = var.optimized_mode_enabled
  # Reference the Atlas registration, not aws_vpc_endpoint, so the string is created after Atlas accepts the endpoint and deleted before it.
  private_endpoint_ids = [mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id]
}

data "mongodbatlas_private_endpoint_connection_string" "this" {
  project_id           = mongodbatlas_private_endpoint_connection_string.this.project_id
  cluster_name         = mongodbatlas_private_endpoint_connection_string.this.cluster_name
  connection_string_id = mongodbatlas_private_endpoint_connection_string.this.connection_string_id
}

data "mongodbatlas_private_endpoint_connection_strings" "self_managed" {
  project_id      = mongodbatlas_private_endpoint_connection_string.this.project_id
  cluster_name    = mongodbatlas_private_endpoint_connection_string.this.cluster_name
  management_type = "SELF_MANAGED"
}

output "srv_connection_string" {
  description = "Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to the cluster through the selected private endpoints."
  value       = mongodbatlas_private_endpoint_connection_string.this.srv_connection_string
}

output "connection_string" {
  description = "Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to the cluster through the selected private endpoints."
  value       = mongodbatlas_private_endpoint_connection_string.this.connection_string
}

output "connection_string_status" {
  description = "Status of the connection string, read with the singular data source."
  value       = data.mongodbatlas_private_endpoint_connection_string.this.status
}

output "self_managed_connection_string_ids" {
  description = "IDs of all self-managed private endpoint connection strings for the cluster, read with the plural data source."
  value       = [for cs in data.mongodbatlas_private_endpoint_connection_strings.self_managed.results : cs.connection_string_id]
}
