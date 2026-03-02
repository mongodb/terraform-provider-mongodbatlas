# MongoDB Atlas Log Integration with Microsoft Azure Blob Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to a Microsoft Azure Blob storage.

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role.
- Microsoft Azure account with permissions to create Blobs.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Cloud Provider Access Setup and Authorization.
- Log Integration configuration.

### Azure
- Azure storage group.
- Azure storage account. 
- Azure storage container


## Usage

**1\. Ensure your Azure and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export AZURE_SUBSCRIPTION_ID = '<AZURE_SUBSCRIPTION_ID>'
export AZURE_CLIENT_ID='<AZURE_CLIENT_ID>'
export AZURE_CLIENT_SECRET='<AZURE_CLIENT_SECRET>'
```

... or the `~/.azure/credentials` file.

```
$ cat ~/.azure/credentials
[default]
subscription_id    = var.azure_subscription_id
client_id = var.azure_client_id
client_secret = var.azure_client_secret
```

... or follow as in the `~/.azure/variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
project_id  = var.project_id
type        = "AZURE_LOG_EXPORT"
log_types   = ["MONGOS_AUDIT"]
prefix_path            = "logs/mongodb/"
role_id   = mongodbatlas_cloud_provider_access_authorization.azure_auth.role_id
storage_account_name   = azurerm_storage_account.log_storage.name
storage_container_name = azurerm_storage_container.log_container.name
```

**2\. Review the Terraform plan.**

Execute the following command and ensure you agree with the plan.

```bash
terraform plan
```

**3\. Execute the Terraform apply.**

Now execute the plan to provision the resources.

```bash
terraform apply
```

**4\. Destroy the resources.**

When you have finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

```bash
terraform destroy
```

## Log Types

The `log_types` attribute supports the following values:
- `MONGOD` - MongoDB server logs.
- `MONGOS` - MongoDB router logs.
- `MONGOD_AUDIT` - MongoDB server audit logs.
- `MONGOS_AUDIT` - MongoDB router audit logs.

