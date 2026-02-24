# MongoDB Atlas Log Integration with Google Cloud Platform Example

This example demonstrates how to configure a log integration to export MongoDB Atlas logs to a Microsoft Azure Blob storage.

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role.
- Google Cloud Platform account with permissions to create Containers and IAM roles.
- Terraform >= `1.0`.

## Resources Created

This example creates the following resources:

### MongoDB Atlas
- Project.
- Cloud Provider Access Setup and Authorization.
- Log Integration configuration.

### Google CLoud Platform
- GCP Container for storing logs.
- IAM role for Atlas to assume.
- IAM policy for Container access.


## Usage

**1\. Ensure your Google Cloud Platform and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

```bash
export GCP_ACCESS_KEY_ID='<GCP_ACCESS_KEY_ID>'
export GCP_SECRET_ACCESS_KEY='<GCP_SECRET_ACCESS_KEY>'
```

... or the `~/.gcp/credentials` file.

```
$ cat ~/.gcp/credentials
[default]
gcp_access_key_id = <GCP_ACCESS_KEY_ID>
gcp_secret_access_key = <GCP_SECRET_ACCESS_KEY>
```

... or follow as in the `~/.gcp/variables.tf` file and create **terraform.tfvars** file with all the variable values:

```hcl
atlas_org_id        = "your-org-id"
atlas_client_id     = "your-service-account-client-id"
atlas_client_secret = "your-service-account-client-secret"
access_key          = "your-gcp-access-key"
secret_key          = "your-gcp-secret-key"
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

- The requesting Service Account or API Key must have the Organization Owner or Project Owner role.
- MongoDB Atlas will add sub-directories based on the log type under the specified `prefix_path`.
