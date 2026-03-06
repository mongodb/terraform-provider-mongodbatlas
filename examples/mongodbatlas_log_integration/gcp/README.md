# MongoDB Atlas Log Integration with Google Cloud Storage Buckets Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to a Google Cloud Storage (GCS) bucket.

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- GCP account with permissions to create GCS buckets and manage IAM bindings.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project.
- Cloud Provider Access Setup and Authorization.
- Log Integration configuration.

### GCP
- GCS bucket for storing logs.
- IAM binding granting the Atlas-managed service account object admin access to the bucket.

## Usage

**1\. Ensure your GCP and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
atlas_org_id        = "your-org-id"
atlas_client_id     = "your-service-account-client-id"
atlas_client_secret = "your-service-account-client-secret"
gcp_project_id      = "your-gcp-project-id"
gcs_bucket_name     = "your-globally-unique-bucket-name"
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
