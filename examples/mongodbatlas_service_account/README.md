# MongoDB Atlas Provider -- Service Account

This example shows how to create a Service Account without an Atlas-generated secret by setting `without_initial_secret = true`, then create its first secret through `mongodbatlas_service_account_secret`.

## Important Notes

Setting `without_initial_secret = true` is the preferred approach: Atlas returns no secret from the create request, so you manage it through `mongodbatlas_service_account_secret`. It is mutually exclusive with `secret_expires_after_hours`, and Atlas rejects a create request that sets both.

The example includes a sensitive output `secret` that captures the secret value. You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw secret
```

For managing and rotating secrets, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

## Prerequisites

- Service Account with Organization Owner permissions used for provider authentication.

## Variables Required to be set

- `atlas_client_id`: MongoDB Atlas Service Account Client ID.
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret.
- `org_id`: Atlas Organization ID where this configuration creates the Service Account.

## Outputs

- `service_account_client_id`: The Client ID of the Service Account
- `service_account_name`: The name of the Service Account, read from the data source
- `secret_id`: The ID of the Service Account secret
- `secret` (sensitive): The secret value
- `service_accounts_results`: All Service Accounts in the organization

## Usage

**1. Create `terraform.tfvars`.**

```hcl
atlas_client_id     = "<ATLAS_CLIENT_ID>"
atlas_client_secret = "<ATLAS_CLIENT_SECRET>"
org_id              = "your-org-id"
```

**2. Plan and apply.**

```bash
terraform plan
terraform apply
```

**3. Read the secret** (**warning**: this prints the secret to your terminal).

```bash
terraform output -raw secret
```

**4. Destroy.**

```bash
terraform destroy
```
