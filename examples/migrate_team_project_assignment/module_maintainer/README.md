# Module Maintainer Migration Example

This example demonstrates how module maintainers should migrate from `mongodbatlas_project.teams` to `mongodbatlas_team_project_assignment`.

For migration steps, see the [Migration Guide](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/guides/atlas-user-management.md).

## v1: Initial State

This example demonstrates the legacy pattern (prior to v2.0.0) for managing team-project assignments using the deprecated `mongodbatlas_project.teams` block. It is intended to show the "before" state for users migrating to the new recommended pattern.

## v2: Migration

- Add `ignore_changes = [teams]` lifecycle rule to the project.
- Add `mongodbatlas_team_project_assignment` resources for each team.
- Terraform doesn't allow import blocks in modules. Document import ID format for users: `PROJECT_ID/TEAM_ID`.
- Expose `project_id` as a module output so users can form import IDs.

See `module_user/v2` for an example of how to consume this module with imports.

## v3: Cleaned Up Configuration

Final module definition after migration is complete:
- Uses `mongodbatlas_team_project_assignment` resources.
- Keep `ignore_changes = [teams]` until the provider removes the teams attribute.
