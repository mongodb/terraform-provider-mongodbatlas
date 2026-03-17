# MongoDB Atlas Provider -- Ephemeral Service Account JWT with AWS Secrets Manager

This example demonstrates how to generate a short-lived Atlas JWT using the `mongodbatlas_service_account_jwt` ephemeral resource and store it securely in AWS Secrets Manager. A second configuration then retrieves the stored token and uses it to authenticate the Atlas provider and create a project. This is useful when passing short-lived credentials between separate Terraform configurations, CI/CD pipeline stages, or teams, for example, a platform team generating a JWT for an application team, or a pipeline injecting a scoped token into downstream stages that provision Atlas resources.

In `step-1-token-generator`, the JWT is never written to Terraform state or plan. The token is persisted in AWS Secrets Manager using a write-only attribute (`secret_string_wo`), so the value is sent to AWS but excluded from Terraform state. In `step-2-token-consumer`, the JWT is read from Secrets Manager via a data source.

## Prerequisites

- Terraform >= 1.11 (required for write-only attributes in step 1).
- A MongoDB Atlas Service Account.
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

## Using a dedicated Service Account for the JWT

By default, the ephemeral resource generates a JWT using the provider's Service Account credentials. To generate a JWT with a different access level, create a dedicated Service Account and pass its credentials explicitly:

```hcl
resource "mongodbatlas_service_account" "jwt_sa" {
  org_id                     = var.org_id
  name                       = "jwt-dedicated-sa"
  description                = "SA used exclusively for ephemeral JWT generation."
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160
}

resource "mongodbatlas_service_account_secret" "jwt_sa" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.jwt_sa.client_id
  secret_expires_after_hours = 2160
}

ephemeral "mongodbatlas_service_account_jwt" "token" {
  client_id     = mongodbatlas_service_account.jwt_sa.client_id
  client_secret = mongodbatlas_service_account_secret.jwt_sa.secret
}
```

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

