---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: access_list_api_key"
sidebar_current: "docs-mongodbatlas-datasource-access-list-api-key"
description: |-
    Displays the access list entries for the specified Atlas Organization API Key. 
---

# Data Source: mongodbatlas_access_list_api_key

`mongodbatlas_access_list_api_key` describes an Access List API Key entry resource. The access list grants access from IPs, CIDRs) to clusters within the Project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

~> **IMPORTANT:**
When you remove an entry from the access list, existing connections from the removed address(es) may remain open for a variable amount of time. How much time passes before Atlas closes the connection depends on several factors, including how the connection was established, the particular behavior of the application or driver using the address, and the connection protocol (e.g., TCP or UDP). This is particularly important to consider when changing an existing IP address or CIDR block as they cannot be updated via the Provider (comments can however), hence a change will force the destruction and recreation of entries.   

~> **IMPORTANT WARNING:** Managing Atlas Programmatic API Keys (PAKs) with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your PAKs via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).


## Example Usage

### Using CIDR Block
```terraform
resource "mongodbatlas_access_list_api_key" "test" {
  org_id = "<ORG-ID>"
  cidr_block = "1.2.3.4/32"
  api_key = "a29120e123cd"
}

data "mongodbatlas_access_list_api_key" "test" {
  org_id     = mongodbatlas_access_list_api_key.test.org_id
  cidr_block = mongodbatlas_access_list_api_key.test.cidr_block
  api_key_id = mongodbatlas_access_list_api_key.test.api_key_id
}
```

### Using IP Address
```terraform
resource "mongodbatlas_access_list_api_key" "test" {
  org_id     = "<ORG-ID>"
  ip_address = "2.3.4.5"
  api_key = "a29120e123cd"
}

data "mongodbatlas_access_list_api_key" "test" {
  org_id = mongodbatlas_access_list_api_key.test.org_id
  ip_address = mongodbatlas_access_list_api_key.test.ip_address
  api_key_id = mongodbatlas_access_list_api_key.test.api_key_id
}
```

## Argument Reference

* `org_id` - (Required) Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
* `cidr_block` - (Optional) Range of IP addresses in CIDR notation to be added to the access list.
* `ip_address` - (Optional) Single IP address to be added to the access list.
* `api_key_id` - (Required) Unique identifier for the Organization API Key for which you want to retrieve an access list entry.
* 
->**NOTE:** You must set either the `cidr_block` attribute or the `ip_address` attribute. Don't set both.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier used by Terraform for internal management and can be used to import.
* `comment` - Comment to add to the access list entry.

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/getApiKeyAccessList)
