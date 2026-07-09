# MongoDB Atlas Metric Integration with OpenTelemetry Example

This example demonstrates how to configure a metric integration to export MongoDB Atlas metrics to an OTLP-compatible endpoint such as Datadog, New Relic, or Dynatrace. It also shows how to read the integration back with the singular and plural data sources.

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- An OTLP-compatible endpoint URL and authentication credentials. This example uses Datadog and creates the Datadog API key with the `datadog` provider.

## Resources Created

This example creates the following resources:

### MongoDB Atlas

- Project.
- Metric Integration configuration.

### Datadog

- API key used to authenticate metric ingestion.

## Usage

**1\. Ensure your MongoDB Atlas and Datadog credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
export DD_API_KEY="<DATADOG_API_KEY>"
export DD_APP_KEY="<DATADOG_APP_KEY>"
```

... or follow as in the `variables.tf` file and create a **terraform.tfvars** file with all the variable values:

```hcl
atlas_org_id        = "your-org-id"
atlas_client_id     = "your-service-account-client-id"
atlas_client_secret = "your-service-account-client-secret"
datadog_api_key     = "your-datadog-api-key"
datadog_app_key     = "your-datadog-app-key"
datadog_endpoint    = "https://otlp.datadoghq.com/v1/metrics"
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
