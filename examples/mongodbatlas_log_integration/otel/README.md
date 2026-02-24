# MongoDB Atlas Log Integration with Datadog

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to a Open Telemetry endpoint.

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role.
- OTel endpoint and headers.
- Terraform >= `1.0`.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project
- OTel Access Setup and Authorization.
- Log Integration configuration.



## Usage

**1\. Ensure your Datadog and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export OTEL_ENDPOINT='<OTEL-ENDPOINT-URL>'
export OTEL-SUPPLIED-HEADERS='<OTEL-ENDPOINT-HEADERS>'
```


... or follow as in the `~/.otel/variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
project_id            = "your-mongodb-project-id"
type                  = "DATADOG_LOG_EXPORT"
log_types             = "[your-log-export-types]"
otel_endpoint         = "your-otel-endpoint-url"
otel_supplied_headers = {
    name  = "your-otel-endpoint-header-name"
    value = "your-otel-endpoint-header-value"
}
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

