# Service Account with a bootstrap secret

Create a Service Account and let Atlas generate its first secret on create, by setting `secret_expires_after_hours`.

Atlas returns the secret value only once, in the create response. The example captures it in a sensitive output, so it is available immediately after the first apply and never again.

## Prerequisites

- `org_id` for the organization that gets the Service Account.
- Organization Owner role for the credentials that run the apply.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Defaults

The example sets `secret_expires_after_hours = 2160` (90 days). Adjust it in `main.tf` to match your organization's secret lifetime policy.

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

**4. Read the bootstrap secret.**

Warning: this prints the secret to your terminal.

```bash
terraform output -raw service_account_first_secret
```

## Outputs

- `service_account_client_id`: The Client ID of the created Service Account
- `service_account_name`: The name of the Service Account
- `service_account_first_secret` (sensitive): The initial secret value, available only at creation
- `service_accounts_results`: All Service Accounts in the organization

## Notes

- `secret_expires_after_hours` cannot be updated after creation. Changing it forces a new Service Account.
- To rotate the secret instead of recreating the account, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).
