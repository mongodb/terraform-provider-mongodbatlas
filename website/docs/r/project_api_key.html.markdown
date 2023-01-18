---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project_api_key"
sidebar_current: "docs-mongodbatlas-resource-project-api-key"
description: |-
    Provides a Project API Key resource.
---

# Resource: mongodbatlas_project_api_key

`mongodbatlas_project_api_key` provides a Project API Key resource. This allows project API Key to be created.

~> **IMPORTANT WARNING:**  Creating, Reading, Updating, or Deleting Atlas API Keys may key expose sensitive organizational secrets to Terraform State. Consider storing sensitive API Key secrets instead via the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

## Example Usage

```terraform
resource "mongodbatlas_project_api_key" "test" {
  description   = "key-name"
  project_id        = "<PROJECT_ID>"
  role_names = ["GROUP_OWNER"]
  }
}
```

## Argument Reference

* `project__id` - Unique identifier for the project whose API keys you want to retrieve. Use the /orgs endpoint to retrieve all organizations to which the authenticated user has access.
* `description` - Description of this Organization API key.
* `role_names` - (Required) List of Project roles that the Programmatic API key needs to have. Ensure you provide: at least one role and ensure all roles are valid for the Project.  You must specify an array even if you are only associating a single role with the Programmatic API key.
 The following are valid roles:
  * `GROUP_OWNER`
  * `GROUP_READ_ONLY`
  * `GROUP_DATA_ACCESS_ADMIN`
  * `GROUP_DATA_ACCESS_READ_WRITE`
  * `GROUP_DATA_ACCESS_READ_ONLY`
  * `GROUP_CLUSTER_MANAGER`  

~> **NOTE:** Project created by API Keys must belong to an existing organization.

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
