---
name: pr-and-documentation-standards
description: Use when creating PRs, writing or editing documentation (schema descriptions, docs markdown, examples), creating or restructuring examples/, reviewing code, or adding changelog entries.
---

# PR and Documentation Standards

## Pull Request Structure

### PR Title Format

PR titles must follow [Conventional Commits](https://www.conventionalcommits.org/): `fix:`, `feat:`, `chore:`, `doc:`, `test:`, `refactor:`, `ci:`, etc. The subject after the colon must start with an uppercase letter.

```text
fix: Emit warning when use_effective_fields and auto-scaling are enabled
feat: Add search deployment resource
```

Do **not** include a scope in parentheses (e.g. `fix(resource/mongodbatlas_advanced_cluster): ...`). The repo convention uses the bare prefix without a scope.

### Separate Refactoring from Feature Changes

Avoid mixing refactoring with functional changes in the same PR. Reviewers need to clearly distinguish which changed lines are behavioral vs structural.

### Changelog Entries

Add a changelog entry (`.changelog/<PR_NUMBER>.txt`) for:
- Bug fixes (`release-note:bug`)
- New features (`release-note:enhancement`)
- Breaking changes (`release-note:breaking-change`)
- New resources/data sources (`release-note:new-resource` / `release-note:new-data-source`)
- Migration guides or user-facing documentation changes

### Changelog Entry Format

- Start with a verb in 3rd person singular (e.g. `Emits`, `Adds`, `Fixes`, `Updates`)
- Do not end with punctuation

## Documentation Style Guide

### Consolidate Admonitions

Avoid excessive NOTE/IMPORTANT/WARNING boxes. Prefer:
1. Inlining short notes into attribute descriptions.
2. Combining multiple notes into a single box.
3. Downgrading from IMPORTANT to NOTE when the content is informational, not action-required.

### CLOUDP Ticket References

Do **not** include CLOUDP ticket references in user-facing documentation. Internal ticket references are acceptable in code comments only when tracking a deliberate technical decision.

### Resource and Data Source Descriptions

Start data source and resource descriptions with the resource name and a clear one-line purpose:

```
`mongodbatlas_log_integration` provides a resource for managing log integration configurations at the project level.
```

## Examples (`examples/` directory)

The following rules are the default for new provider examples. One user journey (flow) stays a flat root, same shape as `examples/mongodbatlas_ai_model_api_key/`. Two or more flows use sibling directories. Copy `examples/mongodbatlas_cloud_backup_collection_restore_job/` for that layout. Read `example-layout.md` in this skill directory for the canonical tree, Atlas project comment, common mistakes, and exceptions.

### Flat vs siblings

- **One flow**: Keep files at `examples/mongodbatlas_<name>/` (`main.tf`, `variables.tf`, `providers.tf`, `versions.tf`, `README.md`).
- **Two or more flows**: Use a parent `README.md` (index only, no parent `main.tf`) plus sibling directories named for what they do.

### Rules

1. **Use-case directories**: Name each example for what it does (`snapshot_restore`, `pit_restore`). Do not use generic `singular-data-source.tf` / `plural-data-source.tf`.
2. **Standalone configs**: Each subdirectory is its own root module and depends only on input variables. Do not share Terraform modules across siblings. Do not use `terraform_remote_state`. Document cross-example flow in the grouped README; users copy outputs into `terraform.tfvars`.
3. **Aligned intro**: Each primary `.tf` file starts with a one- or two-sentence comment stating the goal. Match the sibling README opening line (the first sentence after the H1, not the H1 itself).
4. **Grouped root README**: Parent `README.md` lists siblings, typical flows, and links to the product doc. Collection restore links to [Restore from Selected Databases and Collections](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-db-coll/). Other resources link their own product page.
5. **Docs over duplication**: Keep READMEs operational (prerequisites, defaults, tfvars). Link to official MongoDB docs for product limits and semantics.
6. **Variables**: Every user-facing value is a `variable` with a `description`. Put defaults in `variables.tf` and document them in the README.
7. **Outputs**: Expose only useful post-apply values (`job_state`, `collection_states`, discovery IDs). Each output has a `description`.
8. **Template docs**: Never paste HCL into `templates/**/*.md.tmpl`. Embed with [tffile](https://github.com/hashicorp/terraform-plugin-docs?tab=readme-ov-file#templates) pointing at a real file under `examples/`, for example `{{ tffile "examples/mongodbatlas_cloud_backup_collection_restore_job/snapshot_restore/main.tf" }}`. Further Examples links point at sibling directories.

### Provider baseline

State once, not per sibling:

- No pinned provider version.
- Empty `provider "mongodbatlas" {}`. Credentials come from the environment, not from HCL. Do not set `client_id`, `client_secret`, `public_key`, or `private_key` on the provider.
- README `export` examples use Service Account `MONGODB_ATLAS_CLIENT_ID` / `MONGODB_ATLAS_CLIENT_SECRET`. Prefer those. Programmatic API keys via env also work.
- `versions.tf` sets `source = "mongodb/mongodbatlas"` and `required_version` only.

New examples use `providers.tf`. Leave existing `provider.tf` names alone.
