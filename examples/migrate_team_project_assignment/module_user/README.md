# Module User Migration Example

This example demonstrates how module consumers should migrate when upgrading to a module that uses `mongodbatlas_team_project_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

Uses the legacy module which manages `mongodbatlas_project.teams` inline block.

## v2: Migration

- Upgrade to the new module version (`terraform init -upgrade`).
- Add `import` blocks at the root level for each team-project assignment.
- Import ID format: `PROJECT_ID/TEAM_ID`.
- Run `terraform plan` — expect `will be imported` for each team assignment.

## v3: Cleaned Up Configuration

Final module consumption after migration is complete:
- Uses the new module with `mongodbatlas_team_project_assignment`.
- The `import` blocks from v2 have been removed since migration has been applied.
