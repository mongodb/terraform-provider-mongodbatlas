# Combined Example: Org Invitation → Cloud User Org Assignment

This combined example is organized into step subfolders (v1–v3):

- v1/: Initial state with:
  - a pending `mongodbatlas_org_invitation` (with `teams_ids`), and
  - an accepted (ACTIVE) user present in the org (no invitation in state).
- v2/: Migration step showcasing both paths:
  - moved block for the pending invitation (module-friendly, recommended), and
  - import blocks for accepted (ACTIVE) users and team assignments.
- v3/: Cleaned-up final configuration after v2 is applied:
  - remove the `mongodbatlas_org_invitation` resource,
  - remove moved and import blocks,
  - keep only `mongodbatlas_cloud_user_org_assignment` and `mongodbatlas_cloud_user_team_assignment`.

Navigate into each version folder to see the step-specific configuration.
