# Two-slot Service Account secret rotation

Rotate a Service Account secret while a consumer keeps authenticating: the account carries two secrets, and rotating one leaves the other valid through the switch.

For the model comparison and the expiry behavior, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).
For product limits, see [Rotate Service Account Secrets](https://www.mongodb.com/docs/atlas/tutorial/rotate-service-account-secrets/).

## How it works

A Service Account holds at most two secrets. This example manages both:

- `secret_1` and `secret_2` are `mongodbatlas_service_account_secret` resources on one `mongodbatlas_service_account` created with `without_initial_secret = true`. No import step.
- `secret_2` depends on `secret_1` so the two creates run one after the other. Without that ordering the second POST can fail with `DATA_CONCURRENCY_ERROR` (HTTP 409).
- `current_credentials` is the credential consumers adopt after a rotation. It resolves to the slot with the largest live `expires_at`.
- `expires_at` reports the live expiry of both slots, read through `data.mongodbatlas_service_account`.

## The 7-day overlap

Creating a secret cuts every existing secret on the same Service Account to the shorter of its remaining lifetime and 7 days. The new secret keeps the requested `secret_expires_after_hours`. Two secrets requested at 90 days do not give a 90-day overlap; the real overlap is at most 7 days.

After apply, `expires_at` shows this: the slot created second expires in about 90 days, and the slot created first drops to about 7 days.

## Rotate

Rotate the slot consumers are not using, then roll them onto the new value.

**1. Confirm which slot consumers use.** `current_credentials` resolves to it. On a fresh apply that is `secret_2`.

**2. Replace the other slot.**

```bash
terraform apply -replace="mongodbatlas_service_account_secret.secret_1"
```

**3. Read the new credential.**

```bash
terraform output -json current_credentials
terraform output -json expires_at
```

**4. Roll consumers onto the new value within 7 days.** Update the stored credential in each consumer, redeploy or restart it, and verify it authenticates before the window closes.

**5. Next cycle, replace `secret_2` and roll consumers back.**

`client_id` stays constant through rotation; only the secret changes. Terraform deletes the replaced secret as part of the apply.

## Prerequisites

- Credentials with Organization Owner permissions, exported as `MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET` for the first apply.
- `without_initial_secret` support, added in provider v2.19.0.

## Variables

- `org_id` (required): Atlas Organization ID.
- `service_account_name`: Name of the Service Account to create. Defaults to `example-rotation-service-account`.
- `service_account_roles`: Roles for the Service Account. Defaults to `["ORG_OWNER"]`, which lets the account rotate its own secrets.
- `secret_expires_after_hours`: Requested lifetime of each secret in hours. Defaults to `2160` (90 days).

## Usage

**1. Set credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
export TF_VAR_org_id="<ATLAS_ORG_ID>"
```

**2. Plan and apply.**

```bash
terraform init
terraform plan
terraform apply
```

**3. Read the outputs.**

```bash
terraform output expires_at
terraform output -json current_credentials
```

**4. Rotate.** Replace the slot consumers are not using, then roll them over. See Rotate above. To rotate while authenticated as the Service Account itself, run the replace with `MONGODB_ATLAS_CLIENT_ID` set to its `client_id` and `MONGODB_ATLAS_CLIENT_SECRET` set to the secret of the slot you are not replacing.

**5. Destroy.**

```bash
terraform destroy
```
