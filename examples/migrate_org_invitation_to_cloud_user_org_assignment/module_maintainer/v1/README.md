# module_maintainer/v1 (legacy)

Legacy module definition using `mongodbatlas_org_invitation` with `teams_ids` (deprecated). Use this only as the “before” state; migrate to v2 to adopt `mongodbatlas_cloud_user_org_assignment` and `mongodbatlas_cloud_user_team_assignment`.

Contents: `main.tf`, `variables.tf`, `versions.tf`.

Next step: switch to `module_maintainer/v2` and follow the migration guide (`docs/guides/atlas-user-management.md`) to publish the updated module. Imports must be performed by module consumers at the root.

