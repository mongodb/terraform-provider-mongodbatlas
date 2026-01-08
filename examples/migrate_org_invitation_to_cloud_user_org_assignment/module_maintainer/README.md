# Module Maintainer Migration Example

This example demonstrates how module maintainers should migrate from `mongodbatlas_org_invitation` to `mongodbatlas_cloud_user_org_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

This example demonstrates the legacy pattern (prior to v2.0.0) for managing org invitations using the the `mongodbatlas_org_invitation` resource. It is intended to show the "before" state for users migrating to the new recommended pattern.

## v2: Migration

- The `moved` block handles migration from `mongodbatlas_org_invitation` to `mongodbatlas_cloud_user_org_assignment`.
- Team assignments use the new `mongodbatlas_cloud_user_team_assignment` resource.

Team assignments must be imported at the root level. Import ID format: `ORG_ID/TEAM_ID/USER_ID` (or `ORG_ID/TEAM_ID/USERNAME`).

See `module_user/v2` for an example of how to consume this module with imports.

## v3: Cleaned Up Configuration

Final module definition after migration is complete:
- The `moved` block from v2 has been removed since all users have migrated.
- Uses `mongodbatlas_cloud_user_org_assignment` and `mongodbatlas_cloud_user_team_assignment` resources.

