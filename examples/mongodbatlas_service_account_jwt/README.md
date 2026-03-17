# MongoDB Atlas Provider -- Ephemeral Service Account JWT with AWS Secrets Manager

This example demonstrates how to generate a short-lived Atlas JWT using the `mongodbatlas_service_account_jwt` ephemeral resource and store it securely in AWS Secrets Manager. A second configuration then retrieves the stored token and uses it to authenticate the Atlas provider and create a project.

In `step-1-token-generator`, the JWT is never written to Terraform state or plan. The token is persisted in AWS Secrets Manager using a write-only attribute (`secret_string_wo`), so the value is sent to AWS but excluded from Terraform state. In `step-2-token-consumer`, the JWT is read from Secrets Manager via a data source, which means it will be present in step 2's Terraform state (marked as sensitive).

## Prerequisites

- Terraform >= 1.11 (required for write-only attributes in step 1).
- A MongoDB Atlas Service Account. By default, the JWT is generated using the provider's SA credentials. To use a dedicated SA with specific permissions, see the inline comments in `step-1-token-generator/main.tf`.
- AWS CLI configured with credentials that have `secretsmanager:CreateSecret`, `secretsmanager:PutSecretValue`, and `secretsmanager:GetSecretValue` permissions.

## Structure

| Directory | Purpose |
|---|---|
| `step-1-token-generator/` | Generates an ephemeral JWT and stores it in AWS Secrets Manager. |
| `step-2-token-consumer/` | Reads the JWT from Secrets Manager, configures the Atlas provider with it, and creates a project. |

## Usage

### Step 1: Generate and store the token

```bash
cd step-1-token-generator
```

Set the required variables in `terraform.tfvars`:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID.
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret.

Then apply:

```bash
terraform init
terraform apply
```

Note the `aws_secret_id` output, you will need it for step 2.

### Step 2: Consume the token

```bash
cd ../step-2-token-consumer
```

Set the required variables in `terraform.tfvars`:

- `aws_secret_id`: ARN of the AWS Secrets Manager secret (from step 1 output).
- `org_id`: Organization ID where the project will be created.

Then apply:

```bash
terraform init
terraform apply
```

This reads the JWT from Secrets Manager, authenticates the Atlas provider using `access_token`, and creates a project.

## Alternative: local-exec provisioner (Terraform >= 1.10)

If you are on Terraform 1.10 or your cloud provider does not yet support write-only attributes, see the inline comments in `step-1-token-generator/main.tf` for instructions on switching to a `local-exec` provisioner approach.

## Cleanup

Destroy resources in reverse order:

```bash
cd step-2-token-consumer
terraform destroy

cd ../step-1-token-generator
terraform destroy
```

The AWS Secrets Manager secret has a default 30-day recovery window. To delete it immediately, use:

```bash
aws secretsmanager delete-secret --secret-id <aws_secret_id> --force-delete-without-recovery
```
