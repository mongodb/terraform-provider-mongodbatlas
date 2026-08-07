---
name: pr-and-documentation-standards
description: Standards for pull requests, documentation, and code review in this Terraform provider. Use when creating PRs, writing or editing documentation (schema descriptions, docs markdown, examples), reviewing code, or adding changelog entries. Covers PR structure, docs style guide, example conventions, and changelog practices.
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
- New resources/data sources (`release-note:new-resource` / `release-note:new-datasource`)
- Migration guides or user-facing documentation changes

### Changelog Entry Format

- For `new-resource` and `new-datasource`, include only the Terraform name, such as `mongodbatlas_cluster_adaptive_settings`. Do not add `resource/` or `data-source/` prefixes, a verb, or punctuation.
- For other entry types, use the required resource, data-source, or provider prefix, then start the description with a third-person singular verb such as `Emits`, `Adds`, `Fixes`, or `Updates`.
- Do not end descriptions with punctuation.
- Treat `contributing/changelog-process.md` as authoritative when examples or older entries disagree.

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

### Do Not Pin Provider Versions

Do not pin specific provider versions in examples. This avoids examples becoming outdated and ensures users always get the latest compatible version.

### Use Variables Consistently

All configurable values in examples should use `var.` references with corresponding entries in `variables.tf`. Do not hardcode values that users need to customize.

### Atlas Project Resource Comment

Use a consistent comment for the Atlas project resource block in examples:

```hcl
# Set up MongoDB Atlas Project access
resource "mongodbatlas_project" "project" {
  name   = var.atlas_project_name
  org_id = var.atlas_org_id
}
```
