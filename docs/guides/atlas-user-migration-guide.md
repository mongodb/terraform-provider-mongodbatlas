---
page_title: "Migration Guide: Migrate off deprecated `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users`"
---

# Migration Guide: Migrate off deprecated `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users`

**Objective**: Migrate from the deprecated `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users` data sources to their respective replacements.

## Before you begin

- Ensure you are using MongoDB Atlas Terraform Provider version 2.0.0 or later that includes `mongodbatlas_cloud_user_org_assignment`.

## What’s changing?

- `mongodbatlas_atlas_user` returned a user profile by `user_id` or `username` and is deprecated.
- `mongodbatlas_cloud_user_org_assignment` reads a user’s assignment in a specific organization using either `username` or `user_id` together with `org_id`.
- For details on the new data source, see the `mongodbatlas_cloud_user_org_assignment` data source [documentation](../data-sources/cloud_user_org_assignment)

- `mongodbatlas_atlas_users` returned lists of users by `org_id`, `project_id`, or `team_id` and is deprecated. Replace it with the `users` attribute available on `mongodbatlas_organization`, `mongodbatlas_project`, or `mongodbatlas_team` data sources, respectively.
- Attribute structure differences: The new organization users API does not return `email_address` as a separate field and replaces the consolidated `roles` with structured `org_roles` and `project_role_assignments`.

---

## Migrate reads to `mongodbatlas_cloud_user_org_assignment`

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

### Step 1: Add the new data source alongside the existing one

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

### Step 2: Verify the new data source works

Run `terraform plan` to ensure the new data source will read correctly without errors.

### Step 3: Replace references incrementally

Replace references from `data.mongodbatlas_atlas_user.test` to `data.mongodbatlas_cloud_user_org_assignment.user_1`.

**Important**: Update attribute references as the structure has changed:

Key attribute changes:

| Old Attribute | New Attribute |
|---------------|---------------|
| `email_address` | `username` |
| `roles` (filtered by org_id) | `roles.org_roles` |
| `roles` (filtered by group_id) | `roles.project_role_assignments[*].project_roles` |

**Examples**:
- Email: `data.mongodbatlas_atlas_user.test.email_address` → `data.mongodbatlas_cloud_user_org_assignment.user_1.username`
- Org roles: Use `data.mongodbatlas_cloud_user_org_assignment.user_1.roles.org_roles` directly
- Project roles: Access via `roles.project_role_assignments` list, filtering by `project_id` as needed

### Step 4: Remove the old data source

Once all references are updated and working, remove the old data source from your configuration:

```terraform
# Remove this block
# data "mongodbatlas_atlas_user" "test" {
#   user_id = "<USER_ID>"
# }
```

### Step 5: Apply and verify

Run `terraform plan` to ensure no unexpected changes, then `terraform apply`.

---

## Migrate list reads from `mongodbatlas_atlas_users`

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

### Step 1: Add new data sources alongside existing ones

Add the appropriate replacement data source(s) while keeping the old one temporarily:

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

### Step 2: Verify new data sources work

Run `terraform plan` to ensure the new data sources read correctly and return expected user data.

### Step 3: Replace references incrementally

Replace `data.mongodbatlas_atlas_users.test.results` with the appropriate `...users` collection above.

**Important**: Update attribute references as the structure has changed:

| Old Attribute | New Attribute |
|---------------|---------------|
| `results[*].email_address` | `users[*].username` |
| `results[*].roles` (filtered) | `users[*].roles.org_roles` or `users[*].roles` |

**Examples**:
- Email list: `data.mongodbatlas_atlas_users.test.results[*].email_address` → `data.mongodbatlas_organization.org.users[*].username`
- User list: `data.mongodbatlas_atlas_users.test.results` → `data.mongodbatlas_organization.org.users` (or `.project.proj.users`, `.team.team.users`)
- Org roles: Use `users[*].roles.org_roles` from organization data source
- Project roles: Use `users[*].roles` from project data source, or `users[*].roles.project_role_assignments` from organization data source

### Step 4: Remove the old data source

Once all references are updated and working, remove the old data source from your configuration:

```terraform
# Remove this block
# data "mongodbatlas_atlas_users" "test" {
#   org_id = "<ORG_ID>"
# }
```

### Step 5: Apply and verify

Run `terraform plan` to ensure no unexpected changes, then `terraform apply`.

---

## Examples

For complete, working configurations that demonstrate the migration process, see the examples in the provider repository: [migrate_atlas_user_and_atlas_users](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_atlas_user_and_atlas_users). 

The examples include:
- **v1**: Original configuration using deprecated data sources
- **v2**: Migration phase with side-by-side comparison and validation
- **v3**: Final clean configuration using only new data sources

These examples provide practical validation of the migration steps and demonstrate the attribute mappings in working Terraform code.

---

## Notes

- The new data source requires the `org_id` context to read the user's organization assignment.
- After migration, ensure no remaining references to `mongodbatlas_atlas_user` exist in your configuration.
