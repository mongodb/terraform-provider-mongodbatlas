# MongoDB Atlas Provider -- Encryption At Rest using Customer Key Management with Azure
This example shows how to configure encryption at rest with customer managed keys with Azure Key Vault. 

Note: It is possible to configure Atlas Encryption at Rest to communicate with Azure Key Vault using Azure Private Link, ensuring that all traffic between Atlas and Key Vault takes place over Azure’s private network interfaces. Please review `mongodbatlas_encryption_at_rest_private_endpoint` resource for details.

## Dependencies

* Terraform MongoDB Atlas Provider
* A MongoDB Atlas account 
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

**NOTE**: The Azure application (associated to `azure_client_id`) must have the following permissions associated to the Azure Key Vault (`azure_key_vault_name`):
- GET (Key Management Operation), ENCRYPT (Cryptographic Operation) and DECRYPT (Cryptographic Operation) policy permissions.
- A `Key Vault Reader` role.

**2\. Review the Terraform plan.**

Execute the following command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the following deployments:

- Configure encryption at rest in an existing project using a custom Azure Key.

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

