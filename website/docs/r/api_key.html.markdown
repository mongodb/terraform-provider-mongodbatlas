---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: api_key"
sidebar_current: "docs-mongodbatlas-resource-api-key"
description: |-
    Provides a API Key resource.
---

# Resource: mongodbatlas_api_key

`mongodbatlas_api_key` provides a Organization API key resource. This allows an Organizational API key to be created.

~> **IMPORTANT WARNING:** Managing Atlas Programmatic API Keys (PAKs) with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your PAKs via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

## Example Usage

```terraform
resource "mongodbatlas_api_key" "test" {
  description   = "key-name"
  org_id        = "<ORG_ID>"
  role_names = ["ORG_READ_ONLY"]
  }
}
```

## Argument Reference

* `org_id` - Unique identifier for the organization whose API keys you want to retrieve. Use the /orgs endpoint to retrieve all organizations to which the authenticated user has access.
* `description` - Description of this Organization API key.
* `role_names` - Name of the role. This resource returns all the roles the user has in Atlas.
The following are valid roles:
  * `ORG_OWNER`
  * `ORG_GROUP_CREATOR`
  * `ORG_BILLING_ADMIN`
  * `ORG_READ_ONLY`
  * `ORG_MEMBER`

~> **NOTE:** Project created by API Keys must belong to an existing organization.

### Programmatic API Keys
api_keys allows one to assign an existing organization programmatic API key to a Project. The api_keys attribute is optional.

* `api_key_id` - (Required) The unique identifier of the Programmatic API key you want to associate with the Project.  The Programmatic API key and Project must share the same parent organization.  Note: this is not the `publicKey` of the Programmatic API key but the `id` of the key. See [Programmatic API Keys](https://docs.atlas.mongodb.com/reference/api/apiKeys/) for more.

* `role_names` - (Required) List of Project roles that the Programmatic API key needs to have. Ensure you provide: at least one role and ensure all roles are valid for the Project.  You must specify an array even if you are only associating a single role with the Programmatic API key.
 The following are valid roles:
  * `GROUP_OWNER`
  * `GROUP_READ_ONLY`
  * `GROUP_DATA_ACCESS_ADMIN`
  * `GROUP_DATA_ACCESS_READ_WRITE`
  * `GROUP_DATA_ACCESS_READ_ONLY`
  * `GROUP_CLUSTER_MANAGER`  
 ## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_key_id` - Unique identifier for this Organization API key.
## Import

API Keys must be imported using org ID, API Key ID e.g.

```
$ terraform import mongodbatlas_api_key.test 5d09d6a59ccf6445652a444a-6576974933969669
```
See [MongoDB Atlas API - API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/createOneOrganizationApiKey) - Documentation for more information.
