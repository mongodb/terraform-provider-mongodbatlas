# Module Maintainer Migration Example

This example demonstrates how module maintainers should migrate from `mongodbatlas_project_invitation` to `mongodbatlas_cloud_user_project_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

This example demonstrates the legacy pattern (prior to v2.0.0) for managing project invitations using the `mongodbatlas_project_invitation` resource. It is intended to show the "before" state for users migrating to the new recommended pattern.

## v2: Migration

- Replace `mongodbatlas_project_invitation` with `mongodbatlas_cloud_user_project_assignment`.
- Terraform will plan to delete the old resource and create the new one.

See `module_user/v2` for an example of how to consume this module.

## v3: Cleaned Up Configuration

Final module definition after migration is complete:
- The `moved` block from v2 has been removed since all users have migrated.
- Uses `mongodbatlas_cloud_user_project_assignment` resource.
