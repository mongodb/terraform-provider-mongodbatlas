# MongoDB Atlas Log Integration with Azure Blob Storage Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to an Azure Blob Storage container.

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- Azure account with permissions to create resource groups, storage accounts, and storage containers.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project.
- Cloud Provider Access Setup and Authorization.
- Log Integration configuration.

### Azure
- Resource group for log storage.
- Storage account.
- Storage container.

## Usage

**1\. Ensure your Azure and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export ARM_SUBSCRIPTION_ID="<AZURE_SUBSCRIPTION_ID>"
export ARM_CLIENT_ID="<AZURE_CLIENT_ID>"
export ARM_CLIENT_SECRET="<AZURE_CLIENT_SECRET>"
export ARM_TENANT_ID="<AZURE_TENANT_ID>"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
atlas_org_id               = "your-org-id"
atlas_client_id            = "your-service-account-client-id"
atlas_client_secret        = "your-service-account-client-secret"
azure_subscription_id      = "your-azure-subscription-id"
azure_client_id            = "your-azure-client-id"
azure_client_secret        = "your-azure-client-secret"
azure_tenant_id            = "your-azure-tenant-id"
atlas_azure_app_id         = "your-atlas-azure-app-id"
azure_service_principal_id = "your-azure-service-principal-id"
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

## Notes

- Atlas will add sub-directories based on the log type under the specified `prefix_path`.
