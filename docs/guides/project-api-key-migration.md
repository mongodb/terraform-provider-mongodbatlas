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

## Main Changes Between Patterns

| Old Pattern (`mongodbatlas_project_api_key`) | New Pattern (`mongodbatlas_api_key` + `mongodbatlas_api_key_project_assignment`) |
|----------------------------------------------|--------------------------------------------------------------------------|
| API key creation and assignment are coupled  | API key creation and assignment are decoupled                            |
| Assignments are defined within the resource  | Assignments are managed by separate resources                            |
| Less flexible for multi-project assignments  | Easily assign the same key to multiple projects                          |

## Before and After: Example Configurations

### Old Pattern (Legacy)
```hcl
resource "mongodbatlas_project_api_key" "example" {
  description = "example key"
  project_assignment {
    project_id = var.project_id
    role_names = ["GROUP_READ_ONLY"]
  }
}
```

### New Pattern (Recommended)
```hcl
resource "mongodbatlas_api_key" "example" {
  org_id      = var.org_id
  description = "example key"
  role_names  = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_api_key_project_assignment" "example" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
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
   resource "mongodbatlas_api_key" "example" {
     org_id      = var.org_id
     description = "example key"
     role_names  = ["ORG_READ_ONLY"]
   }

   resource "mongodbatlas_api_key_project_assignment" "example" {
     project_id = var.project_id
     api_key_id = mongodbatlas_api_key.example.api_key_id
     role_names = ["GROUP_READ_ONLY"]
   }
   ```
2. **Import the existing API key into the API key resource:**
   ```shell
   terraform import mongodbatlas_api_key.example <ORG_ID>-<API_KEY_ID>
   ```
3. **Import the existing project assignment into the assignment resource:**
   ```shell
   terraform import mongodbatlas_api_key_project_assignment.example <PROJECT_ID>/<API_KEY_ID>
   ```
4. **Remove the old resource from the Terraform state:**
   ```shell
   terraform state rm mongodbatlas_project_api_key.example
   ```
   This step ensures Terraform will not attempt to delete the underlying Atlas resource.
5. **Remove the old resource block from your configuration**
6. **Run `terraform plan` to review the changes.**
   - Ensure that Terraform does not plan to delete or recreate your API keys or assignments.
   - **Note:** After import, Terraform may show an in-place update for attributes like `org_id` on the `mongodbatlas_api_key` resource. This is expected and only updates the Terraform state to match your configuration; it does not change the actual resource in Atlas. You can safely apply this change.
6. **Run `terraform apply` to apply the migration.**
   - Your resources should now be managed under the new resource types without any disruption.

### Example
Suppose you previously had:
```hcl
resource "mongodbatlas_project_api_key" "example" {
  description = "example key"
  project_assignment {
    project_id = var.project_id
    role_names = ["GROUP_READ_ONLY"]
  }
}
```
You would:
1. Add:
   ```hcl
   resource "mongodbatlas_api_key" "example" {
     org_id      = var.org_id
     description = "example key"
     role_names  = ["ORG_READ_ONLY"]
   }

   resource "mongodbatlas_api_key_project_assignment" "example" {
     project_id = var.project_id
     api_key_id = mongodbatlas_api_key.example.api_key_id
     role_names = ["GROUP_READ_ONLY"]
   }
   ```
2. Import the existing resources:
   ```shell
   terraform import mongodbatlas_api_key.example <ORG_ID>-<API_KEY_ID>
   terraform import mongodbatlas_api_key_project_assignment.example <PROJECT_ID>/<API_KEY_ID>
   ```
3. Remove the old resource from state:
   ```shell
   terraform state rm mongodbatlas_project_api_key.example
   ```
4. Remove the old resource block from your configuration (if not already done).
5. Run `terraform plan` and `terraform apply`.

This process ensures that your existing Atlas API keys and assignments are preserved and managed by Terraform under the new resource types, with no deletion or recreation.

## Step-by-Step Migration Example
1. **Review your current configuration:**
   ```hcl
   resource "mongodbatlas_project_api_key" "example" {
     description = "example key"
     project_assignment {
       project_id = var.project_id
       role_names = ["GROUP_READ_ONLY"]
     }
   }
   ```
2. **Add the new resources:**
   ```hcl
   resource "mongodbatlas_api_key" "example" {
     org_id      = var.org_id
     description = "example key"
     role_names  = ["ORG_READ_ONLY"]
   }

   resource "mongodbatlas_api_key_project_assignment" "example" {
     project_id = var.project_id
     api_key_id = mongodbatlas_api_key.example.api_key_id
     role_names = ["GROUP_READ_ONLY"]
   }
   ```
3. **Remove the old resource:**
   - Delete or comment out the `mongodbatlas_project_api_key` block.
4. **Update references:**
   - Change any references from `mongodbatlas_project_api_key.example.api_key_id` to `mongodbatlas_api_key.example.api_key_id`.
5. **Apply the changes:**
   - Run `terraform plan` to review the changes.
   - Run `terraform apply` to apply the migration.

## FAQ
**Q: Can I assign the same API key to multiple projects?**
A: Yes, simply create multiple `mongodbatlas_api_key_project_assignment` resources for each project.

**Q: Do I need to change anything for existing keys?**
A: Existing keys will continue to work, but we recommend following the migration guide to move to the new pattern.

**Q: What if I have many project assignments?**
A: You can use `for_each` or `count` with `mongodbatlas_api_key_project_assignment` to manage multiple assignments efficiently.

**Q: Where can I find a working example?**
A: See [examples/mongodbatlas_api_key/main.tf](../../examples/mongodbatlas_api_key/main.tf).

## Further Resources
- [Terraform MongoDB Atlas Provider Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs)
- [API Key Project Assignment Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/api_key_project_assignment)
- [Migration Guide Source](https://github.com/mongodb/terraform-provider-mongodbatlas) 