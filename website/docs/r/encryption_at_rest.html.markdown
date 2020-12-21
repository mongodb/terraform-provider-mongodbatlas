---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: encryption_at_rest"
sidebar_current: "docs-mongodbatlas-resource-encryption_at_rest"
description: |-
    Provides an Encryption At Rest resource.
---

# mongodbatlas_encryption_at_rest

`mongodbatlas_encryption_at_rest` Allows management of encryption at rest for an Atlas project with one of the following providers:

[Amazon Web Services Key Management Service](https://docs.atlas.mongodb.com/security-aws-kms/#security-aws-kms)
[Azure Key Vault](https://docs.atlas.mongodb.com/security-azure-kms/#security-azure-kms)
[Google Cloud KMS](https://docs.atlas.mongodb.com/security-gcp-kms/#security-gcp-kms)

After configuring at least one Encryption at Rest provider for the Atlas project, Project Owners can enable Encryption at Rest for each Atlas cluster for which they require encryption. The Encryption at Rest provider does not have to match the cluster cloud service provider.

Atlas does not automatically rotate user-managed encryption keys. Defer to your preferred Encryption at Rest provider’s documentation and guidance for best practices on key rotation. Atlas automatically creates a 365-day key rotation alert when you configure Encryption at Rest using your Key Management in an Atlas project.

See [Encryption at Rest](https://docs.atlas.mongodb.com/security-kms-encryption/index.html) for more information, including prerequisites and restrictions.

~> **IMPORTANT** Atlas encrypts all cluster storage and snapshot volumes, securing all cluster data on disk: a concept known as encryption at rest, by default.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = "<PROJECT-ID>"

  aws_kms = {
    enabled                = true
    access_key_id          = "AKIAIOSFODNN7EXAMPLE"
    secret_access_key      = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    customer_master_key_id = "030gce02-586d-48d2-a966-05ea954fde0g"
    region                 = "US_EAST_1"
  }

  azure_key_vault = {
    enabled             = true
    client_id           = "g54f9e2-89e3-40fd-8188-EXAMPLEID"
    azure_environment   = "AZURE"
    subscription_id     = "0ec944e3-g725-44f9-a147-EXAMPLEID"
    resource_group_name = "ExampleRGName"
    key_vault_name      = "EXAMPLEKeyVault"
    key_identifier      = "https://EXAMPLEKeyVault.vault.azure.net/keys/EXAMPLEKey/d891821e3d364e9eb88fbd3d11807b86"
    secret              = "EXAMPLESECRET"
    tenant_id           = "e8e4b6ba-ff32-4c88-a9af-EXAMPLEID"
  }

  google_cloud_kms = {
    enabled                 = true
    service_account_key     = "{\"type\": \"service_account\",\"project_id\": \"my-project-common-0\",\"private_key_id\": \"e120598ea4f88249469fcdd75a9a785c1bb3\",\"private_key\": \"-----BEGIN PRIVATE KEY-----\\nMIIEuwIBA(truncated)SfecnS0mT94D9\\n-----END PRIVATE KEY-----\\n\",\"client_email\": \"my-email-kms-0@my-project-common-0.iam.gserviceaccount.com\",\"client_id\": \"10180967717292066\",\"auth_uri\": \"https://accounts.google.com/o/oauth2/auth\",\"token_uri\": \"https://accounts.google.com/o/oauth2/token\",\"auth_provider_x509_cert_url\": \"https://www.googleapis.com/oauth2/v1/certs\",\"client_x509_cert_url\": \"https://www.googleapis.com/robot/v1/metadata/x509/my-email-kms-0%40my-project-common-0.iam.gserviceaccount.com\"}"
    key_version_resource_id = "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1"
  }
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier for the project.
* `aws_kms` - (Required) Specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
* `azure_key_vault` - (Required) Specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
* `google_cloud_kms` - (Required) Specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.

### aws_kms
* `enabled` - Specifies whether Encryption at Rest is enabled for an Atlas project, To disable Encryption at Rest, pass only this parameter with a value of false, When you disable Encryption at Rest, Atlas also removes the configuration details.
* `access_key_id` - The IAM access key ID with permissions to access the customer master key specified by customerMasterKeyID.
* `secret_access_key` - The IAM secret access key with permissions to access the customer master key specified by customerMasterKeyID.
* `customer_master_key_id` - The AWS customer master key used to encrypt and decrypt the MongoDB master keys.
* `region` - The AWS region in which the AWS customer master key exists: CA_CENTRAL_1, US_EAST_1, US_EAST_2, US_WEST_1, US_WEST_2, SA_EAST_1
* `role_id` - ID of an AWS IAM role authorized to manage an AWS customer master key. To find the ID for an existing IAM role check the `role_id` attribute of the `mongodbatlas_cloud_provider_access` resource.


### azure_key_vault
* `enabled` - Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
* `client_id` - The client ID, also known as the application ID, for an Azure application associated with the Azure AD tenant.
* `azure_environment` - The Azure environment where the Azure account credentials reside. Valid values are the following: AZURE, AZURE_CHINA, AZURE_GERMANY
* `subscription_id` - The unique identifier associated with an Azure subscription.
* `resource_group_name` - The name of the Azure Resource group that contains an Azure Key Vault.
* `key_vault_name` - The name of an Azure Key Vault containing your key.
* `key_identifier` - The unique identifier of a key in an Azure Key Vault.
* `secret` - The secret associated with the Azure Key Vault specified by azureKeyVault.tenantID.
* `tenant_id` - The unique identifier for an Azure AD tenant within an Azure subscription.

### google_cloud_kms
* `enabled` - Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
* `service_account_key` - String-formatted JSON object containing GCP KMS credentials from your GCP account.
* `key_version_resource_id` - The Key Version Resource ID from your GCP account.


For more information see: [MongoDB Atlas API Reference for Encryption at Rest using Customer Key Management.](https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/)