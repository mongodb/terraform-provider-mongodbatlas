# Resource: mongodbatlas_project_api_key

`mongodbatlas_project_api_key` provides a Project API Key resource. This allows project API Key to be created.

~> **IMPORTANT NOTE** The use of `mongodbatlas_project_api_key` resource is no longer the recommended approach. For new configurations, we recommend using the `mongodbatlas_api_key` resource and the `mongodbatlas_api_key_project_assignment` resource to assign the API Keys to projects. This approach is more flexible and aligns with best practices. For existing configurations, the migration to the new pattern is **not required**. If you want to migrate, see the [Migration Guide](../guides/project-api-key-migration.md) for step-by-step instructions on migrating from `mongodbatlas_project_api_key` to the new pattern.

~> **IMPORTANT WARNING:** Managing Atlas Programmatic API Keys (PAKs) with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your PAKs via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

## Example Usage - Create and Assign PAK Together

```terraform
resource "mongodbatlas_project_api_key" "test" {
  description   = "Description of your API key"
  project_assignment {
    project_id = "64259ee860c43338194b0f8e"
    role_names = ["GROUP_OWNER"]
  }
}
```

## Example Usage - Create and Assign PAK to Multiple Projects

```terraform
resource "mongodbatlas_project_api_key" "test" {
  description   = "Description of your API key"
  
  project_assignment {
    project_id = "64259ee860c43338194b0f8e"
    role_names = ["GROUP_READ_ONLY", "GROUP_OWNER"]
  }
  
  project_assignment {
    project_id = "74259ee860c43338194b0f8e"
    role_names = ["GROUP_READ_ONLY"]
  }
  
}
```

## Argument Reference

* `description` - (Required) Description of this Project API key.

~> **NOTE:** Project created by API Keys must belong to an existing organization.

### project_assignment
List of Project roles that the Programmatic API key needs to have. At least one `project_assignment` block must be defined.

* `project_id` - (Required) Project ID to assign to Access Key
* `role_names` - (Required) List of Project roles that the Programmatic API key needs to have. Ensure you provide: at least one role and ensure all roles are valid for the Project. You must specify an array even if you are only associating a single role with the Programmatic API key. The [MongoDB Documentation](https://www.mongodb.com/docs/atlas/reference/user-roles/#project-roles) describes the valid roles that can be assigned.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_key_id` - Unique identifier for this Project API key.

## Import

API Keys must be imported using project ID, API Key ID e.g.

```
$ terraform import mongodbatlas_project_api_key.test 5d09d6a59ccf6445652a444a-6576974933969669
```
See [MongoDB Atlas API - API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/createProjectApiKey) - Documentation for more information.
