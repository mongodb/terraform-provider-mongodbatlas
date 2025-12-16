# Basic Migration from `mongodbatlas_push_based_log_export` to `mongodbatlas_log_integration`

This example demonstrates how to migrate a `mongodbatlas_push_based_log_export` resource to `mongodbatlas_log_integration` using a create-before-destroy pattern (see alternatives and more details in the [push-based log export to log integration migration guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/push-based-log-export-to-log-integration-migration-guide)).

The migration is shown in three phases:

- **v1/** - Original configuration with `mongodbatlas_push_based_log_export`
- **v2/** - Both resources running during migration (with distinct prefix paths)
- **v3/** - Final configuration with only `mongodbatlas_log_integration`

## Prerequisites

- An AWS account with permissions to create S3 buckets and IAM roles
- A MongoDB Atlas organization and project
- Terraform 1.0 or later
- MongoDB Atlas Provider (latest version recommended)

## Migration Steps

### Step 1: Start with v1 Configuration

If you already have a `mongodbatlas_push_based_log_export` resource, your configuration should look similar to `v1/`.

```bash
cd v1
terraform init
terraform apply
```

### Step 2: Add the New Resource (v2)

1. Copy your configuration to match `v2/`
2. Add the `mongodbatlas_log_integration` resource with a **distinct prefix path**
3. Apply the changes:

```bash
cd ../v2
terraform init
terraform apply
```

4. **Validate** that logs are appearing at the new prefix path in your S3 bucket

### Step 3: Remove the Old Resource (v3)

After validating the new configuration:

1. Remove the `mongodbatlas_push_based_log_export` resource from your configuration
2. Optionally update the prefix path if you want to use the original path
3. Apply the changes:

```bash
cd ../v3
terraform init
terraform apply
```

## Important Notes

- **Prefix Paths**: During migration (v2), use distinct prefix paths to avoid log conflicts
- **Log Duplication**: Expect some log duplication during the overlap period
- **Validation Time**: Allow sufficient time to validate the new configuration before removing the old resource
- **Rollback**: If issues occur during migration, you can remove the `mongodbatlas_log_integration` resource and continue using the original configuration

