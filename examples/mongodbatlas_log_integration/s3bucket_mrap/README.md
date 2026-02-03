# MongoDB Atlas Log Integration with S3 Multi-Region Access Point (MRAP) Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to an AWS S3 Multi-Region Access Point (MRAP) instead of a single S3 bucket.

## Overview

Using a Multi-Region Access Point (MRAP) provides:
- **High availability**: Automatic routing to the closest available S3 bucket
- **Lower latency**: Requests are routed to the nearest regional bucket
- **Resilience**: Continues to work even if one region becomes unavailable

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role.
- AWS account with permissions to create S3 buckets, S3 Multi-Region Access Points, and IAM roles.
- Terraform >= `1.0`.
- AWS Provider >= `5.0`.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Cloud Provider Access Setup and Authorization.
- Log Integration configuration using MRAP.

### AWS
- S3 buckets in multiple regions (us-east-1 and us-west-2).
- S3 Multi-Region Access Point (MRAP).
- IAM role for Atlas to assume.
- IAM policy with permissions for MRAP and backing buckets.

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export AWS_ACCESS_KEY_ID='<AWS_ACCESS_KEY_ID>'
export AWS_SECRET_ACCESS_KEY='<AWS_SECRET_ACCESS_KEY>'
```

... or the `~/.aws/credentials` file.

```
$ cat ~/.aws/credentials
[default]
aws_access_key_id = <AWS_ACCESS_KEY_ID>
aws_secret_access_key = <AWS_SECRET_ACCESS_KEY>
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
project_id  = "your-atlas-project-id"
name_prefix = "my-atlas-logs"
prefix_path = "atlas-logs/"
log_types   = ["MONGOD", "MONGOS"]
```

**2\. Review the Terraform plan.**

Execute the following command and ensure you are happy with the plan.

```bash
terraform plan
```

**3\. Execute the Terraform apply.**

Now execute the plan to provision the resources.

```bash
terraform apply
```

**4\. Destroy the resources.**

When you have finished your testing, ensure you destroy the resources to avoid unnecessary Atlas and AWS charges.

```bash
terraform destroy
```

## Log Types

The `log_types` attribute supports the following values:
- `MONGOD` - MongoDB server logs.
- `MONGOS` - MongoDB router logs.
- `MONGOD_AUDIT` - MongoDB server audit logs.
- `MONGOS_AUDIT` - MongoDB router audit logs.

## Variables

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `project_id` | MongoDB Atlas Project ID | `string` | (required) |
| `name_prefix` | Prefix for naming AWS resources (must be globally unique for S3) | `string` | `"atlas-logs"` |
| `prefix_path` | S3 directory path prefix for log files | `string` | `"atlas-logs/"` |
| `log_types` | Array of log types to export | `list(string)` | `["MONGOD", "MONGOS"]` |

## Outputs

| Name | Description |
|------|-------------|
| `mrap_alias` | The MRAP alias used in bucket_name |
| `mrap_arn` | The MRAP ARN |
| `integration_id` | The ID of the log integration |
| `iam_role_arn` | ARN of the IAM role used for the integration |
| `cloud_provider_role_id` | Atlas Cloud Provider Access Role ID |
| `backing_buckets` | The S3 buckets backing the MRAP |

## Notes

- The requesting Service Account or API Key must have the Organization Owner or Project Owner role.
- MongoDB Cloud will add sub-directories based on the log type under the specified `prefix_path`.
- MRAP creation can take several minutes to complete.
- The IAM policy includes permissions for both the MRAP and the backing S3 buckets, as some operations require direct bucket access.
- Optional: Use `kms_key` to specify an AWS KMS key ID or ARN for server-side encryption.
