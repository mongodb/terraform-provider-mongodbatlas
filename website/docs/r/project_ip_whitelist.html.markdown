---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project_ip_whitelist"
sidebar_current: "docs-mongodbatlas-resource-project-ip-whitelist"
description: |-
    Provides an IP Whitelist resource.
---

# mongodbatlas_project_ip_whitelist

`mongodbatlas_project_ip_whitelist` provides an IP Whitelist entry resource. The whitelist grants access from IPs or CIDRs to clusters within the Project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

~> **IMPORTANT:**
When you remove an entry from the whitelist, existing connections from the removed address(es) may remain open for a variable amount of time. How much time passes before Atlas closes the connection depends on several factors, including how the connection was established, the particular behavior of the application or driver using the address, and the connection protocol (e.g., TCP or UDP). This is particularly important to consider when changing an existing IP address or CIDR block as they cannot be updated via the Provider (comments can however), hence a change will force the destruction and recreation of entries.   


## Example Usage

```hcl
resource "mongodbatlas_project_ip_whitelist" "test" {
  project_id = "<PROJECT-ID>"
  cidr_block = "1.2.3.4/32"
  comment    = "cidr block for tf acc testing"
}
```

```hcl
resource "mongodbatlas_project_ip_whitelist" "test" {
  project_id = "<PROJECT-ID>"
  ip_address = "2.3.4.5"
  comment    = "ip address for tf acc testing"
}
```

## Argument Reference

* `project_id` - (Required) The ID of the project in which to add the whitelist entry.
* `cidr_block` - (Optional) The whitelist entry in Classless Inter-Domain Routing (CIDR) notation. Mutually exclusive with `ip_address`.
* `ip_address` - (Optional) The whitelisted IP address. Mutually exclusive with `cidr_block`.
* `comment` - (Optional) Comment to add to the whitelist entry.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier used for terraform for internal manages and can be used to import.

## Import

IP Whitelist entries can be imported using the `project_id`, e.g.

```
$ terraform import mongodbatlas_project_ip_whitelist.test 5d0f1f74cf09a29120e123cd
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/whitelist/)
