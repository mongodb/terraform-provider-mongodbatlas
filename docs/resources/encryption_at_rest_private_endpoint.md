# Resource: mongodbatlas_encryption_at_rest_private_endpoint

`mongodbatlas_encryption_at_rest_private_endpoint` provides a resource for managing a private endpoint used for encryption at rest with customer-managed keys. This ensures all traffic between Atlas and customer key management systems take place over private network interfaces.

-> **NOTE:** As a prerequisite to configuring a private endpoint for Azure Key Vault or AWS KMS, the corresponding [`mongodbatlas_encryption_at_rest`](encryption_at_rest) resource has to be adjusted by configuring to true [`azure_key_vault_config.require_private_networking`](encryption_at_rest#require_private_networking) or [`aws_kms_config.require_private_networking`](encryption_at_rest#require_private_networking), respectively. This attribute should be updated in place, ensuring the customer-managed keys encryption is never disabled.

-> **NOTE:** This resource does not support update operations. To modify values of a private endpoint the existing resource must be deleted and a new one can be created with the modified values.

## Example Usages

-> **NOTE:** Only Azure Key Vault with Azure Private Link and AWS KMS over AWS PrivateLink is supported at this time.

### Configuring Atlas Encryption at Rest using Azure Key Vault with Azure Private Link
To learn more about existing limitations, see [Manage Customer Keys with Azure Key Vault Over Private Endpoints](https://www.mongodb.com/docs/atlas/security/azure-kms-over-private-endpoint/#manage-customer-keys-with-azure-key-vault-over-private-endpoints).

Make sure to reference the [complete example section](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_encryption_at_rest_private_endpoint/azure) for detailed steps and considerations.

```terraform
resource "mongodbatlas_encryption_at_rest" "ear" {
  project_id = var.atlas_project_id

  azure_key_vault_config {
    require_private_networking = true

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

# Creates private endpoint
resource "mongodbatlas_encryption_at_rest_private_endpoint" "endpoint" {
  project_id     = mongodbatlas_encryption_at_rest.ear.project_id
  cloud_provider = "AZURE"
  region_name    = var.azure_region_name
}

locals {
  key_vault_resource_id = "/subscriptions/${var.azure_subscription_id}/resourceGroups/${var.azure_resource_group_name}/providers/Microsoft.KeyVault/vaults/${var.azure_key_vault_name}"
}

# Approves private endpoint connection from Azure Key Vault
resource "azapi_update_resource" "approval" {
  type      = "Microsoft.KeyVault/Vaults/PrivateEndpointConnections@2023-07-01"
  name      = mongodbatlas_encryption_at_rest_private_endpoint.endpoint.private_endpoint_connection_name
  parent_id = local.key_vault_resource_id

  body = jsonencode({
    properties = {
      privateLinkServiceConnectionState = {
        description = "Approved via Terraform"
        status      = "Approved"
      }
    }
  })
}
```

### Configuring Atlas Encryption at Rest using AWS KMS with AWS PrivateLink

Make sure to reference the [complete example section](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_encryption_at_rest_private_endpoint/aws) for detailed steps and considerations.

```terraform
resource "mongodbatlas_encryption_at_rest" "ear" {
  project_id = var.atlas_project_id

  aws_kms_config {
    require_private_networking = true

    enabled                = true
    customer_master_key_id = var.aws_kms_key_id
    region                 = var.atlas_aws_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }
}

# Creates private endpoint
resource "mongodbatlas_encryption_at_rest_private_endpoint" "endpoint" {
  project_id     = mongodbatlas_encryption_at_rest.ear.project_id
  cloud_provider = "AWS"
  region_name    = var.atlas_aws_region
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_provider` (String) Label that identifies the cloud provider for the Encryption At Rest private endpoint.
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project.
- `region_name` (String) Cloud provider region in which the Encryption At Rest private endpoint is located.

### Read-Only

- `error_message` (String) Error message for failures associated with the Encryption At Rest private endpoint.
- `id` (String) Unique 24-hexadecimal digit string that identifies the Private Endpoint Service.
- `private_endpoint_connection_name` (String) Connection name of the Azure Private Endpoint.
- `status` (String) State of the Encryption At Rest private endpoint.

## Import 
Encryption At Rest Private Endpoint resource can be imported using the project ID, cloud provider, and private endpoint ID. The format must be `{project_id}-{cloud_provider}-{private_endpoint_id}` e.g.

```
$ terraform import mongodbatlas_encryption_at_rest_private_endpoint.test 650972848269185c55f40ca1-AZURE-650972848269185c55f40ca2
$ terraform import mongodbatlas_encryption_at_rest_private_endpoint.test 650972848269185c55f40ca2-AWS-650972848269185c55f40ca3
```

For more information see: 
- [MongoDB Atlas API - Private Endpoint for Encryption at Rest Using Customer Key Management](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getencryptionatrestprivateendpoint) Documentation.
- [Manage Customer Keys with Azure Key Vault Over Private Endpoints](https://www.mongodb.com/docs/atlas/security/azure-kms-over-private-endpoint/).
