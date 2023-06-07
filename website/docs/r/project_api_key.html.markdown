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
  description   = "key-name"
  project_id    = "<PROJECT_ID>"
  role_names    = ["GROUP_OWNER"]
}
```

## Example Usage - Create and Assign PAK to Multiple Projects

```terraform
resource "mongodbatlas_api_key" "test" {
  description   = "key-name"
  org_id        = "<ORG_ID>"
  
 project_assignment {
    project_id = <project_id>
    role_names = ["GROUP_READ_ONLY", "GROUP_OWNER"]
  }
  
  project_assignment {
    project_id = <additional_project_id>
    role_names = ["GROUP_READ_ONLY"]
  }
  
}
```

## Argument Reference

* `project_id` -Unique 24-hexadecimal digit string that identifies your project.
* `description` - Description of this Organization API key.
* `role_names` - (Deprecated use project_assignment) List of Project roles that the Programmatic API key needs to have. Ensure you provide: at least one role and ensure all roles are valid for the Project.  You must specify an array even if you are only associating a single role with the Programmatic API key.
 The following are valid roles:
  * `GROUP_OWNER`
  * `GROUP_READ_ONLY`
  * `GROUP_DATA_ACCESS_ADMIN`
  * `GROUP_DATA_ACCESS_READ_WRITE`
  * `GROUP_DATA_ACCESS_READ_ONLY`
  * `GROUP_CLUSTER_MANAGER`  

~> **NOTE:** Project created by API Keys must belong to an existing organization.

# project_assignment
Project Assignment attribute is optional (Use project_assignment going forward as role_names parameter above is deprecated)

* `project_id` - (Required) Project ID to assign to Access Key
* `role_names` - Name of the role. This resource returns all the roles the user has in Atlas.
The following are valid roles:
  * `GROUP_OWNER`
  * `GROUP_READ_ONLY`
  * `GROUP_DATA_ACCESS_ADMIN`
  * `GROUP_DATA_ACCESS_READ_WRITE`
  * `GROUP_DATA_ACCESS_READ_ONLY`
  * `GROUP_CLUSTER_MANAGER`  
  * 
### Programmatic API Keys
api_keys allows one to assign an existing organization programmatic API key to a Project. The api_keys attribute is optional.

* `api_key_id` - (Required) The unique identifier of the Programmatic API key you want to associate with the Project.  The Programmatic API key and Project must share the same parent organization.  Note: this is not the `publicKey` of the Programmatic API key but the `id` of the key. See [Programmatic API Keys](https://docs.atlas.mongodb.com/reference/api/apiKeys/) for more.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_key_id` - Unique identifier for this Project API key.

## Import

API Keys must be imported using org ID, API Key ID e.g.

```
$ terraform import mongodbatlas_project_api_key.test 5d09d6a59ccf6445652a444a-6576974933969669
```
See [MongoDB Atlas API - API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/createAndAssignOneOrganizationApiKeyToOneProject) - Documentation for more information.
