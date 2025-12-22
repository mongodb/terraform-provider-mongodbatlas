# v2: Updated Log Export Module with Migration Support

This module version supports migration from `mongodbatlas_push_based_log_export` to `mongodbatlas_log_integration` using a feature flag pattern.

## Usage

### Initial Upgrade (both resources active)

```hcl
module "log_export" {
  source = "./module_maintainer/v2"

  project_id      = var.project_id
  bucket_name     = var.bucket_name
  iam_role_id     = var.iam_role_id
  prefix_path     = "atlas-logs"
  new_prefix_path = "atlas-logs-new"
  log_types       = ["MONGOD", "MONGOD_AUDIT"]

  skip_push_based_log_export = false  # Keep old resource during transition
}
```

### After Validation (remove old resource)

```hcl
module "log_export" {
  source = "./module_maintainer/v2"

  project_id      = var.project_id
  bucket_name     = var.bucket_name
  iam_role_id     = var.iam_role_id
  prefix_path     = "atlas-logs"
  new_prefix_path = "atlas-logs"  # Can use same path now
  log_types       = ["MONGOD", "MONGOD_AUDIT"]

  skip_push_based_log_export = true  # Remove old resource
}
```

## Variables

| Name | Description | Required |
|------|-------------|----------|
| `skip_push_based_log_export` | Set to `true` to skip the old resource | No (default: `false`) |
| `new_prefix_path` | Prefix path for the new log integration | Yes |
| `log_types` | Array of log types to export | No (default: all types) |

