# Data Source: mongodbatlas_federated_query_limits

`mongodbatlas_federated_query_limits` provides a Federated Database Instance Query Limits data source. To learn more about Atlas Data Federation see https://www.mongodb.com/docs/atlas/data-federation/overview/. 

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usages


```terraform
data "mongodbatlas_federated_query_limits" "test" {
  project_id = "PROJECT_ID"
  tenant_name = "FEDERATED_DATABASE_INSTANCE_NAME"
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create a Federated Database Instance.
* `tenant_name` - (Required) Name of the Atlas Federated Database Instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `results` - Includes Federated Database instance query limits for each item detailed in the results array section.

### Federated Database Instance Query Limit

* `id` - The Terraform's unique identifier used internally for state management.
* `current_usage` - Amount that indicates the current usage of the limit.
* `default_limit` - Default value of the limit.
* `lastModifiedDate` - Only used for Data Federation limits. Timestamp that indicates when this usage limit was last modified. This field uses the ISO 8601 timestamp format in UTC.
* `maximumLimit` - Maximum value of the limit.
* `name` - Name that identifies the user-managed limit to modify.

## Import

The Federated Database Instance Query Limit can be imported using project ID, name of the instance and limit name, in the format: 
`project_id`--`tenant_name`--`limit_name`, e.g.

```
$ terraform import mongodbatlas_federated_query_limit.example 1112222b3bf99403840e8934--FederatedDatabaseInstance0--bytesProcessed.daily
```

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation) Documentation for more information.
