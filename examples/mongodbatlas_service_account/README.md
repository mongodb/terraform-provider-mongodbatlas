# MongoDB Atlas Provider -- Service Account

This example shows how to create a Service Account without an Atlas-generated secret by setting `without_initial_secret = true`, then create its first secret as a managed resource.

Setting `without_initial_secret = true` means Atlas returns no secret from the create request, so every secret is created and rotated through `mongodbatlas_service_account_secret`, which is what a rotation submodule expects.

## Important Notes

`without_initial_secret` and `secret_expires_after_hours` are mutually exclusive on the Service Account. Set `without_initial_secret = true` and omit `secret_expires_after_hours`, or set `secret_expires_after_hours` and omit `without_initial_secret`. Atlas rejects a create request that sets both. This example sets the expiration only on the secret resource.

The example includes a sensitive output `secret` that captures the secret value. You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw secret
```

For managing and rotating secrets, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

## Prerequisites

- Service Account with Organization Owner permissions used for provider authentication.
- Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Variables Required to be set

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Atlas Organization ID where the Service Account is created

## Outputs

- `service_account_client_id`: The Client ID of the Service Account
- `secret_id`: The ID of the Service Account secret
- `secret` (sensitive): The secret value

## Usage

**1. Create `terraform.tfvars`.**

```hcl
org_id = "your-org-id"
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
