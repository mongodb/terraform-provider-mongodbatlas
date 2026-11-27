---
name: mongodbatlas-3-0-upgrade
description: Use when upgrading Terraform configurations from MongoDB Atlas provider (mongodb/mongodbatlas) 2.x to 3.0, or when `terraform plan` fails after the upgrade with "Unsupported block type", "Unsupported argument", "Invalid index" or import ID errors on mongodbatlas_* resources. Rewrites configs (blocks to nested attributes, references, dynamic blocks, import IDs) and verifies the upgrade with an empty plan. For Terraform users, not provider developers.
---

# Upgrade MongoDB Atlas provider configurations from 2.x to 3.0

In 3.0, several `mongodbatlas_*` resources move to a new implementation generated from the Atlas Admin API. Resource type names do not change and existing state is kept: no `moved` blocks, no `terraform state` commands and no re-import are needed. What changes is the configuration syntax of some resources, a few attribute shapes, and the import ID format.

The goal of every upgrade is the same: **after rewriting the configuration, `terraform plan` shows no changes.** A plan that updates or replaces a real resource means the rewrite is wrong or incomplete. Never apply it.

The authoritative list of affected resources and their changes is the 3.0 upgrade guide in the provider docs (`docs/guides/3.0.0-upgrade-guide.md`). Use this skill together with it: the guide says WHAT changed per resource, this skill says HOW to rewrite and verify.

## Before you start

1. **Upgrade to the latest 2.x release first** and run `terraform plan`. Fix every deprecation warning the provider prints. 2.x already accepts the new forms (for example slash-separated import IDs and replacement resources for deprecated attributes), so most 3.0 changes can be made before the major upgrade, one resource at a time.
2. Make sure the 2.x plan is empty. Do not start from a configuration that already has pending changes: you will not be able to tell upgrade diffs from real ones.
3. Back up the state (`terraform state pull > backup.tfstate`), or use a backend with versioning.
4. Find every affected resource, including the ones inside modules:
   ```bash
   grep -rn 'resource "mongodbatlas_' --include='*.tf' .
   grep -rn 'mongodbatlas_' --include='*.tf' . | grep -v '^.*resource "' # references, outputs, data sources
   ```
   Cross-check each resource type against the 3.0 upgrade guide.

## Rewrite rules

Apply only the rules the upgrade guide lists for each resource. The attribute names and values stay the same unless the guide says otherwise.

### 1. Blocks become list attributes

A repeated block becomes a list of objects: add `=` and wrap the entries in `[ ]`.

```hcl
# 2.x
notification {
  type_name     = "GROUP"
  interval_min  = 5
}
notification {
  type_name     = "EMAIL"
  email_address = "ops@example.com"
}

# 3.0
notification = [
  {
    type_name    = "GROUP"
    interval_min = 5
  },
  {
    type_name     = "EMAIL"
    email_address = "ops@example.com"
  },
]
```

References into lists keep their index (`notification[0].type_name` stays the same).

### 2. Single blocks become objects, and references lose `[0]`

A block that could appear at most once becomes a single object: add `=`, no brackets.

```hcl
# 2.x
threshold_config {
  operator  = "GREATER_THAN"
  threshold = 1
  units     = "HOURS"
}

# 3.0
threshold_config = {
  operator  = "GREATER_THAN"
  threshold = 1
  units     = "HOURS"
}
```

Every reference to it drops the index. Search outputs, locals, other resources and module outputs:

```hcl
# 2.x
value = mongodbatlas_alert_configuration.this.threshold_config[0].threshold
# 3.0
value = mongodbatlas_alert_configuration.this.threshold_config.threshold
```

If you miss one, `terraform plan` fails with `Invalid index` or `Unsupported attribute`, pointing at the file and line.

### 3. `dynamic` blocks become `for` expressions

```hcl
# 2.x
dynamic "notification" {
  for_each = var.notifications
  content {
    type_name    = notification.value.type_name
    interval_min = notification.value.interval_min
  }
}

# 3.0
notification = [
  for n in var.notifications : {
    type_name    = n.type_name
    interval_min = n.interval_min
  }
]
```

Optional fields that the `dynamic` block skipped with `lookup(..., null)` keep working: `null` in an object attribute means "not set".

### 4. Structural changes

Some resources take the API's shape instead of the old one (attributes move under a nested object, or several blocks merge into one list). These are not mechanical: the upgrade guide gives a mapping table per resource. Apply the mapping, then check the plan. Ask the user when a value's meaning is not obvious from the mapping; do not guess.

Example: in `mongodbatlas_backup_compliance_policy`, the per-frequency blocks merge into one list, and the frequency moves from the block name into a `frequency_type` value:

```hcl
# 2.x
policy_item_hourly {
  frequency_interval = 6
  retention_unit     = "days"
  retention_value    = 7
}
policy_item_daily {
  frequency_interval = 1
  retention_unit     = "days"
  retention_value    = 7
}

# 3.0
scheduled_policy_items = [
  { frequency_type = "hourly", frequency_interval = 6, retention_unit = "days", retention_value = 7 },
  { frequency_type = "daily",  frequency_interval = 1, retention_unit = "days", retention_value = 7 },
]
```

Things to check in structural changes:
- **Order** of merged list entries can matter. Use the order the upgrade guide gives (for backup compliance policies: hourly, daily, weekly, monthly, yearly), otherwise the first plan shows an update of the list. Reorder the entries instead of applying it.
- **References** into the old blocks (`policy_item_hourly[0].retention_value`) have no one-to-one replacement. Rewrite them as a filter over the new list, for example `one([for p in r.scheduled_policy_items : p.retention_value if p.frequency_type == "hourly"])`, and confirm the intent with the user.
- **Computed values** from the old blocks (such as per-item `id`) are not copied into the new configuration.

### 5. Removed attributes

Attributes removed in 3.0 were deprecated in 2.x with a replacement. Move the setting to the replacement resource named in the deprecation message (for example `mongodbatlas_project.teams` to `mongodbatlas_team_project_assignment`) before upgrading, while still on 2.x, and confirm the plan is empty.

### 6. Import IDs

Import IDs use `/` as the separator in 3.0 (for example `{project_id}/{name}` instead of `{project_id}-{name}`). This only affects future imports: `import` blocks, `terraform import` commands and scripts. Existing state is not affected. Update them to the format in each resource's docs page.

## Modules

- **Module maintainers**: rewrite the resources inside the module (rules above), release a new major version of the module, and document that it requires the 3.0 provider. Variables usually do not change; `dynamic` blocks become `for` expressions inside the module.
- **Module users**: bump the module version together with the provider. If you only consume the module, your own code changes only where it references the module's resources through outputs that changed shape.

## Verify

1. `terraform init -upgrade` with the provider constraint set to `~> 3.0`.
2. `terraform validate`: fix every syntax error it reports (each one points at a block to rewrite).
3. `terraform plan`. The expected result is **No changes**.
4. If the plan is not empty:
   - **An update or replacement of a real resource**: the rewrite is wrong (a missed attribute, a changed value, a wrong mapping). Fix the configuration. Do not apply.
   - **A one-time change listed in the upgrade guide as expected** (for example a state-only attribute): apply it once, then plan again and confirm **No changes**.
   - **Anything else**: stop and report it, with the plan output, to the provider issue tracker. Do not apply.

## Rules for the agent

- Change only `.tf` files and import scripts. Never edit state files or run commands that modify state (`terraform state mv`, `rm`, `push`, `terraform import`) for this upgrade; `terraform state pull` for a backup is fine.
- Never run `terraform apply` unless the user asks and the plan contains only changes the upgrade guide lists as expected.
- Keep values, comments and formatting. Run `terraform fmt` at the end.
- Rewrite one resource type at a time and plan after each, so a diff can be traced to its rewrite.
- When the upgrade guide and this skill disagree, the upgrade guide wins.
