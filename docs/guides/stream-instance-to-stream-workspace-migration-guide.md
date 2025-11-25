---
page_title: "Migration Guide: Stream Instance to Stream Workspace"
---

# Migration Guide: Stream Instance to Stream Workspace

**Objective**: This guide explains how to replace the deprecated `mongodbatlas_stream_instance` resource with the `mongodbatlas_stream_workspace` resource. For data source migrations, refer to the [output changes](#output-changes) section.

## Why do we have both `mongodbatlas_stream_instance` and `mongodbatlas_stream_workspace` resources?

Both `mongodbatlas_stream_instance` and `mongodbatlas_stream_workspace` resources currently allow customers to manage MongoDB Atlas Stream Processing workspaces. Initially, only `mongodbatlas_stream_instance` existed. However, MongoDB Atlas has evolved its terminology and API to use "workspace" instead of "instance" for stream processing environments. To align with this change and provide a clearer, more consistent naming convention, we created the `mongodbatlas_stream_workspace` resource as a direct replacement.

## If I am using `mongodbatlas_stream_instance`, why should I move to `mongodbatlas_stream_workspace`?

The `mongodbatlas_stream_workspace` resource provides the exact same functionality as `mongodbatlas_stream_instance` but with updated terminology that aligns with MongoDB Atlas's current naming conventions. This change provides:

1. **Consistent Terminology**: Aligns with MongoDB Atlas's current documentation and UI terminology
2. **Future-Proof**: New stream processing features will be developed using the workspace terminology
3. **Clearer Intent**: The term "workspace" better describes the stream processing environment

To maintain consistency with MongoDB Atlas's terminology and ensure you're using the most current resource names, we recommend migrating to `mongodbatlas_stream_workspace`. The `mongodbatlas_stream_instance` resource is deprecated and will be removed in future major versions of the provider.

## How should I move to `mongodbatlas_stream_workspace`?

To move from `mongodbatlas_stream_instance` to `mongodbatlas_stream_workspace` we offer two alternatives:
1. [(Recommended) Use the `moved` block](#migration-using-the-moved-block-recommended)
2. [Manually use the import command with the `mongodbatlas_stream_workspace` resource](#migration-using-import)

### Best Practices Before Migrating

Before doing any migration, create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).

## Migration using the Moved block (recommended)

This is our recommended method to migrate from `mongodbatlas_stream_instance` to `mongodbatlas_stream_workspace`. The [moved block](https://developer.hashicorp.com/terraform/language/moved) is a Terraform feature that allows to move between resource types. It's conceptually similar to running `removed` and `import` commands separately but it brings the convenience of doing it in one step.

**Prerequisites:**
 - Terraform version 1.8 or later is required, more information in the [State Move page](https://developer.hashicorp.com/terraform/plugin/framework/resources/state-move).
 - MongoDB Atlas Provider version 2.1 or later is required.

The basic experience when using the `moved` block is as follows:
1. Before starting, run `terraform plan` to make sure that there are no planned changes.
2. Add the `mongodbatlas_stream_workspace` resource definition.
3. Comment out or delete the `mongodbatlas_stream_instance` resource definition.
4. Update the references from your previous stream instance resource: `mongodbatlas_stream_instance.this.XXXX` to the new `mongodbatlas_stream_workspace.this.XXX`.
   - Change `instance_name` to `workspace_name` in your references
   - Double check [output-changes](#output-changes) to ensure the underlying configuration stays unchanged.
5. Add the `moved` block to your configuration file, e.g.:
```terraform
moved {
  from = mongodbatlas_stream_instance.this
  to   = mongodbatlas_stream_workspace.this
}
```
6. Run `terraform plan` and make sure that there are no planned changes, only the moved block should be shown. This is an example output of a successful plan:
```text
 # mongodbatlas_stream_instance.this has moved to mongodbatlas_stream_workspace.this
     resource "mongodbatlas_stream_workspace" "this" {
         workspace_name           = "my-workspace"
         # (6 unchanged attributes hidden)
     }

 Plan: 0 to add, 0 to change, 0 to destroy.
```

7. Run `terraform apply` to apply the changes. The `mongodbatlas_stream_instance` resource will be removed from the Terraform state and the `mongodbatlas_stream_workspace` resource will be added.
8. Hashicorp recommends to keep the move block in your configuration file to help track the migrations, however you can delete the `moved` block from your configuration file without any adverse impact.

## Migration using import

**Note**: We recommend the [`moved` block](#migration-using-the-moved-block-recommended) method as it's more convenient and less error-prone.

This method uses [Terraform native tools](https://developer.hashicorp.com/terraform/language/import/generating-configuration) and works if you:
1. Have an existing stream workspace without any Terraform configuration and want to import and manage it with Terraform.
2. Have existing `mongodbatlas_stream_instance` resource(s) but you can't use the [recommended approach](#migration-using-the-moved-block-recommended).

The process works as follow:
1. If you have an existing `mongodbatlas_stream_instance` resource, remove it from your configuration and delete it from the state file, e.g.: `terraform state rm mongodbatlas_stream_instance.this`.
2. Find the import IDs of the stream workspaces you want to migrate: `{PROJECT_ID}-{WORKSPACE_NAME}`, such as `664619d870c247237f4b86a6-my-workspace`
3. Import it using the `terraform import` command, e.g.: `terraform import mongodbatlas_stream_workspace.this 664619d870c247237f4b86a6-my-workspace`.
4. Run `terraform plan -generate-config-out=stream_workspace.tf`. This should generate a `stream_workspace.tf` file.
5. Update the references from your previous stream instance resource: `mongodbatlas_stream_instance.this.XXXX` to the new `mongodbatlas_stream_workspace.this.XXX`.
   - Change `instance_name` to `workspace_name` in your references
   - Double check [output-changes](#output-changes) to ensure the underlying configuration stays unchanged.
6. Run `terraform apply`. You should see the resource(s) imported.

## Main Changes Between `mongodbatlas_stream_instance` and `mongodbatlas_stream_workspace`

The primary change is the field name for identifying the stream processing workspace:

1. **Field Name**: `instance_name` is replaced with `workspace_name`
2. **Resource Name**: `mongodbatlas_stream_instance` becomes `mongodbatlas_stream_workspace`
3. **Data Source Names**: 
   - `mongodbatlas_stream_instance` becomes `mongodbatlas_stream_workspace`
   - `mongodbatlas_stream_instances` becomes `mongodbatlas_stream_workspaces`

All other functionality remains identical.

### Example 1: Old Configuration (`mongodbatlas_stream_instance`)

```terraform
resource "mongodbatlas_stream_instance" "this" {
  project_id    = var.project_id
  instance_name = "my-stream-workspace"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

data "mongodbatlas_stream_instance" "this" {
  project_id    = var.project_id
  instance_name = mongodbatlas_stream_instance.this.instance_name
}

data "mongodbatlas_stream_instances" "all" {
  project_id = var.project_id
}
```

### Example 2: New Configuration (`mongodbatlas_stream_workspace`)

```terraform
resource "mongodbatlas_stream_workspace" "this" {
  project_id     = var.project_id
  workspace_name = "my-stream-workspace"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

data "mongodbatlas_stream_workspace" "this" {
  project_id     = var.project_id
  workspace_name = mongodbatlas_stream_workspace.this.workspace_name
}

data "mongodbatlas_stream_workspaces" "all" {
  project_id = var.project_id
}
```

### Output Changes

The only change in outputs is the field name:
- **Field Name Change**: 
  - Before: `mongodbatlas_stream_instance.this.instance_name`
  - After: `mongodbatlas_stream_workspace.this.workspace_name`

All other attributes (`project_id`, `data_process_region`, `stream_config`, `hostnames`, etc.) remain exactly the same.

## Complete Migration Example with Moved Block

Here's a complete example showing the migration process:

### Step 1: Original Configuration
```terraform
resource "mongodbatlas_stream_instance" "example" {
  project_id    = var.project_id
  instance_name = "my-workspace"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

output "stream_hostnames" {
  value = mongodbatlas_stream_instance.example.hostnames
}
```

### Step 2: Add Moved Block and New Resource
```terraform
# Add the moved block
moved {
  from = mongodbatlas_stream_instance.example
  to   = mongodbatlas_stream_workspace.example
}

# Replace with new resource (note: instance_name becomes workspace_name)
resource "mongodbatlas_stream_workspace" "example" {
  project_id     = var.project_id
  workspace_name = "my-workspace"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

# Update output references
output "stream_hostnames" {
  value = mongodbatlas_stream_workspace.example.hostnames
}
```

### Step 3: Apply and Clean Up
After running `terraform apply`, you can optionally remove the `moved` block:

```terraform
resource "mongodbatlas_stream_workspace" "example" {
  project_id     = var.project_id
  workspace_name = "my-workspace"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

output "stream_hostnames" {
  value = mongodbatlas_stream_workspace.example.hostnames
}
```
