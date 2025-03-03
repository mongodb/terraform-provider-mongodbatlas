data "mongodbatlas_encryption_at_rest_private_endpoint" "single" {
  project_id     = var.atlas_project_id
  cloud_provider = "AWS"
  id             = mongodbatlas_encryption_at_rest_private_endpoint.endpoint.id
}

output "status" {
  value = data.mongodbatlas_encryption_at_rest_private_endpoint.single.status
}
