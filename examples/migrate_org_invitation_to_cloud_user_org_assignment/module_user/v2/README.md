# module_user/v2 (migrated)

Module consumption after migrating to `mongodbatlas_cloud_user_org_assignment` (and `mongodbatlas_cloud_user_team_assignment` when teams are used).

Steps:
1) Copy state from v1: `cp ../v1/terraform.tfstate .`
2) Set credentials/vars (see `variables.tf`).
3) `terraform init`.
4) `terraform plan` — you should see:
   - `has moved to` for the org assignment (handled by the `moved` block in the module)
   - `will be imported` for team assignments (import blocks are defined in main.tf)
5) `terraform apply`.
6) Verify with `terraform plan` — should show no changes.

Note: The org assignment migration is handled automatically by the `moved` block inside the module. Only team assignments require root-level imports.

Refer back to `docs/guides/atlas-user-management.md` for full migration context and ID formats.

