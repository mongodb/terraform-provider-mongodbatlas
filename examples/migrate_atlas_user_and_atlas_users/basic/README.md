# Basic Migration Example: Atlas User Data Sources

This example demonstrates how to migrate from the deprecated `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users` data sources to their replacements.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

### v1: Initial State (Deprecated Data Sources)

Shows the original configuration using deprecated data sources:
- `mongodbatlas_atlas_user` for single user reads
- `mongodbatlas_atlas_users` for user lists

### v2: Migration Phase

Runs both old and new data sources side-by-side to validate the migration.

**Key comparisons:**
- Single user: `email_address` → `username`
- Single user: Complex role filtering → Structured `roles.org_roles` and `roles.project_role_assignments`
- User lists: `results[*].email_address` → `users[*].username`

**Usage:**
1. Apply this configuration: `terraform apply`
2. Review the comparison outputs to verify data consistency
3. Check `migration_validation.ready_for_v3` is `true`
4. Once validated, proceed to v3

**Expected outputs:** Side-by-side comparisons of email retrieval methods, role access patterns, user list structures, and count validations.

### v3: Final State (New Data Sources Only)

Clean final configuration using only:
- `mongodbatlas_cloud_user_org_assignment` for single user reads
- `mongodbatlas_organization.users`, `mongodbatlas_project.users`, `mongodbatlas_team.users` for user lists

**Key improvements:**
- **Structured roles**: Organization and project roles are clearly separated
- **Direct access**: No need to filter consolidated role lists
- **Consistent naming**: `username` instead of `email_address`
- **Better organization**: User lists come from their natural containers (org/project/team)
- **Better performance**: Organization context required for user reads (more efficient API calls)

## Usage

1. Start with v1 to understand the original setup
2. Apply v2 to validate new data sources return expected data
3. Apply v3 for the final clean state

## Prerequisites

- MongoDB Atlas Terraform Provider 2.0.0 or later
- Valid MongoDB Atlas organization, project, and team IDs
- Existing users in your organization

## Variables

Set these variables for all versions:

```terraform
client_id     = "<ATLAS_CLIENT_ID>"      # Optional, can use env vars
client_secret = "<ATLAS_CLIENT_SECRET>"  # Optional, can use env vars
org_id        = "your-organization-id"
project_id    = "your-project-id"  
team_id       = "your-team-id"
user_id       = "existing-user-id"
username      = "existing-user@example.com"
```

Alternatively, set environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```
