# v1: Initial State

State:
- `mongodbatlas_org_invitation` manages a pending user (with `teams_ids`).
- An accepted (ACTIVE) user exists in the organization (no invitation in state), referenced via `data.mongodbatlas_organization`.
