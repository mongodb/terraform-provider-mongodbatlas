# MongoDB Atlas Log Integration with Datadog Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to Datadog.

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- A Datadog API key.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project.
- Log Integration configuration.

## Usage

**1\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
atlas_org_id        = "your-org-id"
atlas_client_id     = "your-service-account-client-id"
atlas_client_secret = "your-service-account-client-secret"
datadog_api_key     = "your-datadog-api-key"
datadog_region      = "your-datadog-region"
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
