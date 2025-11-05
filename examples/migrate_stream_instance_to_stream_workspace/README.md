# Migration from `mongodbatlas_stream_instance` to `mongodbatlas_stream_workspace`

This example demonstrates how to migrate a `mongodbatlas_stream_instance` resource to `mongodbatlas_stream_workspace` using the `moved` block. For more details, please refer to the [Migration Guide: Stream Instance to Stream Workspace](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/stream-instance-to-stream-workspace-migration-guide).

## Overview

The `mongodbatlas_stream_workspace` resource is the new preferred way to manage MongoDB Atlas Stream Processing instances. It provides the same functionality as `mongodbatlas_stream_instance` but uses updated terminology that aligns with the MongoDB Atlas UI.

**Key Changes:**
- Resource name: `mongodbatlas_stream_instance` → `mongodbatlas_stream_workspace`
- Field name: `instance_name` → `workspace_name`
- All other fields remain the same

## Migration Steps

This example shows the complete migration process:

1. [Create the `mongodbatlas_stream_instance`](#step-1-create-the-stream-instance) (skip if you already have one)
2. [Add the `mongodbatlas_stream_workspace` configuration](#step-2-add-stream-workspace-configuration)
3. [Add the moved block](#step-3-add-moved-block)
4. [Perform the migration](#step-4-perform-the-migration)
5. [Clean up](#step-5-clean-up)

## Step 1: Create the Stream Instance

**Note**: Skip this step if you already have a `mongodbatlas_stream_instance` resource.

1. Uncomment the code in [stream_instance.tf](stream_instance.tf)
2. Comment the code in [stream_workspace.tf](stream_workspace.tf)
3. Create a `terraform.tfvars` file:
```terraform
project_id = "your-project-id"
workspace_name = "my-stream-workspace"
```
4. Run `terraform init && terraform apply`

## Step 2: Add Stream Workspace Configuration

1. Comment out the `mongodbatlas_stream_instance` in [stream_instance.tf](stream_instance.tf)
2. Uncomment the `mongodbatlas_stream_workspace` in [stream_workspace.tf](stream_workspace.tf)
3. Update any references from `mongodbatlas_stream_instance.example` to `mongodbatlas_stream_workspace.example`

## Step 3: Add Moved Block

The moved block in [stream_workspace.tf](stream_workspace.tf) tells Terraform to migrate the state:

```terraform
moved {
  from = mongodbatlas_stream_instance.example
  to   = mongodbatlas_stream_workspace.example
}
```

## Step 4: Perform the Migration

1. Run `terraform validate` to ensure there are no configuration errors
2. Run `terraform plan` - you should see:
   ```
   Terraform will perform the following actions:
     # mongodbatlas_stream_instance.example has moved to mongodbatlas_stream_workspace.example
         resource "mongodbatlas_stream_workspace" "example" {
             workspace_name = "my-stream-workspace"
             # (other unchanged attributes hidden)
         }

   Plan: 0 to add, 0 to change, 0 to destroy.
   ```
3. Run `terraform apply` and type `yes` to confirm the migration

## Step 5: Clean Up

After successful migration:
1. Remove the commented `mongodbatlas_stream_instance` resource from [stream_instance.tf](stream_instance.tf)
2. The `moved` block can be kept for historical reference or removed after the migration is complete

## Troubleshooting

- **Reference errors**: Ensure all references are updated from `mongodbatlas_stream_instance.example` to `mongodbatlas_stream_workspace.example`
- **Field name errors**: Make sure to use `workspace_name` instead of `instance_name` in the new resource
- **Plan changes**: If you see unexpected changes, verify that all field values match between the old and new resources

## Next Steps

After migration, you can:
- Update any related resources (stream connections, processors) to reference the new workspace
- Use the new `mongodbatlas_stream_workspace` data source for lookups
- Refer to the [stream workspace documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/stream_workspace) for additional configuration options
