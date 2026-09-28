---
page_title: "Guide: Service Account Secret Rotation"
---

# Guide: Service Account Secret Rotation

**Objective**: This guide shows two rotation models for Service Account secrets, and how to choose between them. The choice depends on whether credential consumers outside Terraform need a handoff window while the secret changes.

**Scope**: This guide starts from a new Service Account, created with `without_initial_secret = true`. You cannot set that attribute on an existing Service Account, so use these models for accounts this Terraform stack creates.

## Overview

Service Account secrets expire after the `secret_expires_after_hours` you set, anywhere from 8 hours to 365 days. Terraform cannot update a secret in place, so rotation recreates the `mongodbatlas_service_account_secret` resource with `terraform apply -replace`. Atlas returns the new secret value only once, at creation time.

Two terms used throughout this guide:

- **Slot**: One secret resource. A Service Account and an MCP configuration each hold at most two.
- **Consumer**: Anything that authenticates with the secret, such as an application, a CI job, or another Terraform stack.

This guide applies to both organization-level and project-level service accounts:

- **Organization-level**: Use `mongodbatlas_service_account` and `mongodbatlas_service_account_secret`
- **Project-level**: Use `mongodbatlas_project_service_account` and `mongodbatlas_project_service_account_secret`

**Note**: The steps below use organization-level resources, but the same approach applies to project-level resources.

~> **WARNING:** Service Account secrets expire after the configured `secret_expires_after_hours` period. To avoid losing access to the Atlas Administration API, update your application with the new client secret as soon as possible after rotation. If all secrets expire before being replaced, you will lose access to the organization. For more information, see [Rotate Service Account Secrets](https://www.mongodb.com/docs/atlas/tutorial/rotate-service-account-secrets/).


## Security notes

- Managing Service Accounts with Terraform **exposes sensitive organizational secrets** in Terraform's state. Follow [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data).
- `terraform output -raw` prints the secret value to your terminal.

## Choose a model

- **Single-secret rotation (Model A)**: Use this when no external consumer holds the secret, or a consumer can switch immediately. One secret resource with `create_before_destroy`. This is the simplest path and has no overlap window.
- **Two-slot rotation (Model B)**: Use this when an external consumer holds the secret and needs time to roll over. The account carries two secrets, and rotating one leaves the other valid while consumers switch. Adds an extra resource and a two-cycle schedule.

Both models set `without_initial_secret = true` on the Service Account. Atlas then creates no secret at create time, every secret is managed by `mongodbatlas_service_account_secret`, and there is no import step.

Do not configure the provider to authenticate with the secret that the rotation replaces, unless you use `lifecycle { create_before_destroy = true }` (see Model A). Atlas revokes an OAuth token as soon as the secret that minted it is deleted, so a plain replace, which destroys before it creates, fails with HTTP 401. With `create_before_destroy`, the create runs while the old secret still exists and the provider's token stays valid.

## The 7-day overlap limit

Creating a new secret shortens every existing secret on the same Service Account. Atlas documents this in [Rotate Service Account Secrets](https://www.mongodb.com/docs/atlas/tutorial/rotate-service-account-secrets/):

- Generating a new secret shortens each existing secret to the shorter of its remaining lifetime and **7 days**.
- The new secret keeps the `secret_expires_after_hours` you requested.
- A Service Account holds at most two secrets, and expired secrets are not deleted automatically.

The overlap between the old and new secret is therefore at most 7 days, no matter what `secret_expires_after_hours` says. Two secrets configured with `secret_expires_after_hours = 2160` do not give you a 90-day overlap. Plan every consumer handoff inside that 7-day window.

This limit only applies to Model B, which carries two secrets. Model A replaces its single secret and has no overlap. In Model B, the secret you do not replace still triggers the Atlas `Service Account Secrets are about to expire` alert each cycle. Replace it on schedule or accept the notification.

## Model A: Single-secret rotation

Use this model when no external consumer holds the secret, when a consumer can switch immediately, or when this Service Account authenticates the Terraform stack itself.

### Configuration

```terraform
variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

# Create the Service Account without an Atlas-generated secret.
resource "mongodbatlas_service_account" "this" {
  org_id                 = var.org_id
  name                   = "example-service-account"
  description            = "Example Service Account"
  roles                  = ["ORG_READ_ONLY"]
  without_initial_secret = true
}

resource "mongodbatlas_service_account_secret" "this" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = 2160 # 90 days

  # Create the new secret before the old one is deleted.
  lifecycle {
    create_before_destroy = true
  }
}

output "secret" {
  sensitive = true
  value     = mongodbatlas_service_account_secret.this.secret
}
```

`ORG_READ_ONLY` works when the apply runs with an admin credential. When the provider authenticates as this Service Account, its role must allow managing its own secrets, such as `ORG_OWNER`.

### Rotate

1. Replace the secret:

```shell
terraform apply -replace="mongodbatlas_service_account_secret.this"
```

2. Retrieve and securely store the new secret value (**warning**: this prints the secret to your terminal):

```shell
terraform output -raw secret
```

3. Update every consumer with the new value. A consumer that can switch immediately is the right fit for this model.

### Why `create_before_destroy` matters

There is no overlap window. Terraform deletes the old secret after it creates the new one, and a consumer can keep authenticating with the old secret only until that destroy completes.

That ordering is also what makes this configuration safe when the provider authenticates with the secret it rotates. A plain replace destroys the old secret first, which revokes the OAuth token minted from it, and the create that follows fails with HTTP 401. `create_before_destroy` runs the create first, so the token stays valid and the apply completes.

When the provider authenticates with a different credential, `create_before_destroy` is optional.

## Model B: Two-slot rotation

Use this model when an application, a CI job, or another Terraform stack holds the secret and must keep authenticating through the switch.

The account carries two secrets. Rotating one leaves the other valid while consumers switch.

### Configuration

The [two-slot rotation example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_service_account_secret_rotation) is the full configuration. It creates the account with `without_initial_secret = true`, manages both slots, and reads the live expiry of each slot through `data "mongodbatlas_service_account"`. It exposes `current_credentials`, the credential consumers adopt after a rotation, resolved to the slot with the largest `expires_at`, and `expires_at` for both slots.

Two points to carry over when you build your own:

- `secret_2` needs `depends_on = [mongodbatlas_service_account_secret.secret_1]`. Without it, Terraform creates the two secrets in parallel and the second POST can fail with:

```text
HTTP 409 Conflict (Error code: "DATA_CONCURRENCY_ERROR")
```

- Read each slot's expiry from the data source, not from the resource. A resource `expires_at` refreshes only when that resource is read, so the slot that is not rotating keeps a stale value in state.

Do not set `secret_expires_after_hours` on the Service Account when `without_initial_secret = true`. Atlas rejects a create request that sets both.

### How consumers map to slots

Keep each consumer on the slot you are not rotating, and move it to the new secret after the replace:

- Cycle 1: Consumers authenticate with `secret_2`. Replace `secret_1`, then move consumers to the new `secret_1`.
- Cycle 2: Consumers authenticate with `secret_1`. Replace `secret_2`, then move consumers back.

Alternating the slot that rotates keeps one live secret for consumers while the other is replaced and deployed.

### Rotate

1. Confirm that consumers currently authenticate with the other slot. For the first cycle, that is `secret_2`.
2. Replace the slot:

```shell
terraform apply -replace="mongodbatlas_service_account_secret.secret_1"
```

3. Read the credential to deploy and the deadline (**warning**: this prints the secret to your terminal):

```shell
terraform output -json current_credentials
terraform output -json expires_at
```

`current_credentials` resolves to the slot with the largest `expires_at`, so it is the new secret. Deploy it as-is, without reconstructing the value from `secret_1` or `secret_2`. `expires_at` shows how long each slot has left, including the 7-day cut applied to the slot you did not replace. Use the non-rotated slot's `expires_at` as the handoff deadline.

4. Roll every consumer over to the new value before that deadline.
   - Update the stored credential in each consumer with the `client_secret` from `current_credentials`, such as a CI secret, a secret manager entry, or an environment variable.
   - Redeploy or restart the consumer.
   - Verify that it authenticates before the window closes.
5. For the next cycle, replace `secret_2` instead and roll consumers back.

Replacing `secret_1` shortens `secret_2` to 7 days from the creation of the new `secret_1`, so roll consumers over before `secret_2` expires.

Notes for this model:

- `client_id` stays constant through rotation; only the secret changes.
- Terraform deletes the replaced secret as part of the apply, so you do not revoke it manually.
- The account needs a role that can manage its own secrets only when you rotate while authenticated as that account. Otherwise the admin credential that runs the apply needs it.

## MCP configuration secrets

MCP configuration secrets rotate with the same two-slot pattern. An MCP configuration holds at most two ingress secrets, so define two secret resources and replace them alternately.

The same pattern applies to `mongodbatlas_project_mcp_config_secret`, with `project_id` in place of `org_id`.

Rotate by replacing one slot at a time, alternating between them:

```shell
terraform apply -replace="mongodbatlas_mcp_config_secret.secret_1"
terraform output -raw secret_1
```

The secret value is returned only in the create response, so read the output right after the replace. The snippet assumes an output named `secret_1` for the first slot's value; define one output per slot. The 7-day overlap limit and the consumer handoff steps from [Model B](#model-b-two-slot-rotation) apply here as well.

## Related documentation

- [Two-slot rotation example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_service_account_secret_rotation)
- [`mongodbatlas_service_account`](../resources/service_account)
- [`mongodbatlas_service_account_secret`](../resources/service_account_secret)
- [`mongodbatlas_project_service_account`](../resources/project_service_account)
- [`mongodbatlas_project_service_account_secret`](../resources/project_service_account_secret)
- [`mongodbatlas_mcp_config_secret`](../resources/mcp_config_secret)
- [`mongodbatlas_project_mcp_config_secret`](../resources/project_mcp_config_secret)
