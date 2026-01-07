# v2: Migration

State:
- The `moved` block handles migration from `mongodbatlas_org_invitation` to `mongodbatlas_cloud_user_org_assignment`.
- Team assignments use the new `mongodbatlas_cloud_user_team_assignment` resource.

Team assignments must be imported at the root level. Import ID format: `ORG_ID/TEAM_ID/USER_ID` (or `ORG_ID/TEAM_ID/USERNAME`).

See `module_user/v2` for an example of how to consume this module with imports.

