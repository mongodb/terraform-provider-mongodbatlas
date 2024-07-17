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

 ## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_key_id` - Unique identifier for this Organization API key.
## Import

API Keys must be imported using org ID, API Key ID e.g.

```
$ terraform import mongodbatlas_api_key.test 5d09d6a59ccf6445652a444a-6576974933969669
```
See [MongoDB Atlas API - API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/createApiKey) Documentation for more information.
