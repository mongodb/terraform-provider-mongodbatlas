# Module Maintainer Migration Example

This example demonstrates how module maintainers should migrate from `mongodbatlas_atlas_user` data source to `mongodbatlas_cloud_user_org_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

Module using the deprecated `mongodbatlas_atlas_user` data source:
- Complex role filtering with for expressions
- Separate `email_address` and `username` attributes
- Outputs: user_id, username, email_address, first_name, last_name, org_roles, project_roles

## v2: Migration

- Replace `mongodbatlas_atlas_user` with `mongodbatlas_cloud_user_org_assignment`
- Simplified role access via `roles.org_roles` (no filtering needed for org roles)
- `email_address` now returns `username`
- Output names kept the same

See `module_user/v2` for an example of how to consume this module.

## v3: Cleaned Up Configuration

Final module definition after migration is complete:
- Uses `mongodbatlas_cloud_user_org_assignment` data source
- Same outputs as v1/v2
