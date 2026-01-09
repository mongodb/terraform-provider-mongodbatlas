# Module User Migration Example

This example demonstrates how module consumers should migrate when their module updates from `mongodbatlas_project_invitation` to `mongodbatlas_cloud_user_project_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

Using the legacy module version that relies on `mongodbatlas_project_invitation`.

## v2: Migration

Update module source to v2. Terraform will plan to delete the old `mongodbatlas_project_invitation` and create the new `mongodbatlas_cloud_user_project_assignment`.

## v3: Cleaned Up Configuration

Using the final module version with only `mongodbatlas_cloud_user_project_assignment`.
