---
subcategory: "Organizations"
---

# Data Source: mongodbatlas_roles_org_id

`mongodbatlas_roles_org_id` allows to retrieve the Org ID of the authenticated user.

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

* `org_id` - The ID of the organization you want to retrieve, which is associated with the Service Account or Programmatic API Key (PAK) of the authenticated user.
  
See [MongoDB Atlas API - Role Org ID](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Root/operation/getSystemStatus) Documentation for more information.
