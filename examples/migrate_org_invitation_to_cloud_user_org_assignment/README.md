# Migration Example: Org Invitation → Cloud User Org Assignment

This example demonstrates migrating from `mongodbatlas_org_invitation` to `mongodbatlas_cloud_user_org_assignment`.

## Basic usage (direct resource usage)

- `basic/v1`: Initial state with `mongodbatlas_org_invitation` (with `teams_ids`).
- `basic/v2`: Migration step with `moved` block and `import` blocks for team assignments.
- `basic/v3`: Final clean configuration using only `mongodbatlas_cloud_user_org_assignment` and `mongodbatlas_cloud_user_team_assignment`.

## Module-based examples

- `module_maintainer/v1`: Legacy module using `mongodbatlas_org_invitation` with `teams_ids`.
- `module_maintainer/v2`: Migrated module with `moved` block.
- `module_maintainer/v3`: Final module (no `moved` block needed).
- `module_user/v1`: Legacy module consumption.
- `module_user/v2`: Migration with root-level imports for team assignments.
- `module_user/v3`: Final clean configuration (no import blocks needed).

Navigate into each version folder to see the step-specific configuration.
