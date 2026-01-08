# Basic Migration Example

This example demonstrates direct resource usage for migrating from `mongodbatlas_org_invitation` to `mongodbatlas_cloud_user_org_assignment`.

For migration steps, see the [Migration Guide](../../../docs/guides/atlas-user-management.md).

## v1: Initial State

- `mongodbatlas_org_invitation` manages a pending user (with `teams_ids`).
- An accepted (ACTIVE) user exists in the organization (no invitation in state), referenced via `data.mongodbatlas_organization`.

## v2: Migration

- Pending invitation → move state from `mongodbatlas_org_invitation` to `mongodbatlas_cloud_user_org_assignment` using a Terraform `moved` block (no recreate).
- Accepted (ACTIVE) user → declare the resource and use `import` blocks to adopt the existing assignment (`org_id,user_id`).
- Teams → manage memberships via `mongodbatlas_cloud_user_team_assignment`; import existing mappings (`org_id,team_id,user_id`).

## v3: Cleaned Up Configuration

Final configuration after migration:
- Only `mongodbatlas_cloud_user_org_assignment` and (optionally) `mongodbatlas_cloud_user_team_assignment` remain.
- No `mongodbatlas_org_invitation` resources.
- No import blocks.
