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

### Use-case 1: Existing org invite is still PENDING (resource exists in config)

Original configuration (note: `user_id` does not exist on
`mongodbatlas_org_invitation`):

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
  # teams_ids = local.team_ids  # if applicable, also see Use-case #3 below
}
```

### Option A) [Recommended] Moved block

#### Step 1: Add `mongodbatlas_cloud_user_org_assignment` and `moved` block

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

#### Step 2: Remove `mongodbatlas_org_invitation` from config and state

- With a moved block, `terraform plan` should show the move and no other
  changes. Then `terraform apply`.

#### Module considerations

- For module maintainers: Add the new `mongodbatlas_cloud_user_org_assignment`
  resource inside the module with a `moved {}` block from
  `mongodbatlas_org_invitation` to the new resource, remove current
  `mongodbatlas_org_invitation` resource (Step 2) and publish a new module
  version.
- For module users: Simply bump the module version and run
  `terraform init -upgrade`, then `terraform plan` / `terraform apply`.
  Terraform performs an in-place state move without users writing import blocks
  or touching state.
- Works at any scale (any number of module instances) and keeps the migration
  self-contained within the module. No per-environment import steps are
  required.

### Option B) Import by username

#### Step 1: Add `mongodbatlas_cloud_user_org_assignment` and `import` block

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

#### Step 2: Remove `mongodbatlas_org_invitation` from config and state

- With import, remove the old `mongodbatlas_org_invitation` block and delete it
  from state if still present:
  `terraform state rm mongodbatlas_org_invitation.this`.

#### Module considerations

- Terraform import blocks cannot live inside modules; they must be defined in
  the root module. See
  ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)).
- Module maintainers cannot ship import steps. Each module user must add
  root-level import blocks for every instance to import, which is error-prone
  and repetitive.
- This creates extra coordination for every environment and workspace. Prefer
  Option A whenever you can modify the module source.

---

### Use-case 2: Invitations already ACCEPTED (no `mongodbatlas_org_invitation` in config)

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

#### Module considerations

- Terraform import blocks cannot live inside modules; they must be defined in
  the root module. See
  ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)).

Run `terraform plan` (you should see import operations), then `terraform apply`.

---

### Use-case 3: You also set `teams_ids` on the original invitation

Original configuration where `mongodbatlas_org_invitation` defines `teams_ids`:

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
  teams_ids = local.team_ids
}
```

Migrate team assignments to `mongodbatlas_cloud_user_team_assignment` in
addition to Use-case 1 or 2 above.

```terraform
variable "team_ids" { type = set(string) }

resource "mongodbatlas_cloud_user_team_assignment" "team" {
  for_each = var.team_ids

  org_id  = local.org_id
  team_id = each.key
  user_id = mongodbatlas_cloud_user_org_assignment.this.user_id
}

# Import existing team assignments (root module only)
import {
  for_each = var.team_ids
  to       = mongodbatlas_cloud_user_team_assignment.team[each.key]
  id       = "${local.org_id}/${each.key}/${local.username}" # OR use user_id in place of username
}
```

Run `terraform plan` (you should see import operations), then `terraform apply`.

Finally, remove any remaining `mongodbatlas_org_invitation` references from
config and state.

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
3. Run `terraform plan` followed by `terraform apply`.

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

### From `mongodbatlas_team.usernames` to `mongodbatlas_cloud_user_team_assignment`

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
# Use data source to get team members (with user_id)  
locals {
    usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
    team_assignments = {      
    for user in data.mongodbatlas_team.this.users :      
        user.id => {      
            org_id   = var.org_id
            team_id  = mongodbatlas_team.this.team_id
            user_id  = user.id
        }
    }
}

resource "mongodbatlas_team" "this" {  
    org_id = var.org_id  
    name   = var.team_name
    usernames = local.usernames
} 

data "mongodbatlas_team" "this" {  
    org_id  = var.org_id  
    team_id = mongodbatlas_team.this.team_id  
}
```

#### Step 2: Add `mongodbatlas_cloud_user_team_assignment` and use import blocks

```terraform
locals {
    usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
    team_assignments = {
    for user in data.mongodbatlas_team.this.users :
        user.id => {
            org_id   = var.org_id
            team_id  = mongodbatlas_team.this.team_id
            user_id  = user.id
        }
    }
}

resource "mongodbatlas_team" "this" {
    org_id = var.org_id
    name   = var.team_name
    usernames = local.usernames
}

data "mongodbatlas_team" "this" {
    org_id  = var.org_id
    team_id = mongodbatlas_team.this.team_id
}
  
# New resource for each (user, team) assignment  
resource "mongodbatlas_cloud_user_team_assignment" "this" {           
    for_each = local.team_assignments

    org_id  = each.value.org_id   
    team_id = each.value.team_id     
    user_id = each.value.user_id  # Use user_id instead of username  
}  
  
# Import existing team-user relationships into the new resource  
import {  
    for_each = local.team_assignments

    to = mongodbatlas_cloud_user_team_assignment.this[each.key] 
    id = "${each.value.org_id}/${each.value.team_id}/${each.value.user_id}" 
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

### Migration using Modules

If you are using modules to manage teams and user assignments to teams,
migrating from `mongodbatlas_team` to the new pattern requires some additional
steps. The main consideration is that the old `mongodbatlas_team.usernames`
attribute now maps to the resource `mongodbatlas_cloud_user_team_assignment`,
hence the `moved` block can't be used. This section demonstrates how to migrate
from a module using the `mongodbatlas_team` resource to a module using both
`mongodbatlas_team` and the new `mongodbatlas_cloud_user_team_assignment`
resources.

**Key points for module users:**

- You must use `terraform import` to bring existing user-team assignments into
  the new resources, even when they are managed inside a module.
- The import command must match the resource address as used in your module
  (e.g., `module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>`).
- If you were using a list of usernames in your previous configuration, you also
  need to include the `mongodbatlas_team` data source and use the new `users`
  attribute to retrieve the corresponding user IDs, along with team ID, for the
  import to work correctly.

**Example import blocks for modules**

```terraform
import {
   to = module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>
   # Using USER_ID
   id = "<ORG_ID>/<TEAM_ID>/<USER_ID>"
}

import {
   to = module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>
   # Using USERNAME
   id = "<ORG_ID>/<TEAM_ID>/<USERNAME>"
}
```

**Example import commands for modules:**

```shell
# Using USER_ID
terraform import 'module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>' <ORG_ID>/<TEAM_ID>/<USER_ID>

# Using USERNAME
terraform import 'module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>' <ORG_ID>/<TEAM_ID>/<USERNAME>
```

#### 1. Old Module Usage Example (Using deprecated resources)

```hcl
module "user_team_assignment" {  
  source     = "./old_module"  
  org_id     = var.org_id  
  team_name  = var.team_name  
  usernames  = var.usernames  
}
```

#### 2. New Module Usage Example (Using new resources)

```hcl
data "mongodbatlas_team" "this" {  
  org_id = var.org_id  
  name   = var.team_name
}

locals {  
  team_assigments = {
    for user in data.mongodbatlas_team.this.users :
    user.id => {
      org_id  = var.org_id
      team_id = data.mongodbatlas_team.this.team_id
      user_id = user.id
    }
  }  
}

module "user_team_assignment" {
  source     = "./new_module"
  org_id     = var.org_id
  team_name  = var.team_name
  team_assigments = local.team_assigments
}
```

#### 3. Migration Steps

1. **Add the new module to your configuration:**
   - Add the new module block as shown above, using the same input variables as
     appropriate.
   - Also add the `data.mongodbatlas_team` data source and declare the
     `team_assignments` local variable to retrieve user IDs and team ID.

2. **Import the existing user-team assignments into the new resources:**

- An `import block` (available in Terraform 1.5 and later) can be used to import
  the resource and iterate through a list of users, e.g.:

  ```terraform
  import { 
      for_each = local.team_assigments
      to       = module.user_team_assignment.mongodbatlas_cloud_user_team_assignment.this[each.key]
      id       = "${var.org_id}/${data.mongodbatlas_team.this.team_id}/${each.value.user_id}"
  }
  ```

- Alternatively, use the correct resource addresses for your module and each of
  the user-team assignments:

```shell
terraform import 'module.user_team_assignment.mongodbatlas_cloud_user_team_assignment.this' <ORG_ID>/<TEAM_ID>/<USER_ID>
```

3. **Remove the old module block from your configuration.**
4. **Run `terraform plan` to review the changes.**
   - Ensure that Terraform imports the user-team assignments and does not plan
     to create these.
   - Ensure that Terraform does not plan to destroy and recreate the
     `mongodbatlas_team` resource.
5. **Run `terraform apply` to apply the migration.**

For complete working examples, see:

- [Old module definition](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_maintainer/v1)
  and
  [old module usage](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_user/v1).
- [New module definition](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_maintainer/v2)
  and
  [new module usage](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_user_team_assignment/module_user/v2).

---

### Notes and tips

- **Import format** for `mongodbatlas_cloud_user_team_assignment`:

```
ORG_ID/TEAM_ID/USERNAME
ORG_ID/TEAM_ID/USER_ID
```

- **Importing inside modules:** Terraform import blocks cannot live inside
  modules. See
  ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)). Each
  module user must add root-level import blocks for every instance to import.

- After successful migration, ensure **no references to**
  `mongodbatlas_team.usernames` remain.

---

### FAQ

**Q: Can I assign the same user to multiple teams?** A: Yes, simply create
multiple `mongodbatlas_cloud_user_team_assignment` resources for each team.

**Q: Where can I find a working example?** A: See
[examples/mongodbatlas_cloud_user_team_assignment/main.tf](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_cloud_user_team_assignment/main.tf).

---

### Further Resources

- [Cloud User Team Assignment Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_user_team_assignment)

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

### From `mongodbatlas_project.teams` to `mongodbatlas_team_project_assignment`

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

Replace the `mongodbatlas_project.teams` block with:

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

Then run:

```shell
terraform plan  
terraform apply
```

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

Run `terraform plan` (you should see **import** operations), then
`terraform apply`.

#### Step 3: Verify and clean up

- After successful import and apply, `terraform plan` should show **no
  changes**.
- Keep the `ignore_changes = ["teams"]` lifecycle rule until the provider
  releases a version without the `teams` argument in `mongodbatlas_project`.

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

- **Modules:** Terraform import blocks cannot live inside modules
  ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)).
- If you manage team assignments in modules, import each at the root level using
  the correct resource address (e.g.
  `module.<name>.mongodbatlas_team_project_assignment.<name>`).
- You can use `terraform plan` to confirm imports before applying.

---

### FAQ

**Q: Do I need to delete the old `teams` from state?** A: No — using
`ignore_changes` ensures they remain in Atlas until the provider removes the
field. Then you can drop the lifecycle rule.

---

### Further resources

- [`mongodbatlas_team_project_assignment` docs](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/team_project_assignment)

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

#### Step 1: Add the new resource alongside existing configuration

Add the new resource to re-create the pending invite via the new API:

```terraform
resource "mongodbatlas_cloud_user_project_assignment" "this" {
  project_id = var.project_id
  username   = local.username
  roles      = local.roles
}
```

Use the same `roles` as the original invitation to avoid drift.

#### Step 2: Remove the deprecated resource from the configuration and state

##### Option A) [Recommended] Removed block

Remove the resource block and replace it with a `removed` block to cleanly
remove the old resource from state:

```terraform
removed {
  from = mongodbatlas_project_invitation.this

  lifecycle {
    destroy = false
  }
}
```

##### Option B) Manual state removal

Remove the `mongodbatlas_project_invitation` resource from configuration and
then remove it from the Terraform state using the command line (this does not
affect the actual invitation in Atlas):

```bash
terraform state rm mongodbatlas_project_invitation.this
```

#### Step 3: Apply the changes

Run `terraform apply` to create the assignment with the new resource.
Afterwards, run `terraform plan` and ensure no further changes are pending.

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
