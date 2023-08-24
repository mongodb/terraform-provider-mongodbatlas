---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: federated_database_query_limit"
sidebar_current: "docs-mongodbatlas-resource-federated-query-limit"
description: |-
    Provides a Federated Database Instance Query Limit.
---

# Resource: mongodbatlas_federated_query_limit

`mongodbatlas_federated_query_limit` provides a Federated Database Instance Query Limits resource. To learn more about Atlas Data Federation see https://www.mongodb.com/docs/atlas/data-federation/overview/.


-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usages


```terraform
resource "mongodbatlas_federated_query_limit" "test" {
  project_id = "64707f06c519c20c3a2b1b03"
  tenant_name = "FederatedDatabseInstance0"
  limit_name = "bytesProcessed.weekly"
  overrun_policy = "BLOCK"
  value          = 5147483648
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create a Federated Database Instance.
* `tenant_name` - (Required) Name of the Atlas Federated Database Instance.
* `limit_name` - (Required) String enum that indicates whether the identity provider is active or not. Accepted values are:
    * `bytesProcessed.query`: Limit on the number of bytes processed during a single data federation query.
    * `bytesProcessed.daily`: Limit on the number of bytes processed for the data federation instance for the current day.
    * `bytesProcessed.weekly`: Limit on the number of bytes processed for the data federation instance for the current week.
    * `bytesProcessed.monthly`: Limit on the number of bytes processed for the data federation instance for the current month.
* `overrun_policy` - (Required) String enum that identifies action to take when the usage limit is exceeded. If limit span is set to QUERY, this is ignored because MongoDB Cloud stops the query when it exceeds the usage limit. Accepted values are "BLOCK" OR "BLOCK_AND_KILL"
* `value` - (Required) Amount to set the limit to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

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

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createOneDataFederationQueryLimit) Documentation for more information.
