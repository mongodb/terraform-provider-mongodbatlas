---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project_ip_whitelist"
sidebar_current: "docs-mongodbatlas-resource-project_ip_whitelist"
description: |-
    Provides an IP Whitelist resource.
---

# mongodbatlas_project_ip_whitelist

`mongodbatlas_project_ip_whitelist` provides an IP Whitelist entry resource. The whitelist grants access from IPs or CIDRs to clusters within the Project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_project_ip_whitelist" "cidr" {
  project_id = "<PROJECT-ID>"
  cidr_block = "10.0.0.0/21"
  comment    = "cidr to show in example"
}

resource "mongodbatlas_project_ip_whitelist" "ip" {
  project_id      = "<PROJECT-ID>"
  ip_address = "10.10.10.10"
  comment    = "ip to show in example"
}
```

## Argument Reference

* `project_id` - (Required) The ID of the project in which to add the whitelist entry.
* `cidr_block` - (Optional) The whitelist entry in Classless Inter-Domain Routing (CIDR) notation. Mutually exclusive with `ip_address`.
* `ip_address` - (Optional) The whitelisted IP address. Mutually exclusive with `cidr_block`.
* `comment` - (Optional) Comment to add to the whitelist entry.

-> **IMPORTANT** You cannot set AWS security groups as temporary whitelist entries.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource (Either cidr_block or ip_address).

## Import

IP Whitelist entries can be imported using project ID and CIDR or IP, in the format `PROJECTID-CIDR`, e.g.

```
$ terraform import mongodbatlas_database_user.my_user 1112222b3bf99403840e8934-10.0.0.0/24
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/whitelist/)