---
page_title: "Migration Guide: Push-Based Log Export to Log Integration"
---

# Migration Guide: Push-Based Log Export to Log Integration

**Objective**: This guide explains how to migrate from the `mongodbatlas_push_based_log_export` resource to the `mongodbatlas_log_integration` resource using a **create-before-destroy** pattern to ensure uninterrupted log delivery during the transition.

-> **NOTE:** For comprehensive information about this feature migration from MongoDB Atlas, refer to the official [MongoDB Atlas Log Integration documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-push-based-log-export).

## Why migrate to `mongodbatlas_log_integration`?

The `mongodbatlas_log_integration` resource provides enhanced log export capabilities compared to `mongodbatlas_push_based_log_export`:

1. **Faster Log Export**: Logs are exported every **1 minute** instead of every 5 minutes
2. **Log Type Selection**: Choose specific log types to export (`MONGOD`, `MONGOS`, `MONGOD_AUDIT`, `MONGOS_AUDIT`)
3. **Enhanced Encryption**: Optional AWS KMS key support for server-side encryption
4. **Future-Proof**: New log integration features will be developed using this resource

## Main Differences Between Resources

| Feature | `mongodbatlas_push_based_log_export` | `mongodbatlas_log_integration` |
|---------|--------------------------------------|--------------------------------|
| Export interval | 5 minutes | 1 minute |
| Log type selection | All logs | Configurable via `log_types` |
| KMS encryption | Not available | Optional via `kms_key` |
| Integration type | Implicit S3 | Explicit via `type` |
| Resource identifier | `project_id` only | `project_id` + `integration_id` |

## Migration Approach

This migration uses a **create-before-destroy** pattern rather than a `moved` block. This approach:
- Ensures continuous log delivery during migration
- Allows validation of the new configuration before removing the old one

~> **Important:** Some log duplication may occur during the overlap period. To minimize this, configure distinct prefix paths for the old and new configurations.

## Migration Steps

For complete working examples, see the [basic migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_push_based_log_export_to_log_integration/basic).

### Step 1: Create the New Resource

Add the `mongodbatlas_log_integration` resource to your configuration alongside the existing `mongodbatlas_push_based_log_export` resource:

```terraform
# Existing push-based log export (keep this for now)
resource "mongodbatlas_push_based_log_export" "old" {
  project_id  = var.project_id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "old-logs"
}

# New log integration resource with a distinct prefix path
resource "mongodbatlas_log_integration" "new" {
  project_id  = var.project_id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "new-logs"  # Use a different prefix to avoid conflicts
  type        = "S3_LOG_EXPORT"
  log_types   = ["MONGOD", "MONGOD_AUDIT"]  # Configure as needed
}
```

Run `terraform apply` to create the new log integration.

### Step 2: Test the New Configuration

Validate that the new `mongodbatlas_log_integration` resource is functioning correctly:
- Check for logs appearing in the S3 bucket at the new prefix path
- Verify the expected log types are being exported

### Step 3: Remove the Old Resource

After confirming the new resource is operating successfully, remove the `mongodbatlas_push_based_log_export` resource from your configuration:

```terraform
# Remove or comment out the old resource
# resource "mongodbatlas_push_based_log_export" "old" {
#   ...
# }

# Keep only the new log integration resource
resource "mongodbatlas_log_integration" "new" {
  project_id  = var.project_id
  bucket_name = aws_s3_bucket.log_bucket.bucket
  iam_role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  prefix_path = "new-logs"
  type        = "S3_LOG_EXPORT"
  log_types   = ["MONGOD", "MONGOD_AUDIT"]
}
```

Run `terraform apply` to remove the old resource.

### Step 4: Verify Final Configuration

Confirm the migration is complete:
- Verify logs continue to appear at the new prefix path
- Ensure the old log export is no longer active
- Optionally, consolidate prefix paths after the old resource is fully removed

## Migration for Modules

If you are using modules to manage log exports, the migration follows a similar create-before-destroy approach with a feature flag pattern.

For complete working examples, see:
- [Module maintainer example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_push_based_log_export_to_log_integration/module_maintainer)
- [Module user example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_push_based_log_export_to_log_integration/module_user)

### Module Maintainer Steps

1. **Create a new module version** that uses `mongodbatlas_log_integration` and includes a flag (e.g., `skip_push_based_log_export`) to control the lifecycle of the old resource:

```terraform
# modules/log_export/main.tf

variable "skip_push_based_log_export" {
  description = "Set to true to skip creating the push_based_log_export resource"
  type        = bool
  default     = false
}

# Old resource - conditionally created
resource "mongodbatlas_push_based_log_export" "this" {
  count       = var.skip_push_based_log_export ? 0 : 1
  project_id  = var.project_id
  bucket_name = var.bucket_name
  iam_role_id = var.iam_role_id
  prefix_path = var.prefix_path
}

# New resource - always created in the new module version
resource "mongodbatlas_log_integration" "this" {
  project_id  = var.project_id
  bucket_name = var.bucket_name
  iam_role_id = var.iam_role_id
  prefix_path = var.new_prefix_path  # Use distinct path during migration
  type        = "S3_LOG_EXPORT"
  log_types   = var.log_types
}
```

### Module User Steps

1. **Upgrade to the new module version** with `skip_push_based_log_export = false`:
   - The `mongodbatlas_log_integration` resource is created
   - The `mongodbatlas_push_based_log_export` resource still exists

```terraform
module "log_export" {
  source  = "./modules/log_export"
  version = "2.0.0"  # New version with log_integration support

  project_id     = var.project_id
  bucket_name    = var.bucket_name
  iam_role_id    = var.iam_role_id
  prefix_path    = "old-logs"
  new_prefix_path = "new-logs"
  log_types      = ["MONGOD", "MONGOD_AUDIT"]

  skip_push_based_log_export = false  # Keep old resource during transition
}
```

2. **Set `skip_push_based_log_export = true`** after validating the new configuration:
   - The `mongodbatlas_log_integration` resource continues to exist
   - The `mongodbatlas_push_based_log_export` resource is destroyed

```terraform
module "log_export" {
  source  = "./modules/log_export"
  version = "2.0.0"

  project_id     = var.project_id
  bucket_name    = var.bucket_name
  iam_role_id    = var.iam_role_id
  prefix_path    = "old-logs"
  new_prefix_path = "new-logs"
  log_types      = ["MONGOD", "MONGOD_AUDIT"]

  skip_push_based_log_export = true  # Remove old resource
}
```

## Customer Implications

### Log Duplications

Atlas guarantees **at least once delivery** of logs. During the overlap period when both resources are active, there will be duplicated logs in the destination, particularly if the prefix paths are the same. To manage this:

- **Use distinct paths** for the old and new configurations until the old resource is fully removed

### Delay for Removal

The delay between creating the new resource and destroying the old one ensures:
- Uninterrupted log flow during migration
- Time to verify the new resource's functionality

## Recommended Best Practices

1. **Use distinct prefix paths** when configuring the new `mongodbatlas_log_integration` resource during migration
2. **Perform rigorous validation** and testing of the new resource before destroying the old resource
3. **Monitor log delivery** in your S3 bucket throughout the migration process
4. **Back up your Terraform state file** before starting the migration
5. **Plan for log duplication** during the overlap period and consider how to handle duplicate logs in downstream processing

## Further Resources

- [Migration Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_push_based_log_export_to_log_integration)
- [Log Integration Resource Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/log_integration)
- [Log Integration Data Source Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/log_integration)
- [Push-Based Log Export Resource Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/push_based_log_export)
- [MongoDB Atlas Log Integration API Documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-push-based-log-export)

