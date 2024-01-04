---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: api_key"
sidebar_current: "docs-mongodbatlas-datasource-api-key"
description: |-
    Describes a API Key.
---

# Data Source: mongodbatlas_api_key

`mongodbatlas_api_key` describes a MongoDB Atlas API Key. This represents a API Key that has been created.

~> **IMPORTANT WARNING:** Managing Atlas Programmatic API Keys (PAKs) with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your PAKs via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

-> **NOTE:** You may find org_id in the official documentation.

## Example Usage

### Using org_id attribute to query
```terraform
resource "mongodbatlas_api_key" "test" {
  description   = "key-name"
  org_id        = "<ORG_ID>"
  role_names = ["ORG_READ_ONLY"]
  }
}

data "mongodbatlas_api_key" "test" {
  org_id = "${mongodbatlas_api_key.test.org_id}"
  api_key_id = "${mongodbatlas_api_key.test.api_key_id}"
}
```

## Argument Reference

* `org_id` - (Required) Unique identifier for the organization.
* `api_key_id` - (Required) Unique identifier for the organization API key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Description of this Organization API key.
* `public_key` - Public key for this Organization API key.
* `private_key` - Private key for this Organization API key.
* `role_names` - Name of the role. This resource returns all the roles the user has in Atlas.
The following are valid roles:
  * `ORG_OWNER`
  * `ORG_GROUP_CREATOR`
  * `ORG_BILLING_ADMIN`
  * `ORG_READ_ONLY`
  * `ORG_MEMBER`
    
See [MongoDB Atlas API - Return One Organization API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/getApiKey) - Documentation for more information.
