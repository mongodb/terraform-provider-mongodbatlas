# MongoDB Atlas Log Integration with Datadog

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to a Splunk integration.

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role.
- Splunk HTTP Event Collector Token
- Splunk Event Collector URL
- Terraform >= `1.0`.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project
- Splunk HEC Access Setup and Authorization.
- Log Integration configuration.



## Usage

**1\. Ensure your Datadog and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export HEC_TOKEC='<SPLUNK-HEC-TOKEN>'
export HEC_URL='<SPLUNK-HEC-URL>'
```


... or follow as in the `~/.splunk/variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
project_id    = "your-mongodb-project-id"
type          = "DATADOG_LOG_EXPORT"
log_types     = "[your-log-export-types]"
hec_token     = "your-splunk-hec-token"
hec_url       = "your-splunk-hec-url"
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
- `MONGOD_AUDIT` - MongoDB server audit logs.

