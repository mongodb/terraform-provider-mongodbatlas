# organization2 PoC (`mongodbatlas_organization2`)

Mock-backed Terraform Plugin Framework resource that demonstrates **ModifyPlan**-driven client secret rotation with an optional `client_secret_rotation` block. It is an **anti-pattern reference** for embedded org credentials; production Atlas org design should use child `mongodbatlas_service_account` and `mongodbatlas_service_account_secret` resources instead.

## Purpose

- Show that computed `next_renewal` / `expires_at` alone do not schedule rotation; the rotation **block** opts into ModifyPlan.
- Show plan shape when rotation is due or forced: known `secret_version` increment, known `old_secret_id`, unknown secret fields.
- Show practitioner-forced rotation via a higher `secret_version` in config before the calendar deadline.

## Example

```hcl
resource "mongodbatlas_organization2" "demo" {
  name = "demo-org"

  client_secret_rotation = {
    interval = "240h"
    # secret_version = 2  # optional: force rotation before next_renewal
  }
}
```

Without the `client_secret_rotation` block, the resource still creates mock `client_id` / `client_secret` on create, but no scheduled rotation runs.

## Production pattern (do not copy this resource)

- Keep org identity and settings on `mongodbatlas_organization`.
- Create credentials with `mongodbatlas_service_account` and `mongodbatlas_service_account_secret`.
- Use two secret resources for overlap, or `terraform apply -replace` / optional `force_renew` on the secret resource.
- Use `time_rotating` + `replace_triggered_by` for calendar-driven rotation without provider clocks.

## Limitations

- Mock store persisted to a local JSON file (default `$TMPDIR/mongodbatlas-organization2-poc-store.json`; override with `MONGODB_ATLAS_ORGANIZATION2_POC_STORE`); no Atlas API calls.
- Always registered on the muxed TPF provider (experimental).
- `old_secret_id` in state is teaching-only overlap visibility; production uses two secret resources.
- Removing `client_secret_rotation` from config stops future ModifyPlan rotation; `old_secret_id` may remain stale in state.
- No registry docs, examples, or CHANGELOG entry for this PoC.

## Mock store file

Each provider process loads the mock store from disk on first access. Set `MONGODB_ATLAS_ORGANIZATION2_POC_STORE` to pin a path when debugging (for example next to your Terraform working directory).

```sh
export MONGODB_ATLAS_ORGANIZATION2_POC_STORE=/tmp/mongodbatlas-organization2-demo.json
terraform apply
terraform apply  # Read finds the prior create; plan stays empty when not due
```

## Tests

Unit tests (no Atlas credentials):

```sh
cd code/provider
go test ./internal/service/organization2/ -run 'TestRotationDue|TestModifyPlan' -v
```

Acceptance tests (mock backend; no Atlas env vars required):

```sh
cd code/provider
TF_ACC=1 go test ./internal/service/organization2/ -run TestAccOrganization2 -v
```

Acceptance coverage:

- `TestAccOrganization2_noRotationBlock` — create with `name` only; no rotation attributes in state.
- `TestAccOrganization2_withRotationBlock` — `interval = "2s"`; sleep past due; `secret_version` increments and `old_secret_id` matches prior `current_secret_id`.
- `TestAccOrganization2_forceSecretVersion` — long `interval`; set `secret_version = 2` before due without sleep.

## Branch

Implemented on `CLOUDP-381539_org_resource_sa_rotation_support` until maintainer review.

---

# organization3 PoC (`mongodbatlas_organization3`)

Real-Atlas Terraform Plugin Framework resource that demonstrates **ModifyPlan**-driven client secret rotation with API-backed `expires_at` scheduling. Compare with the mock `mongodbatlas_organization2` PoC above.

## Schema vs organization2

- **Scheduling**: `rotate_before_expiry_hours` (default `expires_after_hours / 2`) and Atlas `current_secret.expires_at`; no `interval`, `next_renewal`, or provider-local `expires_at`.
- **Secret metadata**: Nested `current_secret` and `old_secret` objects (`secret_id`, `created_at`, `expires_at`) refreshed on read. `last_used_at` is omitted to avoid plan noise when the SA is used on refresh.
- **Practitioner inputs**: `expires_after_hours` and `rotate_before_expiry_hours` (no `secret_` prefix inside the rotation block). Maps to Atlas `secret_expires_after_hours` on create and rotation POST.
- **Deletion policy**: No delete on v1→2; from v2→3 onward DELETE `old_secret` before POST so at most two active secrets remain.
- **Destroy**: Uses embedded service account credentials from state when present (same as read/update), otherwise provider-configured credentials.

## Example

```hcl
resource "mongodbatlas_organization3" "demo" {
  name         = "demo-org"
  org_owner_id = var.org_owner_id

  client_secret_rotation = {
    expires_after_hours        = 720
    rotate_before_expiry_hours = 360
  }
}
```

## Acceptance test env

- Provider credentials with permission to create organizations (same as `mongodbatlas_organization` tests).
- `MONGODB_ATLAS_ORG_OWNER_ID`: Atlas user ID for `org_owner_id`.
- Tests call `acc.SkipUnlessHasOrgOwner()` and `acc.SkipTestForCI()`; they create and delete real organizations (cost).

## Tests

Unit tests (no Atlas credentials):

```sh
cd code/provider
go test ./internal/service/organization3/ -run 'TestRenewalDue|TestModifyPlan|TestShouldDeleteOldSecret|TestEffectiveRotateBeforeExpiryHours' -v
```

Acceptance tests (real Atlas):

```sh
cd code/provider
TF_ACC=1 go test ./internal/service/organization3/ -run TestAccOrganization3_rotationLifecycle -v
```

Acceptance coverage (`TestAccOrganization3_rotationLifecycle`):

- Create with `expires_after_hours = 8` only; `secret_version = 1`.
- Empty plan on second apply (unchanged config).
- Force rotation with `secret_version = 2`; `old_secret` matches prior `current_secret`.
- Empty plan after removing `secret_version` from config.
- Widen policy with `rotate_before_expiry_hours = 8`; `secret_version = 3` and `old_secret` matches prior current.

```sh
golangci-lint run ./internal/service/organization3/...
```
