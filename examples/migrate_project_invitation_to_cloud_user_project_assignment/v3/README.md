# v3: Final State

This is the clean, final configuration using only the new `mongodbatlas_cloud_user_project_assignment` resource.

## What changed from v1

### Resource purpose
- **Old**: Managed pending invitations only
- **New**: Manages active project membership

### Lifecycle behavior
- **Old**: Removed from state when user accepted invitation
- **New**: Remains in state, manages ongoing membership

### Available data
- **Old**: Only invitation details (invitation_id, expires_at, etc.)
- **New**: User assignment details including user_id

### Data source support
- **Old**: Had data source for reading invitation details
- **New**: Has data source for reading user assignments

## Key improvements

1. **Persistent management**: Resource doesn't disappear when user accepts invitation
2. **User ID access**: Provides user_id for use in other resources
3. **Import support**: Can import existing project members
4. **Cleaner lifecycle**: No surprise state removals

## Enhanced functionality

The new resource provides additional capabilities:

- **Data source**: Read existing user assignments
- **Import**: Bring existing project members under Terraform management
- **User ID exposure**: Reference users by ID in other resources
- **Active membership**: Manage actual project membership, not just invitations

## Usage patterns

This configuration demonstrates:
- Basic user assignment to project
- Data source usage for reading assignments
- Local values for organizing assignment data
- Output examples showing common use cases

## Migration complete

At this point, you have successfully migrated from the deprecated `mongodbatlas_project_invitation` resource to the modern `mongodbatlas_cloud_user_project_assignment` resource. All references to the old resource have been removed and replaced with the new resource.
