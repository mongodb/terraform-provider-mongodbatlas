---
subcategory: "Programmatic API Keys"
---

# Data Source: mongodbatlas_access_list_api_key

`mongodbatlas_access_list_api_key` describes an Access List API Key entry resource. The access list grants access from IPs or CIDRs to clusters within the Project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

~> **IMPORTANT:**
When you remove an entry from the access list, existing connections from the removed address(es) may remain open for a variable amount of time. How much time passes before Atlas closes the connection depends on several factors, including how the connection was established, the particular behavior of the application or driver using the address, and the connection protocol (e.g., TCP or UDP). This is particularly important to consider when changing an existing IP address or CIDR block as they cannot be updated via the Provider (comments can however), hence a change will force the destruction and recreation of entries.   

~> **IMPORTANT WARNING:** Managing Atlas Programmatic API Keys (PAKs) with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your PAKs via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).


## Example Usage

```terraform
resource "mongodbatlas_access_list_api_key" "test" {
  org_id     = "<ORG_ID>"
  ip_address = "2.3.4.5"
  api_key_id = "a29120e123cd"
}

data "mongodbatlas_access_list_api_key" "test" {
  org_id     = mongodbatlas_access_list_api_key.test.org_id
  ip_address = mongodbatlas_access_list_api_key.test.ip_address
  api_key_id = mongodbatlas_access_list_api_key.test.api_key_id
}
```

## Argument Reference

* `org_id` - (Required) Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
* `api_key_id` - (Required) Unique identifier for the Organization API Key for which you want to retrieve an access list entry.
* `ip_address` - (Required) IP address to retrieve from the access list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cidr_block` - Range of IP addresses in CIDR notation included in the API Key access list.
* `created` - Date and time when the access list entry was created.
* `access_count` - Total number of requests that have originated from this IP address or CIDR block.
* `last_used` - Date and time when the API key was last used from this IP address or CIDR block.
* `last_used_address` - IP address from which the last API request was made.

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/getApiKeyAccessList)
