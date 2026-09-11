# Service Account with a bootstrap secret

Create a Service Account and let Atlas generate its first secret by setting `secret_expires_after_hours`.

The secret value is returned only once, at creation time. The example exposes it through the sensitive output `service_account_first_secret`.

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

**4. Read the initial secret** (**warning**: this prints the secret to your terminal).

```bash
terraform output -raw service_account_first_secret
```

For secret rotation, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

**5. Destroy.**

```bash
terraform destroy
```
