resource "mongodbatlas_stream_privatelink_endpoint" "gcp_pubsub" {
  project_id = var.project_id

  provider_name = "GCP"
  vendor        = "PUBSUB"
  region        = var.gcp_region
}

data "mongodbatlas_stream_privatelink_endpoint" "gcp_pubsub" {
  project_id = var.project_id
  id         = mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.id
}

output "privatelink_endpoint_id" {
  description = "The ID of the MongoDB Atlas Stream Private Link Endpoint"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.id
}

output "privatelink_endpoint_state" {
  description = "The state of the MongoDB Atlas Stream Private Link Endpoint"
  value       = data.mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.state
}

output "dns_domain" {
  description = "The DNS domain computed by the API for the GCP Pub/Sub private link"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_pubsub.dns_domain
}
