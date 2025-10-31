resource "mongodbatlas_stream_privatelink_endpoint" "gcp_confluent" {
  project_id = var.project_id

  provider_name = "GCP"
  vendor        = "CONFLUENT"
  region        = var.gcp_region

  dns_domain     = var.confluent_dns_domain
  dns_sub_domain = var.confluent_dns_subdomains

  service_attachment_uris = [
    "projects/my-project/regions/us-west1/serviceAttachments/confluent-attachment-1",
    "projects/my-project/regions/us-west1/serviceAttachments/confluent-attachment-2"
  ]
}

data "mongodbatlas_stream_privatelink_endpoint" "gcp_confluent" {
  project_id = var.project_id
  id         = mongodbatlas_stream_privatelink_endpoint.gcp_confluent.id
}

output "privatelink_endpoint_id" {
  description = "The ID of the MongoDB Atlas Stream Private Link Endpoint"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_confluent.id
}

output "privatelink_endpoint_state" {
  description = "The state of the MongoDB Atlas Stream Private Link Endpoint"
  value       = data.mongodbatlas_stream_privatelink_endpoint.gcp_confluent.state
}

output "service_attachment_uris" {
  description = "The GCP service attachment URIs used for the private link"
  value       = mongodbatlas_stream_privatelink_endpoint.gcp_confluent.service_attachment_uris
}
