---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: roles_org_id"
sidebar_current: "docs-mongodbatlas-datasource-roles-org-id"
description: |-
    Describes a Roles Org ID.
---

# Data Source: mongodbatlas_roles_org_id

`mongodbatlas_roles_org_id` describes a MongoDB Atlas Roles Org ID. This represents a Roles Org ID.

## Example Usage

### Using data source to query
```terraform
data "mongodbatlas_roles_org_id" "test" {
}
	
output "org_id" {
	value = data.mongodbatlas_roles_org_id.test.org_id
}
```

## Argument Reference

* No parameters required

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `org_id` - The ID of the organization you want to retrieve associated to an API Key.
  
See [MongoDB Atlas API - Role Org ID](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Root/operation/getSystemStatus) Documentation for more information.
