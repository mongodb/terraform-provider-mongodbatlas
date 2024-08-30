# MongoDB Atlas Provider -- Encryption At Rest using Customer Key Management via Private Network Interfaces (Azure)
This example shows how to configure encryption at rest using an Azure ensuring all API calls to their KMS instances travel exclusively over their KMS cloud provider’s private network infrastructure.

## Dependencies

* Terraform MongoDB Atlas Provider v1.19.0 minimum
* A MongoDB Atlas account 
* * Terraform Azure `azapi` provider
* A Microsoft Azure account

## Usage

**1\. Provide the appropriate values for the input variables.**

- `atlas_public_key`: The public API key for MongoDB Atlas
- `atlas_private_key`: The private API key for MongoDB Atlas
- `atlas_project_id`: Atlas Project ID
- `azure_subscription_id`: Azure ID that identifies your Azure subscription
- `azure_client_id`: Azure ID identifies an Azure application associated with your Azure Active Directory tenant
- `azure_client_secret`: Secret associated to the Azure application
- `azure_tenant_id`: Azure ID  that identifies the Azure Active Directory tenant within your Azure subscription
- `azure_resource_group_name`: Name of the Azure resource group that contains your Azure Key Vault
- `azure_key_vault_name`: Unique string that identifies the Azure Key Vault that contains your key
- `azure_key_identifier`: Web address with a unique key that identifies for your Azure Key Vault
- `azure_region_name`: Region in which the Encryption At Rest private endpoint is located


NOTE: The Azure application (associated to `azure_client_id`) must have at least a Key Vault Contributor role assigned in the corresponding Key Vault.

**2\. Review the Terraform plan.**

Execute the following command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the following deployments:

TODO: adjust this
- An AWS IAM Policy
- An AWS IAM Role
- An AWS S3 bucket
- An IAM role policy for the S3 bucket
- Configure Atlas to use your AWS Role
- An Atlas project in the configured Atlas organization
- Configure push-based log export to the S3 bucket for Atlas project

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

