---
page_title: "Migration Guide: Encryption at Rest (GCP) Service Account JSON to Role-based Auth"
---

# Migration Guide: Encryption at Rest (GCP) Service Account JSON to Role-based Auth

**Objective**: Migrate from using a long-lived static Service Account JSON key in `mongodbatlas_encryption_at_rest.google_cloud_kms_config.service_account_key` to role-based authentication using an Atlas-managed service account via `mongodbatlas_encryption_at_rest.google_cloud_kms_config.role_id`.

## Best Practices Before Migrating
- Back up your Terraform state file before making changes.
- Test the migration in a non-production environment if possible.

## Migration Steps

### Current (using Service Account JSON key)
```hcl
resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  google_cloud_kms_config {
    enabled                 = true
    service_account_key     = "{\"type\": \"service_account\",\"project_id\": \"my-project-common-0\", ...}"
    key_version_resource_id = "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1"
  }
}
```

### 1) Obtain the Atlas-managed GCP service account
Add the following resources to enable Atlas Cloud Provider Access for GCP and authorize it for your project:

```hcl
resource "mongodbatlas_cloud_provider_access_setup" "this" {
  project_id    = var.atlas_project_id
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "this" {
  project_id = var.atlas_project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.this.role_id
}
```

The computed attribute `mongodbatlas_cloud_provider_access_authorization.this.gcp[0].service_account_for_atlas` contains the email address of the Google Service Account managed by Atlas.

### 2) Grant KMS permissions to the Atlas service account

If your KMS key is managed with Terraform, IAM bindings can be granted with the following GCP resources:

```hcl
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
```

Alternatively, the Google Cloud CLI can be used:

```shell
gcloud kms keys add-iam-policy-binding \
  <key-name> \
    --location <location> \
    --keyring <keyring-name> \
    --member <ATLAS_OWNED_SERVICE_ACCOUNT_EMAIL> \
    --role="roles/cloudkms.cryptoKeyEncrypterDecrypter"

gcloud kms keys add-iam-policy-binding \
  <key-name> \
    --location <location> \
    --keyring <keyring-name> \
    --member <ATLAS_OWNED_SERVICE_ACCOUNT_EMAIL> \
    --role="roles/cloudkms.viewer"
```

### 3) Update the Encryption at Rest resource to use role-based auth

Replace the `service_account_key` with `role_id` using the value from the authorization resource:

```hcl
resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  google_cloud_kms_config {
    enabled                 = true
    key_version_resource_id = "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1"
    role_id                 = mongodbatlas_cloud_provider_access_authorization.this.role_id
  }
}
```

**Note:** If KMS IAM bindings are being granted within the same apply, a `depends_on` block is required in `mongodbatlas_encryption_at_rest`. This ensures bindings are correctly configured prior to the role being configured in `mongodbatlas_encryption_at_rest` resource.

```
resource "mongodbatlas_encryption_at_rest" "test" {
  ...

  depends_on = [
    google_kms_crypto_key_iam_binding.encrypter_decrypter_binding,
    google_kms_crypto_key_iam_binding.viewer_binding
  ]
}

```


Running `terraform plan` should show a change similar to:

```
# mongodbatlas_encryption_at_rest.test will be updated in-place
  ~ resource "mongodbatlas_encryption_at_rest" "test" {
        id                       = "66d6d2bdb181f8665222509b"
        # (2 unchanged attributes hidden)

      ~ google_cloud_kms_config {
          + role_id                 = "68b0448ac59ddc0496f95fa5"
          - service_account_key     = (sensitive value) -> null
          ~ valid                   = true -> (known after apply)
            # (2 unchanged attributes hidden)
        }
    }
```

### 4) Apply the changes

Run `terraform apply` to complete the migration. Once applied, the migration is complete.

## Additional Resources
- Complete role-based auth example that includes GCP KMS resource creation and IAM binding setup: [examples/mongodbatlas_encryption_at_rest/gcp](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_encryption_at_rest/gcp)

