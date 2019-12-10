---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: organization"
sidebar_current: "docs-mongodbatlas-datasource-organization"
description: |-
    Describes an Organization.
---

# mongodbatlas_organization

`mongodbatlas_organization` describes an Organization.


## Example Usage

```hcl
resource "mongodbatlas_organization" "test" {
	name = "organizationName"
}

data "mongodbatlas_organization" "test" {
	org_id = "${mongodbatlas_organization.test.id}"
}
```

## Argument Reference

* `org_id` - (Required) The unique identifier for the organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the organization.

See detailed information for arguments and attributes: [MongoDB API Organization](https://docs.atlas.mongodb.com/reference/api/organization-get-one/)