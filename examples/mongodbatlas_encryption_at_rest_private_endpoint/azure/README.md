# MongoDB Atlas Provider - Encryption At Rest using Customer Key Management via Private Network Interfaces (Azure)
This example shows how to configure encryption at rest using Azure with customer managed keys ensuring all communication with Azure Key Vault happens exclusively over Azure Private Link.

## Dependencies

* Terraform MongoDB Atlas Provider v1.19.0 minimum
* A MongoDB Atlas account 
* Terraform Azure `azapi` provider
* A Microsoft Azure account

## Usage

**1\. Provide the appropriate values for the input variables.**

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `atlas_project_id`: Atlas Project ID
- `azure_subscription_id`: Azure ID that identifies your Azure subscription
- `azure_client_id`: Azure ID identifies an Azure application associated with your Azure Active Directory tenant
- `azure_client_secret`: Secret associated to the Azure application
- `azure_tenant_id`: Azure ID  that identifies the Azure Active Directory tenant within your Azure subscription
- `azure_resource_group_name`: Name of the Azure resource group that contains your Azure Key Vault
- `azure_key_vault_name`: Unique string that identifies the Azure Key Vault that contains your key
- `azure_key_identifier`: Web address with a unique key that identifies for your Azure Key Vault
- `azure_region_name`: Region in which the Encryption At Rest private endpoint is located


**NOTE**: The Azure application (associated to `azure_client_id`) must have the following permissions associated to the Azure Key Vault (`azure_key_vault_name`):
- GET (Key Management Operation), ENCRYPT (Cryptographic Operation) and DECRYPT (Cryptographic Operation) policy permissions.
- A `Key Vault Reader` role.

**2\. Review the Terraform plan.**

Execute the following command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project will execute the following changes to acheive a successful Azure Private Link for customer managed keys:

- Configure encryption at rest in an existing project using a custom Azure Key. For successful private networking configuration, the `requires_private_networking` attribute in `mongodbatlas_encryption_at_rest` is set to true.
- Create a private endpoint for the existing project under a certain Azure region using `mongodbatlas_encryption_at_rest_private_endpoint`. 
- Approve the connection from the Azure Key Vault. This is being done through terraform with the `azapi_update_resource` resource. Alternatively, the private connection can be approved through the Azure UI or CLI.
    - CLI example command: `az keyvault private-endpoint-connection approve --approval-description {"OPTIONAL DESCRIPTION"} --resource-group {RG} --vault-name {KEY VAULT NAME} –name {PRIVATE LINK CONNECTION NAME}`

**3\. Execute the Terraform apply.**

Now execute the plan to provision the resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

When you have finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```

