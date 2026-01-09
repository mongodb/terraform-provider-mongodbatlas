# Migration Example: Project Invitation to Cloud User Project Assignment

This example demonstrates how to migrate from the deprecated `mongodbatlas_project_invitation` resource to the new `mongodbatlas_cloud_user_project_assignment` resource.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

### v1: Initial State (Deprecated Resource)

Shows the original configuration using deprecated `mongodbatlas_project_invitation` for pending invitations.

### v2: Migration Phase (Re-creation with Removed Block)

Demonstrates the migration approach:
- Adds new `mongodbatlas_cloud_user_project_assignment` resource
- Uses `removed` block to cleanly remove old resource from state
- Shows both removed block and manual state removal options

**Key differences:**

Old resource (`mongodbatlas_project_invitation`):
- Only managed pending invitations
- Removed from state when user accepted invitation
- Limited to invitation lifecycle only

New resource (`mongodbatlas_cloud_user_project_assignment`):
- Manages active project membership
- Exposes `user_id` (not available in old resource)
- Supports import for existing users
- Works for both pending and active users

**Migration approach:**
1. **Add new resource**: Re-creates the pending invitation with new API
2. **Remove old resource**: Uses `removed` block to clean up Terraform state
3. **Validate**: Check that the new resource works as expected

**Alternative removal method:**

If you prefer manual state removal instead of the `removed` block:

```bash
# Remove from configuration first, then:
terraform state rm mongodbatlas_project_invitation.pending_user
```

**Usage:**
1. Apply this configuration: `terraform apply`
2. Review the validation outputs to ensure migration success
3. Check `migration_validation.ready_for_v3` is `true`
4. Once validated, proceed to v3

**Expected behavior:**
- The user should receive a new invitation email (since we're re-creating the invitation)
- The old invitation remains valid until it expires or is accepted
- Terraform state is cleaned up properly

### v3: Final State (New Resource Only)

Clean final configuration using only `mongodbatlas_cloud_user_project_assignment`.

**Key improvements:**
- **Persistent management**: Resource doesn't disappear when user accepts invitation
- **User ID access**: Provides user_id for use in other resources
- **Import support**: Can import existing project members
- **Cleaner lifecycle**: No surprise state removals

**Enhanced functionality:**
- **Data source**: Read existing user assignments
- **Import**: Bring existing project members under Terraform management
- **User ID exposure**: Reference users by ID in other resources
- **Active membership**: Manage actual project membership, not just invitations

## Important Notes

- **Pending invites only**: This migration applies only to PENDING project invitations that still exist in your Terraform configuration
- **Re-creation approach**: The new resources and data sources cannot discover pending invites created by the deprecated resource, so we re-create them with the new resource
- **Accepted invites**: If users already accepted invitations, the provider removed them from state and you should remove them from configuration (no migration needed)

## Usage

1. Start with v1 to understand the original setup with pending invitations
2. Apply v2 configuration to re-create invites with new resource and remove old resource
3. Apply v3 configuration for the final clean state

## Prerequisites

- MongoDB Atlas Terraform Provider 2.0.0 or later
- Valid MongoDB Atlas project ID
- Pending project invitations in your configuration (not yet accepted)

## Variables

Set these variables for all versions:

```terraform
client_id     = "your-service-account-client-id"   # Optional, can use env vars
client_secret = "your-service-account-client-secret"  # Optional, can use env vars
project_id  = "your-project-id"
username    = "user@example.com"                # User with pending invitation
roles       = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
```

Alternatively, set environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```
