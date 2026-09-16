# Service Account without an initial secret

Create a Service Account and let Atlas skip secret generation, by omitting `secret_expires_after_hours`.

Use this flow when something else owns the secret: a rotation submodule, an external secret store, or a separate `mongodbatlas_service_account_secret` resource. The Service Account is created, but the `secrets` list stays empty.

## Prerequisites

- `org_id` for the organization that gets the Service Account.
- Organization Owner role for the credentials that run the apply.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Usage

**1. Set credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**2. Create `terraform.tfvars`.**

```hcl
org_id = "your-org-id"
```

**3. Apply.**

```bash
terraform init
terraform apply
```

## Outputs

- `service_account_client_id`: The Client ID of the created Service Account
- `service_account_name`: The name of the Service Account
- `service_account_secrets`: Empty after create, since no initial secret is generated
- `service_accounts_results`: All Service Accounts in the organization

## Adding a secret later

Terraform cannot see a secret created outside its state, so no drift is reported if one is added in the UI or API. To manage the secret in Terraform, add a [`mongodbatlas_service_account_secret`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/service_account_secret) resource that references `service_account_client_id`.

Adding a secret to an account that already has one cuts the existing secret to at most 7 days. See the [service account rotation guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation) before rotating.
