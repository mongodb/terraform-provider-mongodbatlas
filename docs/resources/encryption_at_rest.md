---
subcategory: "Encryption at Rest using Customer Key Management"
---

# Resource: mongodbatlas_encryption_at_rest

`mongodbatlas_encryption_at_rest` allows management of Encryption at Rest for an Atlas project using Customer Key Management configuration. The following providers are supported:
- [Amazon Web Services Key Management Service](https://docs.atlas.mongodb.com/security-aws-kms/#security-aws-kms)
- [Azure Key Vault](https://docs.atlas.mongodb.com/security-azure-kms/#security-azure-kms)
- [Google Cloud KMS](https://docs.atlas.mongodb.com/security-gcp-kms/#security-gcp-kms)

The [encryption at rest Terraform module](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/encryption-at-rest/mongodbatlas/latest) makes use of this resource and simplifies its use. It is currently limited to AWS KMS.

Atlas does not automatically rotate user-managed encryption keys. Defer to your preferred Encryption at Rest provider’s documentation and guidance for best practices on key rotation. Atlas automatically creates a 90-day key rotation alert when you configure Encryption at Rest using your Key Management in an Atlas project.

See [Encryption at Rest](https://docs.atlas.mongodb.com/security-kms-encryption/index.html) for more information, including prerequisites and restrictions.

~> **IMPORTANT** By default, Atlas enables encryption at rest for all cluster storage and snapshot volumes.

~> **IMPORTANT** Atlas limits this feature to dedicated cluster tiers of M10 and greater. For more information see: https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-encryption-at-rest-using-customer-key-management

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.


-> **IMPORTANT NOTE** To disable the encryption at rest with customer key management for a project all existing clusters in the project must first either have encryption at rest for the provider set to none, e.g. `encryption_at_rest_provider = "NONE"`, or be deleted.

## Enabling Encryption at Rest for existing Atlas cluster

After configuring at least one key management provider for an Atlas project, Project Owners can enable customer key management for each Atlas cluster for which they require encryption. For clusters defined in terraform, the [`encryption_at_rest_provider` attribute](advanced_cluster#encryption_at_rest_provider) can be used in both `mongodbatlas_advanced_cluster` and `mongodbatlas_cluster` resources. The key management provider does not have to match the cluster cloud service provider.

Please reference [Enable Customer Key Management for an Atlas Cluster](https://www.mongodb.com/docs/atlas/security-kms-encryption/#enable-customer-key-management-for-an-service-cluster) documentation for additional considerations.


## Example Usages

### Configuring encryption at rest using customer key management in AWS
The configuration of encryption at rest with customer key management, `mongodbatlas_encryption_at_rest`, needs to be completed before a cluster is created in the project. Force this wait by using an implicit dependency via `project_id` as shown in the example below.

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

  enabled_for_search_nodes = true
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id                  = mongodbatlas_encryption_at_rest.test.project_id
  name                        = "MyCluster"
  cluster_type                = "REPLICASET"
  backup_enabled              = true
  encryption_at_rest_provider = "AWS"

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

data "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_encryption_at_rest.test.project_id
}

output "is_aws_kms_encryption_at_rest_valid" {
  value = data.mongodbatlas_encryption_at_rest.test.aws_kms_config.valid
}
```

**NOTE**  If using the two resources path for cloud provider access, `cloud_provider_access_setup` and `cloud_provider_access_authorization`, you may need to define a `depends_on` statement for these two resources, because terraform is not able to infer the dependency.

```terraform
resource "mongodbatlas_encryption_at_rest" "default" {
  (...)
  depends_on = [mongodbatlas_cloud_provider_access_setup.<resource_name>, mongodbatlas_cloud_provider_access_authorization.<resource_name>]
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

#### Manage Customer Keys with Azure Key Vault Over Private Endpoints
It is possible to configure Atlas Encryption at Rest to communicate with Customer Managed Keys (Azure Key Vault or AWS KMS) over private network interfaces (Azure Private Link or AWS PrivateLink). This requires enabling the `azure_key_vault_config.require_private_networking` or the `aws_kms_config.require_private_networking` attribute, together with the configuration of the `mongodbatlas_encryption_at_rest_private_endpoint` resource. 

Please review the [`mongodbatlas_encryption_at_rest_private_endpoint` resource documentation](encryption_at_rest_private_endpoint) and [complete the example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_encryption_at_rest_private_endpoint/) for details on this functionality.


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
```

### Further Examples
- [AWS KMS Encryption at Rest](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_encryption_at_rest/aws)
- [Azure Key Vault Encryption at Rest](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_encryption_at_rest/azure)


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project.

### Optional

- `aws_kms_config` (Block List) Amazon Web Services (AWS) KMS configuration details and encryption at rest configuration set for the specified project. (see [below for nested schema](#nestedblock--aws_kms_config))
- `azure_key_vault_config` (Block List) Details that define the configuration of Encryption at Rest using Azure Key Vault (AKV). (see [below for nested schema](#nestedblock--azure_key_vault_config))
- `enabled_for_search_nodes` (Boolean) Flag that indicates whether Encryption at Rest for Dedicated Search Nodes is enabled in the specified project.
- `google_cloud_kms_config` (Block List) Details that define the configuration of Encryption at Rest using Google Cloud Key Management Service (KMS). (see [below for nested schema](#nestedblock--google_cloud_kms_config))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--aws_kms_config"></a>
### Nested Schema for `aws_kms_config`

Optional:

- `access_key_id` (String, Sensitive) Unique alphanumeric string that identifies an Identity and Access Management (IAM) access key with permissions required to access your Amazon Web Services (AWS) Customer Master Key (CMK).
- `customer_master_key_id` (String, Sensitive) Unique alphanumeric string that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) you used to encrypt and decrypt the MongoDB master keys.
- `enabled` (Boolean) Flag that indicates whether someone enabled encryption at rest for the specified project through Amazon Web Services (AWS) Key Management Service (KMS). To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.
- `region` (String) Physical location where MongoDB Atlas deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Atlas creates them as part of the deployment. MongoDB Atlas assigns the VPC a CIDR block. To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.
- `require_private_networking` (Boolean) Enable connection to your Amazon Web Services (AWS) Key Management Service (KMS) over private networking.
- `role_id` (String) Unique 24-hexadecimal digit string that identifies an Amazon Web Services (AWS) Identity and Access Management (IAM) role. This IAM role has the permissions required to manage your AWS customer master key.
- `secret_access_key` (String, Sensitive) Human-readable label of the Identity and Access Management (IAM) secret access key with permissions required to access your Amazon Web Services (AWS) customer master key.

Read-Only:

- `valid` (Boolean) Flag that indicates whether the Amazon Web Services (AWS) Key Management Service (KMS) encryption key can encrypt and decrypt data.


<a id="nestedblock--azure_key_vault_config"></a>
### Nested Schema for `azure_key_vault_config`

Optional:

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

Read-Only:

- `valid` (Boolean) Flag that indicates whether the Azure encryption key can encrypt and decrypt data.


<a id="nestedblock--google_cloud_kms_config"></a>
### Nested Schema for `google_cloud_kms_config`

Optional:

- `enabled` (Boolean) Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.
- `key_version_resource_id` (String, Sensitive) Resource path that displays the key version resource ID for your Google Cloud KMS.
- `service_account_key` (String, Sensitive) JavaScript Object Notation (JSON) object that contains the Google Cloud Key Management Service (KMS). Format the JSON as a string and not as an object.

Read-Only:

- `valid` (Boolean) Flag that indicates whether the Google Cloud Key Management Service (KMS) encryption key can encrypt and decrypt data.

## Import 
Encryption at Rest Settings can be imported using project ID, in the format `project_id`, e.g.

```
$ terraform import mongodbatlas_encryption_at_rest.example 1112222b3bf99403840e8934
```

For more information see: [MongoDB Atlas API Reference for Encryption at Rest using Customer Key Management.](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-encryption-at-rest-using-customer-key-management)
