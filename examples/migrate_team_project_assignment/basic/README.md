# Basic Migration Example

This example demonstrates direct resource usage for migrating from `mongodbatlas_project.teams` to `mongodbatlas_team_project_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

Uses the deprecated `mongodbatlas_project.teams` inline block to assign teams to the project.

## v2: Migration

- Add `ignore_changes = [teams]` lifecycle rule to the project.
- Define `mongodbatlas_team_project_assignment` resources for each team.
- Add `import` blocks to import existing team-project assignments.
- Run `terraform plan` — expect `will be imported` for each team assignment.

## v3: Cleaned Up Configuration

Final configuration after migration:
- Uses only `mongodbatlas_team_project_assignment` resources.
- Keep `ignore_changes = [teams]` until the provider removes the teams attribute.
- No import blocks.
