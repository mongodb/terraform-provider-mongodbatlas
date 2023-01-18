---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: access_list_api_key"
sidebar_current: "docs-mongodbatlas-resource-access_list-api-key"
description: |-
    Provides an Access List API Key resource.
---

# Resource: mongodbatlas_project_ip_access_list

`mongodbatlas_access_list_api_key` provides an IP Access List entry resource. The access list grants access from IPs, CIDRs or AWS Security Groups (if VPC Peering is enabled) to clusters within the Project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

~> **IMPORTANT:**
When you remove an entry from the access list, existing connections from the removed address(es) may remain open for a variable amount of time. How much time passes before Atlas closes the connection depends on several factors, including how the connection was established, the particular behavior of the application or driver using the address, and the connection protocol (e.g., TCP or UDP). This is particularly important to consider when changing an existing IP address or CIDR block as they cannot be updated via the Provider (comments can however), hence a change will force the destruction and recreation of entries.  

~> **IMPORTANT WARNING:** Creating, Reading, Updating, or Deleting Atlas API Keys may key expose sensitive organizational secrets to Terraform State. Consider storing sensitive API Key secrets instead via the HashiCorp Vault MongoDB Atlas Secrets Engine.


## Example Usage

### Using CIDR Block
```terraform
resource "mongodbatlas_access_list_api_key" "test" {
  org_id = "<ORG-ID>"
  cidr_block = "1.2.3.4/32"
  api_key_id = "a29120e123cd"
}
```

### Using IP Address
```terraform
resource "mongodbatlas_access_list_api_key" "test" {
  org_id = "<ORG-ID>"
  ip_address = "2.3.4.5"
  api_key_id = "a29120e123cd"
}
```

## Argument Reference

* `org_id` - (Required) Unique identifier for the organinzation to which you want to add one or more access list entries.
* `cidr_block` - (Optional) Range of IP addresses in CIDR notation to be added to the access list. Your access list entry can include only one `cidrBlock`, or one `ipAddress`.
* `ip_address` - (Optional) Single IP address to be added to the access list.
* `api_key_id` - Unique identifier for the Organization API Key for which you want to create a new access list entry.

-> **NOTE:** One of the following attributes must set: `cidr_block`  or `ip_address`.

## Import

IP Access List entries can be imported using the `org_id` , `api_key_id` and `cidr_block` or `ip_address`, e.g.

```
$ terraform import mongodbatlas_access_list_api_key.test 5d0f1f74cf09a29120e123cd-a29120e123cd-10.242.88.0/21
```

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys/operation/returnAllAccessListEntriesForOneOrganizationApiKey)
