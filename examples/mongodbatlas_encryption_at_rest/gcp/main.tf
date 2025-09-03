resource "mongodbatlas_cloud_provider_access_setup" "this" {
  project_id    = var.atlas_project_id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "this" {
  project_id = var.atlas_project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.this.role_id
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  google_cloud_kms_config {
    enabled                 = true
    key_version_resource_id = google_kms_crypto_key.crypto_key.primary[0].name
    role_id = mongodbatlas_cloud_provider_access_authorization.this.role_id
  }

  depends_on = [ 
	google_kms_crypto_key_iam_binding.encrypter_decrypter_binding, 
	google_kms_crypto_key_iam_binding.viewer_binding 
]
}
