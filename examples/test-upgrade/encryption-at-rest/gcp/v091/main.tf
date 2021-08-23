resource "mongodbatlas_project" "test" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_project.test.id

  google_cloud_kms = {
    enabled                 = true
    service_account_key     = var.service_account_key
    key_version_resource_id = var.gcp_key_version_resource_id
  }
}

output "project_id" {
  value = mongodbatlas_project.test.id
}
