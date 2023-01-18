---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project_api_key"
sidebar_current: "docs-mongodbatlas-datasource-project-api-key"
description: |-
    Describes a Project API Key.
---

# Data Source: mongodbatlas_project_api_key

`mongodbatlas_project_api_key` describes a MongoDB Atlas Project API Key. This represents a Project API Key that has been created.

~> **IMPORTANT WARNING:**  Creating, Reading, Updating, or Deleting Atlas API Keys may key expose sensitive organizational secrets to Terraform State. For best security practices consider storing sensitive API Key secrets instead via the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

-> **NOTE:** You may find project_id in the official documentation.

## Example Usage

### Using org_id attribute to query
```terraform
resource "mongodbatlas_project_api_key" "test" {
  description   = "key-name"
  project_id    = "<PROJECT_ID>"
  role_names = ["GROUP_READ_ONLY"]
  }
}

data "mongodbatlas_project_api_key" "test" {
  project_id = "${mongodbatlas_api_key.test.project_id}"
  api_key_id = "${mongodbatlas_api_key.test.api_key_id}"
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `project_id` - Unique identifier for the project whose API keys you want to retrieve. Use the /groups endpoint to retrieve all projects to which the authenticated user has access.
* `description` - Description of this Project API key.
* `public_key` - Public key for this Organization API key.
* `private_key` - Private key for this Organization API key.
* `role_names` - Name of the role. This resource returns all the roles the user has in Atlas.
The following are valid roles:
  * `GROUP_OWNER`
  * `GROUP_READ_ONLY`
  * `GROUP_DATA_ACCESS_ADMIN`
  * `GROUP_DATA_ACCESS_READ_WRITE`
  * `GROUP_DATA_ACCESS_READ_ONLY`
  * `GROUP_CLUSTER_MANAGER`  
    
See [MongoDB Atlas API - API Key](https://www.mongodb.com/docs/atlas/reference/api/projectApiKeys/get-all-apiKeys-in-one-project/) - Documentation for more information.
