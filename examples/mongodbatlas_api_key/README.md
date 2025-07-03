# MongoDB Atlas Provider — Programmatic API Key Example

This example demonstrates the **recommended approach** for managing MongoDB Atlas Programmatic API Keys (PAKs) using the Terraform MongoDB Atlas Provider:

- **Create API keys independently** using the `mongodbatlas_api_key` resource.
- **Assign API keys to projects** using the `mongodbatlas_api_key_project_assignment` resource.

This pattern provides greater flexibility and clarity, allowing you to manage API keys and their project assignments separately.

## Example Overview

The included [`main.tf`](./main.tf) shows how to:

1. **Create an API key** at the organization level with `mongodbatlas_api_key`.
2. **Create a project** with `mongodbatlas_project`.
3. **Assign the API key to the project** with `mongodbatlas_api_key_project_assignment`, specifying project-level roles.

```hcl
resource "mongodbatlas_api_key" "test" {
  org_id      = var.org_id
  description = "Test API Key"
  role_names  = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_project" "test1" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_api_key_project_assignment" "test1" {
  project_id  = mongodbatlas_project.test1.id
  api_key_id  = mongodbatlas_api_key.test.api_key_id
  role_names  = ["GROUP_OWNER"]
}
```

## Why Use This Pattern?
- **Flexibility:** Assign the same API key to multiple projects or update assignments independently.
- **Clarity:** Separate resources for key creation and project assignment.
- **Future-proof:** This is the preferred and supported method going forward.

## Migrating from the Old Pattern
If you are currently using `mongodbatlas_project_api_key` resource, see the [Migration Guide](../../docs/guides/project-api-key-migration.md) for step-by-step instructions on updating your configuration.

## Additional Notes
- All API keys in Atlas are organization-level keys. Assigning them to a project grants project-level roles for that project.
- For more details, see the [Terraform MongoDB Atlas Provider documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/api_key_project_assignment).
