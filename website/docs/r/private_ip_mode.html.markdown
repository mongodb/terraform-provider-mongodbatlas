---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_ip_mode"
sidebar_current: "docs-mongodbatlas-resource-private-ip-mode"
description: |-
    Provides a Private IP Mode resource.
---

# mongodbatlas_private_ip_mode

`mongodbatlas_private_ip_mode` provides a Private IP Mode resource. This allows one to disable Connect via Peering Only mode for a MongoDB Atlas Project.

~> **Deprecated Feature**: <br> This feature has been deprecated. Use [Split Horizon connection strings](https://dochub.mongodb.org/core/atlas-horizon-faq) to connect to your cluster. These connection strings allow you to connect using both VPC/VNet Peering and whitelisted public IP addresses. To learn more about support for Split Horizon, see [this FAQ](https://dochub.mongodb.org/core/atlas-horizon-faq). You need this endpoint to [disable Peering Only](https://docs.atlas.mongodb.com/reference/faq/connection-changes/#disable-peering-mode).


## Example Usage

```terraform
resource "mongodbatlas_private_ip_mode" "my_private_ip_mode" {
    project_id = "<YOUR PROJECT ID>"
	enabled    = false
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to enable Only Private IP Mode.
* `enabled` - (Required) Indicates whether Connect via Peering Only mode is enabled or disabled for an Atlas project


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project id.

## Import

Project must be imported using project ID, e.g.

```
$ terraform import mongodbatlas_private_ip_mode.my_private_ip_mode 5d09d6a59ccf6445652a444a
```
For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/get-private-ip-mode-for-project/)
