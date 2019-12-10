---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: organization"
sidebar_current: "docs-mongodbatlas-resource-organization"
description: |-
		Provides a Organization resource.
---

# mongodb_atlas_organization

`mongodbatlas_organization` provides a Organization resource. This allows organization to be created.

## Example Usage

```hcl
resource "mongodbatlas_organization" "my_organization" {
	name   = "testacc-organization"
}
```

## Argument Reference

* `name` - (Required) The name of the organization you want to create.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for the organization.
* `name` - The name of the organization.

## Import

Organization must be imported using organization ID, e.g.

```
$ terraform import mongodbatlas_organization.my_organization 5d09d6a59ccf6445652a444a
```
For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/organizations/)