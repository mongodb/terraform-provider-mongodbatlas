# v2: Migrate using Moved and Import blocks

State:
- Pending invitation → move state from `mongodbatlas_org_invitation` to `mongodbatlas_cloud_user_org_assignment` using a Terraform `moved` block (no recreate).
- Accepted (ACTIVE) user → declare the resource and use `import` blocks to adopt the existing assignment (`org_id,user_id`).
- Teams → manage memberships via `mongodbatlas_cloud_user_team_assignment`; import existing mappings (`org_id,team_id,user_id`).
