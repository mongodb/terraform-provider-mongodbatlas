# Module User Migration Example

This directory contains examples for **module users** showing how to migrate when using a module that manages log exports.

## Overview

As a module user, the migration process involves two steps after your module maintainer releases a version supporting the migration:

1. **Upgrade to the new module version** with `skip_push_based_log_export = false`
2. **Set `skip_push_based_log_export = true`** after validating the new configuration

## Examples

- **v1/** - Initial upgrade: Both resources active during validation
- **v2/** - Final state: Old resource removed after successful validation

## Migration Steps

### Step 1: Upgrade Module (v1)

Update your module source to the new version and set `skip_push_based_log_export = false`:

```bash
cd v1
terraform init -upgrade
terraform apply
```

At this point:
- ✅ `mongodbatlas_log_integration` is created with the new prefix path
- ✅ `mongodbatlas_push_based_log_export` still exists
- ⚠️ Both resources are exporting logs (expect some duplication)

**Validate** that logs appear at the new prefix path before proceeding.

### Step 2: Remove Old Resource (v2)

After validation, set `skip_push_based_log_export = true`:

```bash
cd ../v2
terraform apply
```

At this point:
- ✅ `mongodbatlas_log_integration` continues operating
- ✅ `mongodbatlas_push_based_log_export` is destroyed
- ✅ Migration complete!

