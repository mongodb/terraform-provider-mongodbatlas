# MongoDB Atlas Organization Log Integration with OpenTelemetry Example

This example demonstrates how to configure an organization-level log integration to export MongoDB Atlas organization events to an OpenTelemetry (OTel) collector endpoint.

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner role.
- An OpenTelemetry collector accessible over HTTPS that accepts event data via the OTLP HTTP protocol.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Organization Log Integration configuration.

It also shows how to use the `mongodbatlas_org_log_integration` and `mongodbatlas_org_log_integrations` data sources to read the configuration (see `singular-data-source.tf` and `plural-data-source.tf`).

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
otel_endpoint       = "https://your-otel-collector.com:4318/v1/logs"
otel_supplied_headers = [
  {
    name  = "Authorization"
    value = "Bearer your-token"
  }
]
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

When you have finished your testing, ensure you destroy the resources.

```bash
terraform destroy
```

## Log Types

The `log_types` attribute supports the following values:
- `EVENTS` - Organization events.
