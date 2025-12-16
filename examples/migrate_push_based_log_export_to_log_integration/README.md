# Migration from `mongodbatlas_push_based_log_export` to `mongodbatlas_log_integration`

This directory contains examples demonstrating how to migrate from `mongodbatlas_push_based_log_export` to `mongodbatlas_log_integration` using a **create-before-destroy** pattern. For more details, please refer to the [Migration Guide: Push-Based Log Export to Log Integration](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/push-based-log-export-to-log-integration-migration-guide).

## Overview

The `mongodbatlas_log_integration` resource provides enhanced log export capabilities compared to `mongodbatlas_push_based_log_export`:
- Faster log export intervals (1 minute vs 5 minutes)
- Configurable log types (`MONGOD`, `MONGOS`, `MONGOD_AUDIT`, `MONGOS_AUDIT`)
- Optional AWS KMS key support for server-side encryption

## Migration Approach

Unlike migrations that use the `moved` block, this migration follows a **create-before-destroy** pattern:
1. Create the new `mongodbatlas_log_integration` resource alongside the existing one
2. Validate that logs are appearing at the new configuration
3. Remove the old `mongodbatlas_push_based_log_export` resource

This approach ensures continuous log delivery during migration without gaps.

## Examples

The examples are organized as follows:

- **For users directly utilizing the `mongodbatlas_push_based_log_export` resource**: please check the [basic/](./basic/README.md) folder. This shows a step-by-step migration through three phases:
  - `v1/` - Original configuration with `mongodbatlas_push_based_log_export`
  - `v2/` - Both resources running during migration
  - `v3/` - Final configuration with only `mongodbatlas_log_integration`

- **For users employing modules to manage log exports**: please see the [module_maintainer/](./module_maintainer/README.md) and [module_user/](./module_user/README.md) folders. These folders illustrate the migration process from both the maintainer's and the user's perspectives, using a feature flag approach.

## Important Considerations

- **Log Duplication**: During the overlap period when both resources are active, there may be duplicate logs. Use distinct prefix paths to manage this.
- **Validation**: Always verify logs are appearing at the new configuration before removing the old resource.
- **Rollback**: If issues occur, you can simply remove the new resource and keep using the old one until resolved.

