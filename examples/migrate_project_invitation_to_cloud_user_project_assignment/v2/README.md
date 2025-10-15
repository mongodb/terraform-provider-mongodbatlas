# v2: Migration Phase

This configuration demonstrates the migration approach using the `removed` block to cleanly transition from `mongodbatlas_project_invitation` to `mongodbatlas_cloud_user_project_assignment`.

## What this shows

- **Re-creation**: The new resource re-creates the pending invitation using the new API
- **Clean removal**: Uses `removed` block to remove old resource from state without destroying the Atlas invitation
- **Validation**: Outputs that verify the new resource works correctly
- **New capabilities**: Shows additional features available with the new resource

## Key differences

### Old resource (`mongodbatlas_project_invitation`)
- Only managed pending invitations
- Removed from state when user accepted invitation
- Limited to invitation lifecycle only

### New resource (`mongodbatlas_cloud_user_project_assignment`)
- Manages active project membership
- Exposes `user_id` (not available in old resource)
- Supports import for existing users
- Works for both pending and active users

## Migration approach

1. **Add new resource**: Re-creates the pending invitation with new API
2. **Remove old resource**: Uses `removed` block to clean up Terraform state
3. **Validate**: Check that the new resource works as expected

## Alternative removal method

If you prefer manual state removal instead of the `removed` block:

```bash
# Remove from configuration first, then:
terraform state rm mongodbatlas_project_invitation.pending_user
```

## Usage

1. Apply this configuration: `terraform apply`
2. Review the validation outputs to ensure migration success
3. Check `migration_validation.ready_for_v3` is `true`
4. Once validated, proceed to v3

## Expected behavior

- The user should receive a new invitation email (since we're re-creating the invitation)
- The old invitation remains valid until it expires or is accepted
- Terraform state is cleaned up properly
