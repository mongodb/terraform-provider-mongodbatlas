---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_ip_mode"
sidebar_current: "docs-mongodbatlas-resource-private-ip-mode"
description: |-
    Provides a Private IP Mode resource.
---

# mongodbatlas_private_ip_mode

`mongodbatlas_private_ip_mode` provides a Private IP Mode resource. This allows one to enable/disable Connect via Peering Only mode for a MongoDB Atlas Project.


~> **IMPORTANT**: <br>**What is Connect via Peering Only Mode?** <br>Connect via Peering Only mode prevents clusters in an Atlas project from connecting to any network destination other than an Atlas Network Peer. Connect via Peering Only mode applies only to **GCP** and **Azure-backed** dedicated clusters. This setting disables the ability to: <br><br>• Deploy non-GCP or Azure-backed dedicated clusters in an Atlas project, and
<br>• Use MongoDB Stitch with dedicated clusters in an Atlas project.


-> **NOTE:** You should create one private_ip_mode per project.

## Example Usage

```hcl
resource "mongodbatlas_private_ip_mode" "my_private_ip_mode" {
    project_id = "<YOUR PROJECT ID>"
	enabled    = true
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to enable Only Private IP Mode.
* `enabled` - (Required) Indicates whether Connect via Peering Only mode is enabled or disabled for an Atlas project.




## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project id.

## Import

Project must be imported using project ID, e.g.

```
$ terraform import mongodbatlas_private_ip_mode.my_private_ip_mode 5d09d6a59ccf6445652a444a
```
For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/get-private-ip-mode-for-project/)