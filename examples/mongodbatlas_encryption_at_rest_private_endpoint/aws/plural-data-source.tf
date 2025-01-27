data "mongodbatlas_encryption_at_rest_private_endpoints" "plural" {
  project_id     = var.atlas_project_id
  cloud_provider = "AWS"
}

output "number_of_endpoints" {
  value = length(data.mongodbatlas_encryption_at_rest_private_endpoints.plural.results)
}
