# Module Maintainer Migration Example

This directory contains examples for **module maintainers** showing how to update a log export module to support migration from `mongodbatlas_push_based_log_export` to `mongodbatlas_log_integration`.

## Overview

As a module maintainer, you need to provide a migration path for your users. This example demonstrates the **feature flag pattern** that allows users to migrate at their own pace.

## Module Versions

- **v1/** - Original module using only `mongodbatlas_push_based_log_export`
- **v2/** - Updated module with:
  - New `mongodbatlas_log_integration` resource (always created)
  - Feature flag `skip_push_based_log_export` to control the old resource lifecycle
  - Required `prefix_path` variable (no default) forcing users to specify a distinct S3 path for the new resource, preventing log conflicts during migration

## Migration Workflow

1. **Release v2 of your module** with the feature flag defaulting to `false`
2. **Users upgrade to v2** with `skip_push_based_log_export = false`:
   - Both resources are created
   - Users validate the new configuration
3. **Users set `skip_push_based_log_export = true`**:
   - Old resource is destroyed
   - New resource continues operating
4. **Future major version**: Remove the old resource and feature flag entirely

## Key Design Decisions

- **Default to `false`**: Ensures backward compatibility when users upgrade
- **Distinct prefix paths**: Prevents log conflicts during the overlap period
- **Clear variable descriptions**: Help users understand the migration process

