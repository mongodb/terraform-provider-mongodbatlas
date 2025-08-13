---
page_title: "Migration Guide: Org Invitation to Cloud User Org Assignment"
---

# Migration Guide: Org Invitation to Cloud User Org Assignment

**Objective**: Migrate from the deprecated `mongodbatlas_org_invitation` resource and data source to the `mongodbatlas_cloud_user_org_assignment` resource. If you previously assigned teams via `teams_ids`, also migrate those to `mongodbatlas_cloud_user_team_assignment`.

## Before you begin

- Back up your Terraform state file.
- Ensure you are using the MongoDB Atlas Terraform Provider 2.0.0 version that includes `mongodbatlas_cloud_user_org_assignment` and `mongodbatlas_cloud_user_team_assignment` resources.

## What’s changing?

- `mongodbatlas_org_invitation` only managed invitations and is deprecated. It didn’t manage the actual user membership or expose `user_id`.
- `mongodbatlas_cloud_user_org_assignment` manages the user’s organization membership (pending or active) and exposes both `username` and `user_id`. It supports import using either `ORG_ID/USERNAME` or `ORG_ID/USER_ID`.
- If you previously used `teams_ids` on invitations, use `mongodbatlas_cloud_user_team_assignment` to manage team membership per user and team.

---

## Use-case 1: Existing org invite is still PENDING (resource exists in config)

Original configuration (note: `user_id` does not exist on `mongodbatlas_org_invitation`):

```terraform
locals {
  org_id  = "<ORG_ID>"
  username = "user1@email.com"
  roles    = ["ORG_MEMBER"]
}

resource "mongodbatlas_org_invitation" "this" {
  username  = local.username
  org_id    = local.org_id
  roles     = local.roles
  # teams_ids = local.team_ids  # if applicable
}
```

### Step 1: Add `mongodbatlas_cloud_user_org_assignment`

**Option A) [Recommended]** Moved block (module-friendly)

Why this is module-friendly and recommended:
- For module maintainers: Add the new `mongodbatlas_cloud_user_org_assignment` resource inside the module, include a `moved {}` block from `mongodbatlas_org_invitation` to the new resource, and publish a new module version.
- For module users: Simply bump the module version and run `terraform init -upgrade`, then `terraform plan` / `terraform apply`. Terraform performs an in-place state move without users writing import blocks or touching state.
- Works at any scale (any number of module instances) and keeps the migration self-contained within the module. No per-environment import steps are required.

```terraform
resource "mongodbatlas_cloud_user_org_assignment" "this" {
  org_id   = local.org_id
  username = local.username
  roles    = { org_roles = local.roles }
}

moved {
  from = mongodbatlas_org_invitation.this
  to   = mongodbatlas_cloud_user_org_assignment.this
}
```

**Option B)** Import by username (not module-friendly)

Why this is NOT module-friendly:
- Terraform import blocks cannot live inside modules; they must be defined in the root module. See `https://github.com/hashicorp/terraform/issues/33474`.
- Module maintainers cannot ship import steps. Each module user must add root-level import blocks for every instance to import, which is error-prone and repetitive.
- This creates extra coordination for every environment and workspace. Prefer Option A whenever you can modify the module source.

```terraform
resource "mongodbatlas_cloud_user_org_assignment" "this" {
  org_id   = local.org_id
  username = local.username
  roles    = { org_roles = local.roles }
}

import {
  to = mongodbatlas_cloud_user_org_assignment.this
  id = "${local.org_id}/${local.username}"
}
```

### Step 2: Remove `mongodbatlas_org_invitation` from config and state

- With a moved block, `terraform plan` should show the move and no other changes. Then `terraform apply`.
- If you used import, remove the old `mongodbatlas_org_invitation` block and delete it from state if still present: `terraform state rm mongodbatlas_org_invitation.this`.

---

## Use-case 2: Invitations already ACCEPTED (no `mongodbatlas_org_invitation` in config)

When an invite is accepted, Atlas deletes the underlying invitation. To manage these users going forward, import them into `mongodbatlas_cloud_user_org_assignment`.

### Step 1: Fetch active org users (optional helper)

```terraform
data "mongodbatlas_organization" "org" {
  org_id = var.org_id
}

locals {
  active_users = {
    for u in data.mongodbatlas_organization.org.users :
    u.id => u if u.org_membership_status == "ACTIVE"
  }
}
```

### Step 2: Define and import `mongodbatlas_cloud_user_org_assignment`

- Terraform import blocks cannot live inside modules; they must be defined in the root module. See `https://github.com/hashicorp/terraform/issues/33474`.

Use the `local.active_users` map defined in Step 1 so you don’t have to manually curate a list:

```terraform
resource "mongodbatlas_cloud_user_org_assignment" "user" {
  for_each = local.active_users  # key = user_id, value = user object from data source

  org_id   = var.org_id
  username = each.value.username

  # Keep roles aligned with current assignments to avoid drift after import
  roles = {
    org_roles = each.value.roles[0].org_roles
  }
}

# Import existing users (root module only)
import {
  for_each = local.active_users
  to       = mongodbatlas_cloud_user_org_assignment.user[each.key]
  id       = "${var.org_id}/${each.key}"  # org_id/user_id
}
```

Run `terraform plan` (you should see import operations), then `terraform apply`.

---

## Use-case 3: You also set `teams_ids` on the original invitation

Migrate team assignments to `mongodbatlas_cloud_user_team_assignment` in addition to Use-case 1 or 2 above.

```terraform
variable "team_ids" { type = set(string) }

resource "mongodbatlas_cloud_user_team_assignment" "team" {
  for_each = var.team_ids

  org_id  = local.org_id
  team_id = each.key
  user_id = mongodbatlas_cloud_user_org_assignment.this.id
}

# Import existing team assignments (root module only)
import {
  for_each = var.team_ids
  to       = mongodbatlas_cloud_user_team_assignment.team[each.key]
  id       = "${local.org_id}/${each.key}/${local.username}" # OR use user_id in place of username
}
```

Run `terraform plan` (you should see import operations), then `terraform apply`.

Finally, remove any remaining `mongodbatlas_org_invitation` references from config and state.

---

## Data source migration

Original configuration:

```terraform
locals {
  org_id  = "<ORG_ID>"
  username = "user1@email.com"
}

data "mongodbatlas_org_invitation" "test" {
  org_id        = local.org_id
  username      = local.username
  invitation_id = mongodbatlas_org_invitation.test.invitation_id
}
``;

Replace with the new data source:

```terraform
data "mongodbatlas_cloud_user_org_assignment" "user_1" {
  org_id   = local.org_id
  username = local.username
}
```

Then:

1. Run `terraform apply` to ensure the new data source reads correctly.
2. Replace all usages of `data.mongodbatlas_org_invitation.test` with `data.mongodbatlas_cloud_user_org_assignment.user_1`.
3. Run `terraform plan` followed by `terraform apply`.

---

## Notes and tips

- Import formats:
  - Org assignment: `ORG_ID/USERNAME` or `ORG_ID/USER_ID`.
  - Team assignment: `ORG_ID/TEAM_ID/USERNAME` or `ORG_ID/TEAM_ID/USER_ID`.
- If you use modules, keep in mind import blocks must be placed at the root module.
- After successful migration, ensure no references to `mongodbatlas_org_invitation` remain.

