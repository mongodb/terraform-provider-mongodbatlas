---
page_title: "Migration Guide: Project API Key to API Key + Assignment"
---

# Migration Guide: Project API Key to API Key + Assignment

**Objective**: This guide explains how to migrate from the legacy `mongodbatlas_project_api_key` resource/data source to the more flexible and future-proof pattern of managing API keys and project assignments separately using `mongodbatlas_api_key` and `mongodbatlas_api_key_project_assignment`.

## Why do we have both `mongodbatlas_project_api_key` and the new pattern?

Historically, the `mongodbatlas_project_api_key` resource allowed users to create and assign API keys to projects in a single step. However, this approach limited flexibility, especially for organizations managing many projects or requiring more granular control. The new pattern—creating API keys independently and assigning them to projects with `mongodbatlas_api_key_project_assignment`—offers greater flexibility, clarity, and aligns with best practices for infrastructure as code.

## Why should I migrate?
- **Flexibility:** Manage API keys and assignments independently.
- **Clarity:** Clear separation of key creation and project assignment.
- **Best Practices:** Aligns with Terraform and Atlas recommendations for resource management.

## Will `mongodbatlas_project_api_key` continue to work?

While this is not our recommended approach, you can still continue to use the `mongodbatlas_project_api_key` resource. If you are creating a new configuration, use the `mongodbatlas_api_key_project_assignment` resource.

## Main Changes Between Patterns

| Old Pattern (`mongodbatlas_project_api_key`) | New Pattern (`mongodbatlas_api_key` + `mongodbatlas_api_key_project_assignment`) |
|----------------------------------------------|--------------------------------------------------------------------------|
| API key creation and assignment are coupled  | API key creation and assignment are decoupled                            |
| Assignments are defined within the resource  | Assignments are managed by separate resources                            |
| Less flexible for multi-project assignments  | Easily assign the same key to multiple projects                          |

## Before and After: Example Configurations

### Old Pattern (Legacy)
```hcl
resource "mongodbatlas_project_api_key" "old" {
  description = "example key"
  project_assignment {
    project_id = var.project_id
    role_names = ["GROUP_READ_ONLY"]
  }
}
```

### New Pattern (Recommended)
```hcl
resource "mongodbatlas_api_key" "new" {
  org_id      = var.org_id
  description = "example key"
  role_names  = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_api_key_project_assignment" "new" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.new.api_key_id
  role_names = ["GROUP_READ_ONLY"]
}
```

## Migration using import

If you are migrating from `mongodbatlas_project_api_key` resources already managed in Terraform, **import is required** because each `mongodbatlas_project_api_key` resource is equivalent to one `mongodbatlas_api_key` and one `mongodbatlas_api_key_project_assignment`. This ensures that your existing Atlas API keys and assignments are not deleted and recreated during the migration.

### Best Practices Before Importing
- **Backup your Terraform state file** before making any changes.
- **Review your current usage** of `mongodbatlas_project_api_key` and plan the migration for all affected resources.
- **Test the migration in a non-production environment** if possible.

### Step-by-Step Migration Guide

1. **Add the new resources to your configuration:**
   ```hcl
   resource "mongodbatlas_api_key" "new" {
     org_id      = var.org_id
     description = "example key"
     role_names  = ["ORG_READ_ONLY"]
   }

   resource "mongodbatlas_api_key_project_assignment" "new" {
     project_id = var.project_id
     api_key_id = mongodbatlas_api_key.new.api_key_id
     role_names = ["GROUP_READ_ONLY"]
   }
   ```
2. **Import the existing API key into the API key resource:**
   ```shell
   terraform import mongodbatlas_api_key.new <ORG_ID>-<API_KEY_ID>
   ```
3. **Import the existing project assignment into the assignment resource:**
   ```shell
   terraform import mongodbatlas_api_key_project_assignment.new <PROJECT_ID>/<API_KEY_ID>
   ```
4. **Remove the old resource from the Terraform state:**
   ```shell
   terraform state rm mongodbatlas_project_api_key.old
   ```
   This step ensures Terraform will not attempt to delete the underlying Atlas resource. Alternatively a `removed block` (available in Terraform 1.7 and later) can be used to delete it from the state file, e.g.:
   ```terraform
   removed {
    from = mongodbatlas_project_api_key.old
     lifecycle {
      destroy = false
     }
   }
  ```
5. **Remove the old resource block from your configuration**
6. **Run `terraform plan` to review the changes.**
   - Ensure that Terraform does not plan to delete or recreate your API keys or assignments.
   - **Note:** After import, Terraform may show an in-place update for attributes like `org_id` on the `mongodbatlas_api_key` resource. This is expected and only updates the Terraform state to match your configuration; it does not change the actual resource in Atlas. You can safely apply this change.
6. **Run `terraform apply` to apply the migration.**
   - Your resources should now be managed under the new resource types without any disruption.

This process ensures that your existing Atlas API keys and assignments are preserved and managed by Terraform under the new resource types, with no deletion or recreation.

## Migration using Modules

If you are using modules to manage your API key assignments, migrating from `mongodbatlas_project_api_key` to the new pattern requires special attention. Because the old resource corresponds to two new resources (`mongodbatlas_api_key` and `mongodbatlas_api_key_project_assignment`), you cannot simply move the resource block inside your module and expect Terraform to handle the migration automatically. This section demonstrates how to migrate from a module using the legacy `mongodbatlas_project_api_key` resource to a module using the new `mongodbatlas_api_key` and `mongodbatlas_api_key_project_assignment` resources.

**Key points for module users:**
- You must use `terraform import` to bring existing API keys and assignments into the new resources, even when they are managed inside a module.
- The import command must match the resource address as used in your module (e.g., `module.<module_name>.mongodbatlas_api_key.<name>`).
- After import, remove the old resource from your configuration and state as described below.

**Example import commands for modules:**
```shell
terraform import 'module.<module_name>.mongodbatlas_api_key.<name>' <ORG_ID>-<API_KEY_ID>
terraform import 'module.<module_name>.mongodbatlas_api_key_project_assignment.<name>' <PROJECT_ID>/<API_KEY_ID>
```

### 1. Old Module Usage (Legacy)

```hcl
module "project_api_key" {
  source     = "./old_module"
  project_id = var.project_id
  role_names = var.role_names
}
```

### 2. New Module Usage (Recommended)

```hcl
module "api_key_assignment" {
  source     = "./new_module"
  org_id     = var.org_id
  project_id = var.project_id
  role_names = var.role_names
}
```

### 3. Migration Steps

1. **Add the new module to your configuration:**
   - Add the new module block as shown above, using the same input variables as appropriate.
2. **Import the existing API key and assignment into the new resources:**
   - Use the correct resource addresses for your module:
   ```shell
   terraform import 'module.api_key_assignment.mongodbatlas_api_key.this' <ORG_ID>-<API_KEY_ID>
   terraform import 'module.api_key_assignment.mongodbatlas_api_key_project_assignment.this' <PROJECT_ID>/<API_KEY_ID>
   ```
3. **Remove the old resource from the Terraform state:**
   ```shell
   terraform state rm 'module.project_api_key.mongodbatlas_project_api_key.this'
   ```
   Alternatively a `removed block` (available in Terraform 1.7 and later) can be used to delete it from the state file.
4. **Remove the old module block from your configuration.**
5. **Run `terraform plan` to review the changes.**
   - Ensure that Terraform does not plan to delete or recreate your API keys or assignments.
6. **Run `terraform apply` to apply the migration.**

For complete working examples, see:
- [Old module example](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_api_key_assignment/module/old_module/)
- [New module example](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_api_key_assignment/module/new_module/)


## FAQ
**Q: Can I assign the same API key to multiple projects?**
A: Yes, simply create multiple `mongodbatlas_api_key_project_assignment` resources for each project.

**Q: Do I need to change anything for existing keys?**
A: Existing keys will continue to work, but we recommend following the migration guide to move to the new pattern.

**Q: What if I have many project assignments?**
A: You can use [`for_each`](https://developer.hashicorp.com/terraform/language/meta-arguments/for_each) or [`count`](http://developer.hashicorp.com/terraform/language/meta-arguments/count) with `mongodbatlas_api_key_project_assignment` to manage multiple assignments efficiently.

**Q: Where can I find a working example?**
A: See [examples/mongodbatlas_api_key/main.tf](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_api_key/main.tf).

## Further Resources
- [API Key Project Assignment Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/api_key_project_assignment)
