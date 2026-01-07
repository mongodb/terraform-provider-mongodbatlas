# module_user/v1 (legacy)

Legacy module consumption using `mongodbatlas_org_invitation` with `teams_ids`. No imports are needed because the deprecated resource already manages state.

Steps:
1) Set credentials/vars (see `variables.tf`).
2) `terraform init && terraform plan && terraform apply`.
3) Use v1 only as the “before” state; migrate to v2 to adopt the new resources.

