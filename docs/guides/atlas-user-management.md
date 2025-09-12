---
page_title: "Migration Guide: Atlas User Management"
---

# Migration Guide: Atlas User Management

## Overview

With MongoDB Atlas Terraform Provider `2.0.0`, several attributes and resources
were deprecated in favor of new, assignment-based resources. These changes
improve **clarity, separation of concerns, and alignment with Atlas APIs**. This
guide covers migrating to the new resources/attributes for Atlas user management
in context of **organization, teams, and projects**:

## Quick Finder: What changed

- **Org membership:** The `mongodbatlas_org_invitation` resource is deprecated.
  Use `mongodbatlas_cloud_user_org_assignment`. See
  [Org Invitation to Cloud User Org Assignment](#org-invitation-to-cloud-user-org-assignment).

- **Team membership:** The `usernames` attribute on `mongodbatlas_team` is
  deprecated. Use `mongodbatlas_cloud_user_team_assignment`. See
  [Team Usernames to Cloud User Team Assignment](#team-usernames-to-cloud-user-team-assignment).

- **Project team assignments:** The `teams` block inside `mongodbatlas_project`
  is deprecated. Use `mongodbatlas_team_project_assignment`. See
  [Project Teams to Team Project Assignment](#project-teams-to-team-project-assignment).

- **Project membership:** The `mongodbatlas_project_invitation` resource is
  deprecated. Use `mongodbatlas_cloud_user_project_assignment`. See
  [Project Invitation to Cloud User Project Assignment](#project-invitation-to-cloud-user-project-assignment).

- **Atlas User details:** The `mongodbatlas_atlas_user` and
  `mongodbatlas_atlas_users` data sources are deprecated. Use
  `mongodbatlas_cloud_user_org_assignment` for a single user in an org, and the
  `users` attributes on `mongodbatlas_organization`, `mongodbatlas_project`, or
  `mongodbatlas_team` for listings multiple users. See
  [Atlas User/Users Data Sources](#atlas-userusers-data-sources).

These updates ensure that **organization membership, team membership, and
project assignments** are modeled as explicit and independent resources — giving
you more flexible control over Atlas access management.

## Before You Begin

- Backup your
  [Terraform state](https://developer.hashicorp.com/terraform/cli/commands/state)
  file.
- Use MongoDB Atlas Terraform Provider **v2.0.0+** or later.
- Terraform version requirements:
  - **v1.5+** for
    **[import blocks](https://developer.hashicorp.com/terraform/language/import)**
    (earlier versions can use
    [`terraform import`](https://developer.hashicorp.com/terraform/cli/import))
  - **v1.1+** for
    **[moved blocks](https://developer.hashicorp.com/terraform/language/moved)**
    (useful for modules)
  - **v1.7+** for
    **[removed blocks](https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources)**
    (earlier versions can use
    [`terraform state rm`](https://developer.hashicorp.com/terraform/cli/commands/state/rm))

---

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Org Membership</span></summary>

## Org Invitation to Cloud User Org Assignment

**Objective**: Migrate from the deprecated `mongodbatlas_org_invitation`
resource and data source to the `mongodbatlas_cloud_user_org_assignment`
resource. If you previously assigned teams via `teams_ids`, also migrate those
to `mongodbatlas_cloud_user_team_assignment`.

### What’s changing?

- `mongodbatlas_org_invitation` only managed invitations and is deprecated. It
  didn’t manage the actual user membership or expose `user_id`.
- New `mongodbatlas_cloud_user_org_assignment` manages the user’s organization
  membership (pending or active) and exposes both `username` and `user_id`. It
  supports import using either `ORG_ID/USERNAME` or `ORG_ID/USER_ID`.
- If you previously used `teams_ids` on invitations, use
  `mongodbatlas_cloud_user_team_assignment` to manage team membership for each
  user.

---

### Use-case 1: Pending invites with `teams_ids`
When an invite is still pending and you have `teams_ids` defined in `mongodbatlas_org_invitation`, migrate both the org assignment and the team assignments.

#### Step 1: Replace `mongodbatlas_org_invitation` with `mongodbatlas_cloud_user_org_assignment`

- Original configuration:

```terraform
locals {
  org_id  = "<ORG_ID>"
  username = "user1@email.com"
  roles    = ["ORG_MEMBER"]
  team_ids = ["<TEAM_ID_1>", "<TEAM_ID_2>", "<TEAM_ID_3>"]
}

resource "mongodbatlas_org_invitation" "this" {
  username  = local.username
  org_id    = local.org_id
  roles     = local.roles
  teams_ids = local.team_ids
}
```

- New configuration:

```terraform
resource "mongodbatlas_cloud_user_org_assignment" "this" {
  org_id   = local.org_id
  username = local.username
  roles    = { org_roles = local.roles }
}
```

- Add a `moved` block (recommended) or an `import` block (if you cannot change module code):
```terraform
# Option A: moved block (recommended)
moved {
  from = mongodbatlas_org_invitation.this
  to   = mongodbatlas_cloud_user_org_assignment.this
}

# Option B: import block (use only if you can't use moved blocks; root module only)
import {
  to = mongodbatlas_cloud_user_org_assignment.this
  id = "${local.org_id}/${local.username}"
}

```


#### Step 2: Add `mongodbatlas_cloud_user_team_assignment`
Since `teams_ids` are no longer part of the org invitation, we need to manage them separately:
```terraform
resource "mongodbatlas_cloud_user_team_assignment" "team" {
  for_each = local.team_ids

  org_id  = local.org_id
  team_id = each.key
  user_id = mongodbatlas_cloud_user_org_assignment.this.user_id
}

# Import existing team assignments (root module only)
import {
  for_each = var.team_ids
  to       = mongodbatlas_cloud_user_team_assignment.team[each.key]
  id       = "${local.org_id}/${each.key}/${local.username}" # or use user_id
}

```

#### Step 3: Apply and clean up
- Run `terraform plan` (you should see import & moved operations), then `terraform apply`.
- Finally, remove any remaining `mongodbatlas_org_invitation` references from
config and state:
  ```terraform
  removed {
    from = mongodbatlas_org_invitation.this

    lifecycle {
      destroy = false
    }
  }
  ```
  - Alternatively, use the Terraform CLI command: `terraform state rm mongodbatlas_org_invitation.this`.

#### Module considerations

- **Module maintainers**
  - Add `mongodbatlas_cloud_user_org_assignment` inside the module and a `moved` block from `mongodbatlas_org_invitation`; remove the old resource and publish a new version.
  - If `teams_ids` were used, model them as `mongodbatlas_cloud_user_team_assignment` resources in the module that will be imported by module users.
  - Terraform doesn’t allow import blocks in the module ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)). Document the import ID formats for users:
        - Org assignment: `org_id/user_id`
        - Team assignment (if applicable): `org_id/team_id/user_id`

- **Module users**
  - Upgrade the module (`terraform init -upgrade`) and run `terraform plan` **but do not apply**.
  - Org assignment moves happen automatically via the module’s moved {}—no imports or state edits needed.
  - For team assignments, if applicable, add **root-level** `import {}` blocks (or run `terraform import`) for each existing:
        - Team assignment: `org_id/team_id/user_id`
  - Re-run `terraform plan` to confirm import & moved operations, then `terraform apply`.

  
---

### Use-case 2: Pending invites without `team_ids`

#### Step 1: Replace the org invite with `mongodbatlas_cloud_user_org_assignment` (same as Use-case 1 → Step 1)

```terraform
resource "mongodbatlas_cloud_user_org_assignment" "this" {
  org_id   = local.org_id
  username = local.username
  roles    = { org_roles = local.roles }
}

# Option A (recommended): moved block
moved {
  from = mongodbatlas_org_invitation.this
  to   = mongodbatlas_cloud_user_org_assignment.this
}

# Option B: import block (use only if you can't use moved blocks; root module only)
import {
  to = mongodbatlas_cloud_user_org_assignment.this
  id = "${local.org_id}/${local.username}"
}

```

#### Step 2: Apply and clean up
- Run `terraform plan` (you should see moved operation or imports if using import blocks), then `terraform apply`.
- Finally, remove any remaining `mongodbatlas_org_invitation` references from
config and state:
`terraform state rm mongodbatlas_org_invitation.this`.
  ```terraform
  removed {
    from = mongodbatlas_org_invitation.this

    lifecycle {
      destroy = false
    }
  }
  ```
  - Alternatively, use the Terraform CLI command: `terraform state rm mongodbatlas_org_invitation.this`.

#### Module considerations

- **Module maintainers**
  - Add `mongodbatlas_cloud_user_org_assignment` inside the module and a `moved` block from `mongodbatlas_org_invitation`; remove the old resource and publish a new version.
  
- **Module users**
  - Simply bump the module version and run `terraform init -upgrade`, then `terraform plan` and `terraform apply`. Terraform performs an in-place state move without users touching state.

---
### Use-case 3: Invitations already ACCEPTED (no `mongodbatlas_org_invitation` in config)

When an invite is accepted, Atlas deletes the underlying invitation. To manage
these users going forward, import them into
`mongodbatlas_cloud_user_org_assignment`.

#### Step 1: Fetch active org users (optional helper)

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

#### Step 2: Define and import `mongodbatlas_cloud_user_org_assignment`

Use the `local.active_users` map defined in Step 1 so you don’t have to manually
curate a list:

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

#### Step 3 (Optional): Add team assignments if needed
- If you also need teams, reuse Use-case 1 → Step 2 with a `for_each` over your team IDs per user.

#### Step 4: Apply and clean up
- Run `terraform plan` (you should see import operations planned), then `terraform apply`.
- Finally, remove any remaining `mongodbatlas_org_invitation` references from
config and state:
`terraform state rm mongodbatlas_org_invitation.this`.
  ```terraform
  removed {
    from = mongodbatlas_org_invitation.this

    lifecycle {
      destroy = false
    }
  }
  ```
  - Alternatively, use the Terraform CLI command: `terraform state rm mongodbatlas_org_invitation.this`.

#### Module considerations

- **Module maintainers**
  - Add `mongodbatlas_cloud_user_org_assignment` in the module. Since invites are already **accepted**, these existing org users need to be imported to be managed with Terraform going forward.
  - If teams are in scope, define `mongodbatlas_cloud_user_team_assignment` in the module as well.
  - Terraform doesn’t allow import blocks in the module ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)). Document the import ID formats for users:
      - Org assignment: `org_id/user_id`
      - Team assignment (if applicable): `org_id/team_id/user_id`
  - Publish a new module version.

- **Module users**
  - Upgrade the module (`terraform init -upgrade`) and run `terraform plan` **but do not apply**.
  - Add **root-level** `import {}` blocks (or run `terraform import`) for each existing:
      - Org assignment: `org_id/user_id`
      - Team assignment (if applicable): `org_id/team_id/user_id`
  - Re-run `terraform plan` to confirm import operations, then `terraform apply`.


---

### Data source migration

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
```

Replace with the new data source:

```terraform
data "mongodbatlas_cloud_user_org_assignment" "user_1" {
  org_id   = local.org_id
  username = local.username
}
```

Then:

1. Run `terraform apply` to ensure the new data source reads correctly.
2. Replace all usages of `data.mongodbatlas_org_invitation.test` with
   `data.mongodbatlas_cloud_user_org_assignment.user_1`.
3. Run `terraform plan`, then `terraform apply`.



### Examples

For complete, working configurations that mirror the use-cases above, see the
examples in the provider repository:
[migrate_org_invitation_to_cloud_user_org_assignment](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_org_invitation_to_cloud_user_org_assignment).
These include root-level setups for multiple approaches (e.g., moved blocks and
imports) across different versions.

### Notes and tips

- Import formats:
  - Org assignment: `ORG_ID/USERNAME` or `ORG_ID/USER_ID`.
  - Team assignment: `ORG_ID/TEAM_ID/USERNAME` or `ORG_ID/TEAM_ID/USER_ID`.
- If you use modules, keep in mind import blocks must be placed at the root
  module.
- After successful migration, ensure no references to
  `mongodbatlas_org_invitation` remain.
- [Cloud User Org Assignment Resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_user_org_assignment)

</details>

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Team Membership</span></summary>

## Team Usernames to Cloud User Team Assignment

**Objective**: Migrate from the deprecated `usernames` attribute on the
`mongodbatlas_team` resource to the new
`mongodbatlas_cloud_user_team_assignment` resource.

### Why should I migrate?

- **Future Compatibility:** The `usernames` attribute on `mongodbatlas_team` is
  deprecated and may be removed in future provider versions. Migrating ensures
  your Terraform configuration remains functional.
- **Flexibility:** Manage teams and user assignments independently, without
  coupling membership changes to team creation or updates.
- **Clarity:** Clear separation between the `mongodbatlas_team` resource (team
  definition) and `mongodbatlas_cloud_user_team_assignment` (membership
  management).

### What’s changing?

- `mongodbatlas_team` included a `usernames` argument that allowed assigning
  users to a team directly inside the resource. This argument is now deprecated.
- New attribute `users` in `mongodbatlas_team` data source can be used to
  retrieve information about all the users assigned to that team.
- `mongodbatlas_cloud_user_team_assignment` manages the user’s team membership
  (pending or active) and exposes both `username` and `user_id`. It supports
  import using either `ORG_ID/TEAM_ID/USERNAME` or `ORG_ID/TEAM_ID/USER_ID`.

---

### Migrate from `mongodbatlas_team.usernames` to `mongodbatlas_cloud_user_team_assignment`

#### Original configuration

```terraform
locals {
  usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
}

resource "mongodbatlas_team" "this" {  
  org_id    = var.org_id  
  name      = var.team_name
  usernames = local.usernames
}
```

#### Step 1: Use `mongodbatlas_team` data source to retrieve user IDs

We first need to retrieve each user's `user_id` via the new `users` attribute in
`mongodbatlas_team` data source.

```terraform 
locals {
    usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
    
    user_ids = toset([for u in data.mongodbatlas_team.this.users : u.id])   # Use data source to get team members (with user_id) 
}

data "mongodbatlas_team" "this" {  
    org_id  = var.org_id  
    team_id = mongodbatlas_team.this.team_id  
}

resource "mongodbatlas_team" "this" {  
    org_id = var.org_id  
    name   = var.team_name
    usernames = local.usernames
} 
```

#### Step 2: Add `mongodbatlas_cloud_user_team_assignment` and use import blocks

```terraform
# New resource for each (user, team) assignment  
resource "mongodbatlas_cloud_user_team_assignment" "this" {
  for_each = local.user_ids

  org_id  = var.org_id
  team_id = mongodbatlas_team.this.team_id
  user_id = each.value         # Use user_id instead of username  
}
  
# Import existing team-user relationships into the new resources (root module only)
import {  
    for_each = local.user_ids

    to = mongodbatlas_cloud_user_team_assignment.this[each.key] 
    id = "${var.org_id}/${mongodbatlas_team.this.team_id}/${each.value}" 
}
```

#### Step 3: Remove deprecated `usernames` from `mongodbatlas_team`

Once the new resources are in place:

```terraform
resource "mongodbatlas_team" "this" {  
  org_id = var.org_id  
  name   = "this"  
  # usernames = local.usernames  # Remove this line
}
```

#### Step 4: Run migration

Run `terraform plan` (you should see **import** operations), then
`terraform apply`.

#### Step 5: Update any references to `mongodbatlas_team.usernames`

Before:

```terraform
output "team_usernames" {  
  value = mongodbatlas_team.this.usernames  
}
```

After:

```terraform
output "team_usernames" {  
  value = [for u in data.mongodbatlas_team.this.users : u.username]  
}
```

Run `terraform plan`. There should be **no changes**.

---

#### Module considerations
The legacy `mongodbatlas_team.usernames` list maps to individual
`mongodbatlas_cloud_user_team_assignment` resources, so a `moved` block
cannot be used. Existing team memberships must be imported.

- **Module maintainers**
  - Define `mongodbatlas_cloud_user_team_assignment` inside the module.
  - Example **old** module implementation:
    ```terraform
      variable "org_id"    { type = string }
    variable "team_name" { type = string }
    variable "usernames" { type = list(string) }

    resource "mongodbatlas_team" "this" {
      org_id    = var.org_id
      name      = var.team_name
      usernames = var.usernames  # deprecated
    }
    ```
  - Example **new** module implementation:
    ```terraform
    variable "org_id"    { type = string }
    variable "team_name" { type = string }
    variable "user_ids"  { type = set(string) }

    resource "mongodbatlas_team" "this" {
      org_id = var.org_id
      name   = var.team_name
      # removed deprecated usernames
    }

    resource "mongodbatlas_cloud_user_team_assignment" "this" {
      for_each = var.user_ids
      
      org_id   = var.org_id
      team_id  = mongodbatlas_team.this.team_id
      user_id  = each.value
    }
    ```
  - Terraform doesn’t allow import blocks in the module ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)). Document the import ID formats for users:
      - Team assignment: `org_id/team_id/user_id` (or `org_id/team_id/username`)
  - Publish a new module version.

- **Module users**
  - Upgrade to the new module version (`terraform init -upgrade`) and run terraform plan but **do not apply**.
  - Example **old** module usage (using deprecated resources):
    ```hcl
    module "user_team_assignment" {  
      source     = "./old_module"  
      org_id     = var.org_id  
      team_name  = var.team_name  
      usernames  = var.usernames 
    }
    ```
  - Example **new** module usage:
    ```hcl
    data "mongodbatlas_team" "this" {  
      org_id = var.org_id  
      name   = var.team_name
    }

    locals {  
      user_ids = toset([
        for user in data.mongodbatlas_team.this.users : user.id
      ]) 
    }

    module "user_team_assignment" {
      source     = "./new_module"
      org_id     = var.org_id
      team_name  = var.team_name
      user_ids = local.user_ids   # replaced deprecated usernames
    }
    ```
  - Add an `import block` (or `terraform import`) to import the resources and iterate through the list of users:
    ```terraform
    import { 
        for_each = local.team_assignments
        to       = module.user_team_assignment.mongodbatlas_cloud_user_team_assignment.this[each.key]
        id       = "${var.org_id}/${data.mongodbatlas_team.this.team_id}/${each.value}"
    }
    ```
  - Run `terraform plan` to review the changes.
      - Ensure that Terraform imports the user-team assignments and does not plan to create these.
      - Ensure that Terraform does not plan to modify the `mongodbatlas_team` resource.
  - Run `terraform apply` to apply the migration.

For complete working examples, see:

- [Old module definition](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_maintainer/v1)
  and
  [old module usage](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_user/v1).
- [New module definition](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_maintainer/v2)
  and
  [new module usage](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_user/v2).
- [mongodbatlas_cloud_user_team_assignment](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_cloud_user_team_assignment/main.tf).

---

### Data source migration

If you previously used the `usernames` attribute in the `data.mongodbatlas_team`
data source:

**Original:**

```terraform
output "team_usernames" {  
  description = "Usernames in the MongoDB Atlas team"  
  value       = data.mongodbatlas_team.this.usernames  
}
```

**Replace with:**

```terraform
output "team_usernames" { 
  description = "Usernames in the MongoDB Atlas team"  
  value = [for u in data.mongodbatlas_team.this.users : u.username]  
}
```

Run `terraform plan`. There should be **no changes**.

---

### Notes and tips

- **Import format** for `mongodbatlas_cloud_user_team_assignment`:

```
ORG_ID/TEAM_ID/USERNAME
ORG_ID/TEAM_ID/USER_ID
```

- After successful migration, ensure **no references to**
  `mongodbatlas_team.usernames` remain.

- [Cloud User Team Assignment Resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_user_team_assignment)

</details>

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Project Team Assignment</span></summary>

## Project Teams to Team Project Assignment

**Objective:** Migrate from the deprecated `teams` attribute on the
`mongodbatlas_project` resource to the new
`mongodbatlas_team_project_assignment` resource.

### Why should I migrate?

- **Future compatibility:** The `teams` attribute inside `mongodbatlas_project`
  is deprecated and will be removed in a future provider release.
- **Separation of concerns:** Manage projects and team-to-project role
  assignments independently.
- **Clearer diffs:** Role or team modifications won't require re‑applying the
  entire project resource.

### What's changing?

Historically, `mongodbatlas_project` accepted an inline `teams` block to assign
one or more teams to a project with specific roles. Now, each project-team role
mapping must be managed with `mongodbatlas_team_project_assignment`.

---

### Migrate from `mongodbatlas_project.teams` to `mongodbatlas_team_project_assignment`

#### Original configuration

```hcl
locals {  
  team_map = { # team_id => set(role_names)
    <TEAM_ID_1>  = ["GROUP_OWNER"]
    <TEAM_ID_2>  = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
  }
}

resource "mongodbatlas_project" "this" {
  name             = var.project_name
  org_id           = var.org_id
  project_owner_id = var.project_owner_id

  dynamic "teams" {
    for_each = local.team_map
    content {  
      team_id    = teams.key  
      role_names = teams.value  
    }  
  }  
}
```

#### Step 1: Ignore `teams` and remove from configuration

-> **Note:** The `teams` attribute is a `SetNestedBlock` and cannot be marked
`Optional`/`Computed` for a smooth migration. For now, `ignore_changes` is
required during Step 1. Support for removing `teams` entirely will come in a
future Atlas Provider release.

- Replace the `mongodbatlas_project.teams` block with:

```hcl
resource "mongodbatlas_project" "this" {  
  name             = var.project_name
  org_id           = var.org_id
  project_owner_id = var.project_owner_id
  
  lifecycle {  
    # Ignore `teams` field as it's deprecated.
    # It can now be managed with the new `mongodbatlas_team_project_assignment` resources
    ignore_changes = ["teams"]  
  }  
}
```

- Run `terraform plan`, then `terraform apply`.


This removes the `teams` block from the config but keeps the assignments in
Atlas unchanged until we explicitly manage them in new resources.

#### Step 2: Add the new `mongodbatlas_team_project_assignment` resources

```hcl
resource "mongodbatlas_project" "this" {  
  name             = var.project_name
  org_id           = var.org_id
  project_owner_id = var.project_owner_id
  
  lifecycle {  
    ignore_changes = ["teams"]  
  }  
}

resource "mongodbatlas_team_project_assignment" "this" {  
  for_each = local.team_map  
  
  project_id = mongodbatlas_project.this.id  
  team_id    = each.key  
  role_names = each.value  
}  
 
import {  
  for_each = local.team_map

  to       = mongodbatlas_team_project_assignment.this[each.key]
  id       = "${mongodbatlas_project.this.id}/${each.key}"
}
```

- Run `terraform plan` (you should see **import** operations), then
`terraform apply`.

#### Step 3: Verify and clean up

- After successful import and apply, `terraform plan` should show **no
  changes**.
- Keep the `ignore_changes = ["teams"]` lifecycle rule until the provider
  releases a version without the `teams` argument in `mongodbatlas_project`.

#### Module considerations
Inline `mongodbatlas_project.teams` now maps to separate
`mongodbatlas_team_project_assignment` resources, so no `moved` block is possible.
Existing assignments must be imported at the root module. 

Keep
`ignore_changes = ["teams"]` on the project until the provider removes that field.

- **Module maintainers**
  - Replace the inline `mongodbatlas_project.teams` block with explicit `mongodbatlas_team_project_assignment` resources in the module and add a lifecycle rule to ignore `teams` in `mongodbatlas_project` as mentioned in Step #1 and #2 above.
  - Expose the `project_id` as a module output so users can form import IDs.
  - Terraform doesn’t allow import blocks in the module ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)). Document the import ID formats for users:
      - `project_id/team_id`
  - Publish a new module version.

- **Module users**
  - Upgrade the module (`terraform init -upgrade`) and run `terraform plan` **but do not apply**.
  - Similar to original configuration above, you can have a mapping of team IDs → role names for the project. Alternatively, this can be done using the `data.mongodbatlas_project.teams` attribute to get the existing team IDs → role names mapping. 
  -  Similar to Step #2, add **root-level** `import {}` blocks (or run `terraform import`) for existing project–team assignments:
    - Target the module resource address for each team assignment, for example:
      ```terraform 
      # Import each existing PROJECT_ID/TEAM_ID into the module resource address
      import {
        for_each = var.team_map   # team_id => set(role_names)
        
        to       = module.project.mongodbatlas_team_project_assignment.this[each.key]  # each.key = TEAM_ID
        id       = "${module.project.project_id}/${each.key}"                          # PROJECT_ID/TEAM_ID
      }
      ```
  - Re-run `terraform plan` to confirm import operations, then `terraform apply`.


---

### Examples

For complete, working configurations that demonstrate the migration process, see
the examples in the provider repository:
[migrate_team_project_assignment](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_team_project_assignment).

The examples include:

- **v1**: Original configuration using deprecated `teams` attribute in
  `mongodbatlas_project` resource.
- **v2**: Final configuration using `mongodbatlas_team_project_assignment`
  resource for team-to-project assignments.

---

### Notes and tips

- **Import format** for `mongodbatlas_team_project_assignment`:

```
PROJECT_ID/TEAM_ID
```

- [Atlas Team Project Assignment Resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/team_project_assignment)

---

### FAQ

**Q: Do I need to delete the old `teams` from state?** A: No — using
`ignore_changes` ensures they remain in Atlas until the provider removes the
field. Then you can drop the lifecycle rule.


</details>

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Project Membership</span></summary>

## Project Invitation to Cloud User Project Assignment

**Objective**: Migrate from the deprecated `mongodbatlas_project_invitation`
resource and data source to the `mongodbatlas_cloud_user_project_assignment`
resource.

### What’s changing?

- `mongodbatlas_project_invitation` only managed invitations and is deprecated.
  When the user accepted the invitation and became a project member, the
  underlying invitation entity went away and you needed to remove it from your
  configuration as well. See the resource
  [documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project_invitation)
  for more details.
- `mongodbatlas_cloud_user_project_assignment` manages the user’s project
  membership (both invited and active members).
- Pending project invitations are not discoverable with the new APIs. The only
  migration path for existing PENDING invites is to re-create them using
  `mongodbatlas_cloud_user_project_assignment` with the same `username` and
  `roles`.
- For details on the new resource, see the
  `mongodbatlas_cloud_user_project_assignment` resource documentation:
  https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_user_project_assignment

---

### Migrating PENDING invitations

Original configuration:

```terraform
locals {
  username = "user1@email.com"
  roles    = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

resource "mongodbatlas_project_invitation" "this" {
  project_id = var.project_id
  username   = local.username
  roles      = local.roles
}
```

#### Step 1: Add the new resource to re-create the pending invite via the new API:

```terraform
resource "mongodbatlas_cloud_user_project_assignment" "this" {
  project_id = var.project_id
  username   = local.username
  roles      = local.roles
}
```

Use the same `roles` as the original invitation to avoid drift.

#### Step 2: Delete the deprecated `mongodbatlas_project_invitation` resource block


#### Step 3: Apply the changes

Run `terraform apply` to create the assignment with the new resource & delete the current `mongodbatlas_project_invitation` resource.

---

#### Module considerations

- **Module maintainers**
  - Replace `mongodbatlas_project_invitation` with `mongodbatlas_cloud_user_project_assignment` inside the module.
  - Keep inputs consistent (`project_id`, `username`, `roles`) so the new resource re-creates the pending invite with the same roles.
  - Remove the deprecated `mongodbatlas_project_invitation` resource block from the module.
  - Publish a new module version.

- **Module users**
  - Upgrade to the new module version and run `terraform plan`.
  - Expect to see planned creation `mongodbatlas_cloud_user_project_assignment` and deletion of `mongodbatlas_project_invitation`.
  - Run `terraform apply`.

---

### Examples

For complete, working configurations that demonstrate the migration process, see
the examples in the provider repository:
[migrate_project_invitation_to_cloud_user_project_assignment](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_project_invitation_to_cloud_user_project_assignment).

The examples include:

- **v1**: Original configuration using deprecated
  `mongodbatlas_project_invitation`
- **v2**: Migration phase with re-creation using new resource and clean state
  removal
- **v3**: Final clean configuration using only
  `mongodbatlas_cloud_user_project_assignment`

These examples provide practical validation of the migration steps and
demonstrate the re-creation approach for pending invitations.

---

### Notes and tips

- After successful migration, ensure no references to
  `mongodbatlas_project_invitation` remain in configuration or state.
- Pending invitations are not discoverable by the new APIs and resources; there
  is no data source replacement for reading pending invites. Re-create them
  using the new resource as shown above.
- For additional details on how accepted invitations are handled, see the
  `mongodbatlas_project_invitation` resource
  [documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project_invitation).
- [Cloud User Project Assignment Resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_user_project_assignment)

</details>

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Atlas User details</span></summary>

## Atlas User/Users Data Sources

**Objective**: Migrate from the deprecated `mongodbatlas_atlas_user` and
`mongodbatlas_atlas_users` data sources to their respective replacements.

### What’s changing?

- `mongodbatlas_atlas_user` returned a user profile by `user_id` or `username`
  and is deprecated. Replace it with `mongodbatlas_cloud_user_org_assignment`
  which reads a user's assignment in a specific organization using either
  `username` or `user_id` together with `org_id`. For details, see the
  `mongodbatlas_cloud_user_org_assignment` data source
  [documentation](../data-sources/cloud_user_org_assignment).

- `mongodbatlas_atlas_users` returned lists of users by `org_id`, `project_id`,
  or `team_id` and is deprecated. Replace it with the `users` attribute
  available on `mongodbatlas_organization`, `mongodbatlas_project`, or
  `mongodbatlas_team` data sources, respectively. The resulting `users` list now
  includes both active and pending users.
- Attribute structure differences: The new organization users API does not
  return `email_address` as a separate field and replaces the consolidated
  `roles` with structured `org_roles` and `project_role_assignments`.

---

### Migrate reads to `mongodbatlas_cloud_user_org_assignment`

Original configuration:

```terraform
data "mongodbatlas_atlas_user" "test" {
  user_id = "<USER_ID>"
}

# OR

data "mongodbatlas_atlas_user" "test" {
  username = "<USERNAME>"
}
```

#### Step 1: Add the new data source alongside the existing one

Use either `username` or `user_id` with the target `org_id`:

```terraform
# Keep existing data source temporarily
data "mongodbatlas_atlas_user" "test" {
  user_id = "<USER_ID>"
}

# Add new data source
data "mongodbatlas_cloud_user_org_assignment" "user_1" {
  user_id = "<USER_ID>"
  org_id  = "<ORGANIZATION_ID>"
}
```

#### Step 2: Verify the new data source works

Run `terraform plan` to ensure the new data source will read correctly without
errors.

#### Step 3: Replace references incrementally

Replace references from `data.mongodbatlas_atlas_user.test` to
`data.mongodbatlas_cloud_user_org_assignment.user_1`.

**Important**: Update attribute references as the structure has changed:

Key attribute changes:

| Old Attribute                  | New Attribute                                     |
| ------------------------------ | ------------------------------------------------- |
| `email_address`                | `username`                                        |
| `roles` (filtered by org_id)   | `roles.org_roles`                                 |
| `roles` (filtered by group_id) | `roles.project_role_assignments[*].project_roles` |

**Examples**:

- Email: `data.mongodbatlas_atlas_user.test.email_address` →
  `data.mongodbatlas_cloud_user_org_assignment.user_1.username`
- Org roles: Use
  `data.mongodbatlas_cloud_user_org_assignment.user_1.roles.org_roles` directly
- Project roles: Access via `roles.project_role_assignments` list, filtering by
  `project_id` as needed

#### Step 4: Remove the old data source

Once all references are updated and working, remove the old data source from
your configuration:

```terraform
# Remove this block
# data "mongodbatlas_atlas_user" "test" {
#   user_id = "<USER_ID>"
# }
```

#### Step 5: Apply and verify

Run `terraform plan` to ensure no unexpected changes, then `terraform apply`.

---

### Migrate list reads from `mongodbatlas_atlas_users`

Original configuration:

```terraform
data "mongodbatlas_atlas_users" "test" {
  org_id = "<ORG_ID>"
}

# OR

data "mongodbatlas_atlas_users" "test" {
  project_id = "<PROJECT_ID>"
}

# OR

data "mongodbatlas_atlas_users" "test" {
  team_id = "<TEAM_ID>"
  org_id  = "<ORG_ID>"
}
```

#### Step 1: Add new data sources alongside existing ones

Add the appropriate replacement data source(s) while keeping the old one
temporarily:

Organization users:

```terraform
# Keep existing temporarily
data "mongodbatlas_atlas_users" "test" {
  org_id = "<ORG_ID>"
}

# Add new data source
data "mongodbatlas_organization" "org" {
  org_id = "<ORG_ID>"
}

locals {
  org_users = data.mongodbatlas_organization.org.users
}
```

Project users:

```terraform
# Keep existing temporarily  
data "mongodbatlas_atlas_users" "test" {
  project_id = "<PROJECT_ID>"
}

# Add new data source
data "mongodbatlas_project" "proj" {
  project_id = "<PROJECT_ID>"
}

locals {
  project_users = data.mongodbatlas_project.proj.users
}
```

Team users:

```terraform
# Keep existing temporarily
data "mongodbatlas_atlas_users" "test" {
  team_id = "<TEAM_ID>"
  org_id  = "<ORG_ID>"
}

# Add new data source
data "mongodbatlas_team" "team" {
  team_id = "<TEAM_ID>"
  org_id  = "<ORG_ID>"
}

locals {
  team_users = data.mongodbatlas_team.team.users
}
```

#### Step 2: Verify new data sources work

Run `terraform plan` to ensure the new data sources read correctly and return
expected user data.

#### Step 3: Replace references incrementally

Replace `data.mongodbatlas_atlas_users.test.results` with the appropriate
`...users` collection above.

**Important**: Update attribute references as the structure has changed:

| Old Attribute                 | New Attribute                                  |
| ----------------------------- | ---------------------------------------------- |
| `results[*].email_address`    | `users[*].username`                            |
| `results[*].roles` (filtered) | `users[*].roles.org_roles` or `users[*].roles` |

**Examples**:

- Email list: `data.mongodbatlas_atlas_users.test.results[*].email_address` →
  `data.mongodbatlas_organization.org.users[*].username`
- User list: `data.mongodbatlas_atlas_users.test.results` →
  `data.mongodbatlas_organization.org.users` (or `.project.proj.users`,
  `.team.team.users`)
- Org roles: Use `users[*].roles.org_roles` from organization data source
- Project roles: Use `users[*].roles` from project data source, or
  `users[*].roles.project_role_assignments` from organization data source

#### Step 4: Remove the old data source

Once all references are updated and working, remove the old data source from
your configuration:

```terraform
# Remove this block
# data "mongodbatlas_atlas_users" "test" {
#   org_id = "<ORG_ID>"
# }
```

#### Step 5: Apply and verify

Run `terraform plan` to ensure no unexpected changes, then `terraform apply`.

---

#### Module considerations
Since data sources don’t live in state, in this case migration is about replacing data sources and updating attribute references (and, if needed, module inputs/outputs).

- **Module maintainers**
  - Replace deprecated data sources with the new resources as mentioned in above steps.
  - Update attribute references as mentioned above.
  - Publish a new module version.

- **Module users**
  - Upgrade to the new module version and run `terraform plan`.
  - Update your references to the module’s outputs/variables to match the new attribute structure (use the mapping above).
  - Re-run `terraform plan` to confirm reads succeed and the output shape is as expected, then proceed as usual.



---

### Examples

For complete, working configurations that demonstrate the migration process, see
the examples in the provider repository:
[migrate_atlas_user_and_atlas_users](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_atlas_user_and_atlas_users).

The examples include:

- **v1**: Original configuration using deprecated data sources
- **v2**: Migration phase with side-by-side comparison and validation
- **v3**: Final clean configuration using only new data sources

These examples provide practical validation of the migration steps and
demonstrate the attribute mappings in working Terraform code.

---

### Notes

- The new data source requires the `org_id` context to read the user's
  organization assignment.
- After migration, ensure no remaining references to `mongodbatlas_atlas_user`
  exist in your configuration.

</details>
