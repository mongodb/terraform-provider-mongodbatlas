# MongoDB Atlas Log Integration Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to an Amazon S3 bucket.

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role
- AWS account with permissions to create S3 buckets and IAM roles
- Terraform >= 1.0

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project
- Cloud Provider Access Setup and Authorization
- Log Integration configuration

### AWS
- S3 bucket for storing logs
- IAM role for Atlas to assume
- IAM policy for S3 access

## Usage

1. Set the required variables in a `terraform.tfvars` file or via environment variables:

```hcl
atlas_org_id        = "your-org-id"
atlas_client_id     = "your-service-account-client-id"
atlas_client_secret = "your-service-account-client-secret"
access_key          = "your-aws-access-key"
secret_key          = "your-aws-secret-key"
```

2. Initialize and apply:

```bash
terraform init
terraform apply
```

## Log Types

The `log_types` attribute supports the following values:
- `MONGOD` - MongoDB server logs
- `MONGOS` - MongoDB router logs
- `MONGOD_AUDIT` - MongoDB server audit logs
- `MONGOS_AUDIT` - MongoDB router audit logs

## Notes

- The requesting Service Account or API Key must have the Organization Owner or Project Owner role
- MongoDB Cloud will add sub-directories based on the log type under the specified `prefix_path`
- Optional: Use `kms_key` to specify an AWS KMS key ID or ARN for server-side encryption

