---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project_api_key"
sidebar_current: "docs-mongodbatlas-resource-project-api-key"
description: |-
    Creates and assigns the specified Atlas Organization API Key to the specified Project. Users with the Project Owner role in the project associated with the API key can use the organization API key to access the resources.
---

# Resource: mongodbatlas_project_api_key

`mongodbatlas_project_api_key` provides a Project API Key resource. This allows project API Key to be created.

~> **IMPORTANT WARNING:** Managing Atlas Programmatic API Keys (PAKs) with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your PAKs via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

## Example Usage - Create and Assign PAK Together

```terraform
resource "mongodbatlas_project_api_key" "test" {
  description   = "Description of your API key"
  project_id    = "64259ee860c43338194b0f8e"
  role_names    = ["GROUP_OWNER"]
}
```

## Example Usage - Create and Assign PAK to Multiple Projects

```terraform
resource "mongodbatlas_project_api_key" "test" {
  description   = "Description of your API key"
  project_id  = "64259ee860c43338194b0f8e"
  
 project_assignment {
    project_id = "64259ee860c43338194b0f8e"
    role_names = ["GROUP_READ_ONLY", "GROUP_OWNER"]
  }
  
  project_assignment {
    project_id = "64229ee820c42228194b0f4a"
    role_names = ["GROUP_READ_ONLY"]
  }
  
}
```

## Example Usage - Create Org PAK and Assign it to Multiple Projects

```terraform
resource "mongodbatlas_project" "atlas-project" {
  name   = "ProjectTest"
  org_id = "60ddf55c27a5a20955a707d7"
}

resource "mongodbatlas_project_api_key" "api_1" {
  description = "test api_key multi"
  project_id  = mongodbatlas_project.atlas-project.id

  // NOTE: The `project_id` of the first `project_assignment` element must be the same as the `project_id` of the resource.
  project_assignment {
    project_id = mongodbatlas_project.atlas-project.id
    role_names = ["ORG_BILLING_ADMIN", "GROUP_READ_ONLY"]
  }

  project_assignment {
    project_id = "63dcfc256af00a5934e60924"
    role_names = ["GROUP_READ_ONLY"]
  }

  project_assignment {
    project_id = "64c23af6f133166c39176cbf"
    role_names = ["GROUP_OWNER"]
  }
}
```

## Argument Reference

* `project_id` -Unique 24-hexadecimal digit string that identifies your project.
* `description` - Description of this Project API key.

~> **NOTE:** Project created by API Keys must belong to an existing organization.

### project_assignment
List of Project roles that the Programmatic API key needs to have. `project_assignment` attribute is optional.

* `project_id` - (Required) Project ID to assign to Access Key
* `role_names` - (Required) List of Project roles that the Programmatic API key needs to have. Ensure you provide: at least one role and ensure all roles are valid for the Project. You must specify an array even if you are only associating a single role with the Programmatic API key. The [MongoDB Documentation](https://www.mongodb.com/docs/atlas/reference/user-roles/#project-roles) describes the valid roles that can be assigned.

~> **NOTE:** The `project_id` of the first `project_assignment` element must be the same as the `project_id` of the resource.

~> **NOTE:** The organization level roles can be defined only in the first `project_assignment` element.

~> **NOTE:** The `ORG_READ_ONLY` role at the organization level is invalid in this context. When the `project_assignment``` lacks organizational roles, the `mongodbatlas_project_api_key` resource generates an organization API key with the `ORG_READ_ONLY` role and associates it with `GROUP_*` roles. Consequently, the resource does not permit the use of `ORG_READ_ONLY` to ensure consistency between configuration and state.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_key_id` - Unique identifier for this Project API key.

## Import

API Keys must be imported using org ID, API Key ID e.g.

```
$ terraform import mongodbatlas_project_api_key.test 5d09d6a59ccf6445652a444a-6576974933969669
```
See [MongoDB Atlas API - API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/createAndAssignOneOrganizationApiKeyToOneProject) - Documentation for more information.
