data "mongodbatlas_encryption_at_rest_private_endpoint" "single" {
  project_id     = var.atlas_project_id
  cloud_provider = "AZURE"
  id             = mongodbatlas_encryption_at_rest_private_endpoint.endpoint.id
}

output "endpoint_connection_name" {
  value = data.mongodbatlas_encryption_at_rest_private_endpoint.single.private_endpoint_connection_name
}
