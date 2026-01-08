# Module User Migration Example

This example demonstrates how module consumers should migrate when upgrading to a module that uses `mongodbatlas_cloud_user_org_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

Uses the legacy module which manages `mongodbatlas_org_invitation` with `teams_ids`.

## v2: Migration

- Upgrade to the new module version (`terraform init -upgrade`).
- The `moved` block in the module handles the org assignment migration automatically.
- Add `import` blocks for team assignments at the root level.
- Run `terraform plan` — expect `has moved to` for org assignment and `will be imported` for team assignments.

## v3: Cleaned Up Configuration

Final module consumption after migration is complete:
- Only `mongodbatlas_cloud_user_org_assignment` remains in state.
- The `import` blocks from v2 have been removed since migration has been applied.
