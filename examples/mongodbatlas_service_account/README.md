# MongoDB Atlas Provider -- Service Account

This example shows how to create a Service Account without an Atlas-generated secret by setting `without_initial_secret = true`.

Create the secrets for the Service Account separately with the [`mongodbatlas_service_account_secret`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/service_account_secret) resource, so this configuration owns their lifecycle. See [`mongodbatlas_service_account_secret`](../mongodbatlas_service_account_secret/README.md) for a complete example.

Setting `without_initial_secret = true` means Atlas returns no secret from the create request. Every secret is created and rotated through `mongodbatlas_service_account_secret`, which is what a rotation submodule expects.

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

**3. Destroy.**

```bash
terraform destroy
```
