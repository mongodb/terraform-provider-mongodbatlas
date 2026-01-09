# Module User Migration Example

This example demonstrates how module consumers should migrate when their module updates from `mongodbatlas_atlas_user` to `mongodbatlas_cloud_user_org_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

Uses the legacy module version that relies on `mongodbatlas_atlas_user`:
- Exposes all module outputs (`user_id`, `username`, `email_address`, `first_name`, `last_name`, `org_roles`, `project_roles`)

## v2: Migration

Update module source to v2. Since data sources don't have state:
- Simply change module source from v1 to v2
- No state migration needed
- Output names remain the same
- Note: `email_address` now returns the same value as `username`

Run `terraform plan` to verify outputs are correct.

## v3: Cleaned Up Configuration

Using the final module version with only `mongodbatlas_cloud_user_org_assignment`:
- Clean configuration
- Same outputs as before
