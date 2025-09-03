# Create the Key Ring
resource "google_kms_key_ring" "key_ring" {
  name     = var.key_ring_name
  location = var.location
}

# Create the Crypto Key in the Key Ring
resource "google_kms_crypto_key" "crypto_key" {  
  name     = var.crypto_key_name
  key_ring = google_kms_key_ring.key_ring.id
  purpose  = "ENCRYPT_DECRYPT"
}

# IAM Binding: Grant 'cryptoKeyEncrypterDecrypter' role
resource "google_kms_crypto_key_iam_binding" "encrypter_decrypter_binding" {  
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:${mongodbatlas_cloud_provider_access_authorization.this.gcp[0].service_account_for_atlas}"
  ]  
}
  
# IAM Binding: Grant 'viewer' role
resource "google_kms_crypto_key_iam_binding" "viewer_binding" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.viewer"

  members = [
    "serviceAccount:${mongodbatlas_cloud_provider_access_authorization.this.gcp[0].service_account_for_atlas}"
  ]
}
