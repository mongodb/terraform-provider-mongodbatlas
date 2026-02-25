# MongoDB Atlas Log Integration with Microsoft Azure Blob Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to a Microsoft Azure Blob storage.

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role.
- Microsoft Azure account with permissions to create Blobs.
- Terraform >= `1.0`.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project.
- Cloud Provider Access Setup and Authorization.
- Log Integration configuration.

### Azure
- Azure Blob for storing logs.
- IAM role for Atlas to assume.
- IAM policy for Blob access.


## Usage

**1\. Ensure your Azure and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export AZURE_REGION = '<AZURE_REGION>'
export AZURE_ACCESS_KEY_ID='<AZURE_ACCESS_KEY_ID>'
export AZURE_SECRET_ACCESS_KEY='<AZURE_SECRET_ACCESS_KEY>'
```

... or the `~/.azure/credentials` file.

```
$ cat ~/.azure/credentials
[default]
region     = var.azure_region
access_key = var.access_key
secret_key = var.secret_key
```

... or follow as in the `~/.azure/variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
project_id  = var.project_id
type        = "AZURE_LOG_EXPORT"
log_types   = ["MONGOS_AUDIT"]
prefix_path            = "logs/mongodb/"
service_principal_id   = mongodbatlas_cloud_provider_access_authorization.azure_auth.role_id
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

