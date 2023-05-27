---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: organization"
sidebar_current: "docs-mongodbatlas-datasource-organization"
description: |-
    Describes an Organization.
---

# Data Source: mongodbatlas_organization

`mongodbatlas_organization` describes a MongoDB Atlas Organization. This represents a organization that has been created.

## Example Usage

### Using project_id attribute to query
```terraform

data "mongodbatlas_organization" "test" {
  org_id = "<org_id>"
}
```

## Argument Reference

* `org_id` - Unique 24-hexadecimal digit string that identifies the organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Human-readable label that identifies the organization.
* `id` - Unique 24-hexadecimal digit string that identifies the organization.
* `is_deleted` - Flag that indicates whether this organization has been deleted.

  
See [MongoDB Atlas API - Organization](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/getOrganization) Documentation for more information.
