# Example layout reference

Read this with `pr-and-documentation-standards`. Canonical family: `examples/mongodbatlas_cloud_backup_collection_restore_job/`.

## Canonical sibling tree

```
examples/mongodbatlas_cloud_backup_collection_restore_job/
  README.md
  snapshot_discovery/
  snapshot_restore/
  pit_restore/
  restore_job_inspection/
```

Each sibling is its own root: `main.tf` (or named `.tf` when one sibling has more than one surface), `variables.tf`, `providers.tf`, `versions.tf`, `README.md`. No parent `main.tf`.

## Atlas project resource

When the example creates a project, use this comment. Do not repeat the resource type in the name (`"project"` on `mongodbatlas_project`). Use `this` when the root has a single resource of that type, or when no more descriptive name exists. Use a descriptive name only when there are several of the same type.

```hcl
# Set up MongoDB Atlas Project access
resource "mongodbatlas_project" "this" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}
```

Collection restore does not create a project; it points at an existing cluster.

## Common mistakes

- Generic `singular-data-source.tf` / `plural-data-source.tf` names for new work.
- Shared module or `terraform_remote_state` between siblings.
- Pinning `mongodbatlas` version or putting keys in the provider block.
- README repeating Atlas product limits instead of linking.
- Outputs or variables missing `description`.
- `main.tf` comment that describes locals instead of the README goal.
- Inlined HCL in a `.md.tmpl`.
- Sibling directories when there is only one flow.
- Naming a single `mongodbatlas_project` `project` (repeats the type). New examples use `this`. Leave existing `"project"` labels alone.

## Exceptions

- **`restore_job_inspection/`**: Uses `restore_job.tf` and `restore_job_collections.tf` instead of `main.tf` because one sibling covers two data-source surfaces. The goal comment still sits at the top of the primary file and matches the README opening line.
- **Old examples**: `provider.tf` and generic DS filenames already in the repo stay as they are. Do not migrate them while adding a new resource.
- **Existing inlined tmpls**: For example `templates/data-sources/ai_model_org_api_keys.md.tmpl`. Do not retrofit unless the task is to fix that template.
