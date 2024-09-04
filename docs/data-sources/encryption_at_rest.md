# Data Source: mongodbatlas_encryption_at_rest

`mongodbatlas_encryption_at_rest` describes encryption at rest configuration for an Atlas project with one of the following providers:

[Amazon Web Services Key Management Service](https://docs.atlas.mongodb.com/security-aws-kms/#security-aws-kms)
[Azure Key Vault](https://docs.atlas.mongodb.com/security-azure-kms/#security-azure-kms)
[Google Cloud KMS](https://docs.atlas.mongodb.com/security-gcp-kms/#security-gcp-kms)


~> **IMPORTANT** Atlas encrypts all cluster storage and snapshot volumes, securing all cluster data on disk: a concept known as encryption at rest, by default.

~> **IMPORTANT** Atlas limits this feature to dedicated cluster tiers of M10 and greater. For more information see: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Encryption-at-Rest-using-Customer-Key-Management

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.


## Example Usages

### Configuring encryption at rest using customer key management in AWS
```terraform
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.atlas_project_id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.atlas_project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  aws_kms_config {
    enabled                = true
    customer_master_key_id = aws_kms_key.kms_key.id
    region                 = var.atlas_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id                  = mongodbatlas_encryption_at_rest.test.project_id
  name                        = "MyCluster"
  cluster_type                = "REPLICASET"
  backup_enabled              = true
  encryption_at_rest_provider = "AWS"

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

data "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_encryption_at_rest.test.project_id
}

output "is_aws_kms_encryption_at_rest_valid" {
  value = data.mongodbatlas_encryption_at_rest.test.aws_kms_config.valid
}
```

### Configuring encryption at rest using customer key management in Azure
```terraform
resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  azure_key_vault_config {
    enabled           = true
    azure_environment = "AZURE"

    tenant_id       = var.azure_tenant_id
    subscription_id = var.azure_subscription_id
    client_id       = var.azure_client_id
    secret          = var.azure_client_secret

    resource_group_name = var.azure_resource_group_name
    key_vault_name      = var.azure_key_vault_name
    key_identifier      = var.azure_key_identifier
  }
}

data "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_encryption_at_rest.test.project_id
}

output "is_azure_encryption_at_rest_valid" {
  value = data.mongodbatlas_encryption_at_rest.test.azure_key_vault_config.valid
}
```

-> **NOTE:** It is possible to configure Atlas Encryption at Rest to communicate with Azure Key Vault using Azure Private Link, ensuring that all traffic between Atlas and Key Vault takes place over Azure’s private network interfaces. Please review `mongodbatlas_encryption_at_rest_private_endpoint` resource for details.

### Configuring encryption at rest using customer key management in GCP
```terraform
resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  google_cloud_kms_config {
    enabled                 = true
    service_account_key     = "{\"type\": \"service_account\",\"project_id\": \"my-project-common-0\",\"private_key_id\": \"e120598ea4f88249469fcdd75a9a785c1bb3\",\"private_key\": \"-----BEGIN PRIVATE KEY-----\\nMIIEuwIBA(truncated)SfecnS0mT94D9\\n-----END PRIVATE KEY-----\\n\",\"client_email\": \"my-email-kms-0@my-project-common-0.iam.gserviceaccount.com\",\"client_id\": \"10180967717292066\",\"auth_uri\": \"https://accounts.google.com/o/oauth2/auth\",\"token_uri\": \"https://accounts.google.com/o/oauth2/token\",\"auth_provider_x509_cert_url\": \"https://www.googleapis.com/oauth2/v1/certs\",\"client_x509_cert_url\": \"https://www.googleapis.com/robot/v1/metadata/x509/my-email-kms-0%40my-project-common-0.iam.gserviceaccount.com\"}"
    key_version_resource_id = "projects/my-project-common-0/locations/us-east4/keyRings/my-key-ring-0/cryptoKeys/my-key-0/cryptoKeyVersions/1"
  }
}

data "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_encryption_at_rest.test.project_id
}

output "is_gcp_encryption_at_rest_valid" {
  value = data.mongodbatlas_encryption_at_rest.test.google_cloud_kms_config.valid
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project.

### Read-Only

- `aws_kms_config` (Attributes) Amazon Web Services (AWS) KMS configuration details and encryption at rest configuration set for the specified project. (see [below for nested schema](#nestedatt--aws_kms_config))
- `azure_key_vault_config` (Attributes) Details that define the configuration of Encryption at Rest using Azure Key Vault (AKV). (see [below for nested schema](#nestedatt--azure_key_vault_config))
- `google_cloud_kms_config` (Attributes) Details that define the configuration of Encryption at Rest using Google Cloud Key Management Service (KMS). (see [below for nested schema](#nestedatt--google_cloud_kms_config))
- `id` (String) The ID of this resource.

<a id="nestedatt--aws_kms_config"></a>
### Nested Schema for `aws_kms_config`

Read-Only:

- `access_key_id` (String, Sensitive) Unique alphanumeric string that identifies an Identity and Access Management (IAM) access key with permissions required to access your Amazon Web Services (AWS) Customer Master Key (CMK).
- `customer_master_key_id` (String, Sensitive) Unique alphanumeric string that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) you used to encrypt and decrypt the MongoDB master keys.
- `enabled` (Boolean) Flag that indicates whether someone enabled encryption at rest for the specified project through Amazon Web Services (AWS) Key Management Service (KMS). To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.
- `region` (String) Physical location where MongoDB Atlas deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. When MongoDB Atlas deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Atlas creates them as part of the deployment. MongoDB Atlas assigns the VPC a CIDR block. To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.
- `role_id` (String) Unique 24-hexadecimal digit string that identifies an Amazon Web Services (AWS) Identity and Access Management (IAM) role. This IAM role has the permissions required to manage your AWS customer master key.
- `secret_access_key` (String, Sensitive) Human-readable label of the Identity and Access Management (IAM) secret access key with permissions required to access your Amazon Web Services (AWS) customer master key.
- `valid` (Boolean) Flag that indicates whether the Amazon Web Services (AWS) Key Management Service (KMS) encryption key can encrypt and decrypt data.


<a id="nestedatt--azure_key_vault_config"></a>
### Nested Schema for `azure_key_vault_config`

Read-Only:

- `azure_environment` (String) Azure environment in which your account credentials reside.
- `client_id` (String, Sensitive) Unique 36-hexadecimal character string that identifies an Azure application associated with your Azure Active Directory tenant.
- `enabled` (Boolean) Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.
- `key_identifier` (String, Sensitive) Web address with a unique key that identifies for your Azure Key Vault.
- `key_vault_name` (String) Unique string that identifies the Azure Key Vault that contains your key.
- `require_private_networking` (Boolean) Enable connection to your Azure Key Vault over private networking.
- `resource_group_name` (String) Name of the Azure resource group that contains your Azure Key Vault.
- `secret` (String, Sensitive) Private data that you need secured and that belongs to the specified Azure Key Vault (AKV) tenant (**azureKeyVault.tenantID**). This data can include any type of sensitive data such as passwords, database connection strings, API keys, and the like. AKV stores this information as encrypted binary data.
- `subscription_id` (String, Sensitive) Unique 36-hexadecimal character string that identifies your Azure subscription.
- `tenant_id` (String, Sensitive) Unique 36-hexadecimal character string that identifies the Azure Active Directory tenant within your Azure subscription.
- `valid` (Boolean) Flag that indicates whether the Azure encryption key can encrypt and decrypt data.


<a id="nestedatt--google_cloud_kms_config"></a>
### Nested Schema for `google_cloud_kms_config`

Read-Only:

- `enabled` (Boolean) Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.
- `key_version_resource_id` (String, Sensitive) Resource path that displays the key version resource ID for your Google Cloud KMS.
- `service_account_key` (String, Sensitive) JavaScript Object Notation (JSON) object that contains the Google Cloud Key Management Service (KMS). Format the JSON as a string and not as an object.
- `valid` (Boolean) Flag that indicates whether the Google Cloud Key Management Service (KMS) encryption key can encrypt and decrypt data.

# Import 
Encryption at Rest Settings can be imported using project ID, in the format `project_id`, e.g.

```
$ terraform import mongodbatlas_encryption_at_rest.example 1112222b3bf99403840e8934
```

For more information see: [MongoDB Atlas API Reference for Encryption at Rest using Customer Key Management.](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Encryption-at-Rest-using-Customer-Key-Management)