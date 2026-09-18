# Service Account without an initial secret

Create a Service Account with `without_initial_secret = true`, then create its first managed secret with `mongodbatlas_service_account_secret`.

Use this flow when you want to manage the secrets for a Service Account yourself, for example when a rotation submodule owns the secret lifecycle. Atlas does not generate a bootstrap secret, so nothing is returned that the configuration cannot manage.

This is the recommended create flow. The alternative is to set `secret_expires_after_hours` and let the create request generate a bootstrap secret, which you can only read from the create response and rotate later.

## Prerequisites

- Service Account with Organization Owner permissions used for provider authentication.
- Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Variables

- `org_id`: Atlas Organization ID where the Service Account is created.

## Usage

**1. Set credentials and the organization ID.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**2. Create `terraform.tfvars`.**

```hcl
org_id = "your-org-id"
```

**3. Plan and apply.**

```bash
terraform plan
terraform apply
```

`without_initial_secret` and `secret_expires_after_hours` are mutually exclusive. Set `without_initial_secret = true` and omit `secret_expires_after_hours`, or set `secret_expires_after_hours` and omit `without_initial_secret`. Atlas rejects a create request that sets both.

**4. Read the secret.**

```bash
terraform output -raw secret
```

For managing and rotating secrets, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

**5. Destroy.**

```bash
terraform destroy
```
