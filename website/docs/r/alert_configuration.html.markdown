---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: alert_configuration"
sidebar_current: "docs-mongodbatlas-resource-alert-configuration"
description: |-
    Provides an Alert Configuration resource.
---

# mongodbatlas_alert_configuration

`mongodbatlas_alert_configuration` provides an Alert Configuration resource to define the conditions that trigger an alert and the methods of notification within a MongoDB Atlas project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_alert_configuration" "test" {
  project_id = "<PROJECT-ID>"
  event_type = "OUTSIDE_METRIC_THRESHOLD"
  enabled    = true

  notification {
    type_name     = "GROUP"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = false
    email_enabled = true
    roles         = ["GROUP_CHARTS_ADMIN", "GROUP_CLUSTER_MANAGER"]
  }

  matcher {
    field_name = "HOSTNAME_AND_PORT"
    operator   = "EQUALS"
    value      = "SECONDARY"
  }

  metric_threshold = {
    metric_name = "ASSERT_REGULAR"
    operator    = "LESS_THAN"
    threshold   = 99.0
    units       = "RAW"
    mode        = "AVERAGE"
  }
}
```

## Argument Reference

* `project_id` - (Required) The ID of the project where the alert configuration will create.
* `enabled` - It is not required, but If the attribute is omitted, by default will be false, and the configuration would be disabled. You must set true to enable the configuration.
* `event_type` - (Required) The type of event that will trigger an alert.
  Alert type 	Possible values:
    * Host 	
      - `OUTSIDE_METRIC_THRESHOLD`
      - `HOST_RESTARTED`
      - `HOST_UPGRADED`
      - `HOST_NOW_SECONDARY`
      - `HOST_NOW_PRIMARY`
    * Replica set 	
      - `NO_PRIMARY`
      - `TOO_MANY_ELECTIONS`
    * Sharded cluster 	
      - `CLUSTER_MONGOS_IS_MISSING`
      - `User` 	
      - `JOINED_GROUP`
      - `REMOVED_FROM_GROUP`
      - `USER_ROLES_CHANGED_AUDIT`
    * Project 	
      - `USERS_AWAITING_APPROVAL`
      - `USERS_WITHOUT_MULTI_FACTOR_AUTH`
      - `GROUP_CREATED`
    * Team 	
      - `JOINED_TEAM`
      - `REMOVED_FROM_TEAM`
    * Organization 	
      - `INVITED_TO_ORG`
      - `JOINED_ORG`
    * Data Explorer 	
      - `DATA_EXPLORER`
      - `DATA_EXPLORER_CRUD`
    * Billing 	
      - `CREDIT_CARD_ABOUT_TO_EXPIRE`
      - `CHARGE_SUCCEEDED`
      - `INVOICE_CLOSED`

    -> **NOTE:** If this is set to OUTSIDE_METRIC_THRESHOLD, the metricThreshold field must also be set.

### Matchers
Rules to apply when matching an object against this alert configuration. Only entities that match all these rules are checked for an alert condition. You can filter using the matchers array only when the eventTypeName specifies an event for a host, replica set, or sharded cluster.

* `field_name` - Name of the field in the target object to match on.
  Host alerts support these fields:
    - `TYPE_NAME`
    - `HOSTNAME`
    - `PORT`
    - `HOSTNAME_AND_PORT`
    - `REPLICA_SET_NAME`
  Replica set alerts support these fields:
    - `REPLICA_SET_NAME`
    - `SHARD_NAME`
    - `CLUSTER_NAME`
  Sharded cluster alerts support these fields:
    - `CLUSTER_NAME`
    - `SHARD_NAME`

  All other types of alerts do not support matchers.

* `operator` - If omitted, the configuration is disabled.
* `value` - If omitted, the configuration is disabled.


* `operator` - The operator to test the field’s value. 
  Accepted values are:
    - `EQUALS`
    - `NOT_EQUALS`
    - `CONTAINS`
    - `NOT_CONTAINS`
    - `STARTS_WITH`
    - `ENDS_WITH`
    - `REGEX`

* `value` - Value to test with the specified operator. If `field_name` is set to TYPE_NAME, you can match on the following values:
    - `PRIMARY`
    - `SECONDARY`
    - `STANDALONE`
    - `CONFIG`
    - `MONGOS`

### Metric Threshold
The threshold that causes an alert to be triggered. Required if `event_type_name` : "OUTSIDE_METRIC_THRESHOLD".

* `metric_name` - Name of the metric to check.
* `operator` - Operator to apply when checking the current metric value against the threshold value. 
  Accepted values are:
    - `GREATER_THAN`
    - `LESS_THAN`

* `threshold` - Threshold value outside of which an alert will be triggered.
* `units` - The units for the threshold value. Depends on the type of metric.
  Accepted values are:
    - `RAW`
    - `BITS`
    - `BYTES`
    - `KILOBITS`
    - `KILOBYTES`
    - `MEGABITS`
    - `MEGABYTES`
    - `GIGABITS`
    - `GIGABYTES`
    - `TERABYTES`
    - `PETABYTES`
    - `MILLISECONDS`
    - `SECONDS`
    - `MINUTES`
    - `HOURS`
    - `DAYS`

* `mode` - This must be set to AVERAGE. Atlas computes the current metric value as an average.

### Notifications
Notifications to send when an alert condition is detected.

* `api_token` - Slack API token. Required for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
* `channel_name` - Slack channel name. Required for the SLACK notifications type.
* `datadog_api_key` - Datadog API Key. Found in the Datadog dashboard. Required for the DATADOG notifications type.
* `datadog_region` - Region that indicates which API URL to use. Accepted regions are: `US`, `EU`. The default Datadog region is US.
* `delay_min` - Number of minutes to wait after an alert condition is detected before sending out the first notification.
* `email_address` - Email address to which alert notifications are sent. Required for the EMAIL notifications type.
* `email_enabled` - Flag indicating if email notifications should be sent. Configurable for `ORG`, `GROUP`, and `USER` notifications types.
* `flowdock_api_token` - The Flowdock personal API token. Required for the `FLOWDOCK` notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
* `flow_name` - Flowdock flow name in lower-case letters. Required for the `FLOWDOCK` notifications type
* `interval_min` - Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved. The minimum value is 5.
* `mobile_number` - Mobile number to which alert notifications are sent. Required for the SMS notifications type.
* `ops_genie_api_key` - Opsgenie API Key. Required for the `OPS_GENIE` notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
* `ops_genie_region` - Region that indicates which API URL to use. Accepted regions are: `US` ,`EU`. The default Opsgenie region is US.
* `org_name` - Flowdock organization name in lower-case letters. This is the name that appears after www.flowdock.com/app/ in the URL string. Required for the FLOWDOCK notifications type.
* `service_key` - PagerDuty service key. Required for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
* `sms_enabled` - Flag indicating if text message notifications should be sent. Configurable for `ORG`, `GROUP`, and `USER` notifications types.
* `team_id` - Unique identifier of a team.
* `type_name` - Type of alert notification. 
  Accepted values are:
    - `DATADOG`
    - `EMAIL`
    - `FLOWDOCK`
    - `GROUP` (Project)
    - `OPS_GENIE`
    - `ORG`
    - `PAGER_DUTY`
    - `SLACK`
    - `SMS`
    - `TEAM`
    - `USER`
    - `VICTOR_OPS`
    - `WEBHOOK`

* `username` - Name of the Atlas user to which to send notifications. Only a user in the project that owns the alert configuration is allowed here. Required for the `USER` notifications type.
* `victor_ops_api_key` - VictorOps API key. Required for the `VICTOR_OPS` notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
* `victor_ops_routing_key` - VictorOps routing key. Optional for the `VICTOR_OPS` notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
* `Roles` - Roles. Required. The following roles grant privileges within a project.
  Accepted values are:
    - `GROUP_CHARTS_ADMIN`
    - `GROUP_CLUSTER_MANAGER`
    - `GROUP_DATA_ACCESS_ADMIN`
    - `GROUP_DATA_ACCESS_READ_ONLY`
    - `GROUP_DATA_ACCESS_READ_WRITE`
    - `GROUP_OWNER`
    - `GROUP_READ_ONLY`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier used for terraform for internal manages and can be used to import.
* `alert_configuration_id` - Unique identifier for the alert configuration.
* `group_id` - Unique identifier of the project that owns this alert configuration.
* `created` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was created.
* `updated` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was last updated.

## Import

Alert Configuration can be imported using the `project_id-alert_configuration_id`, e.g.

```
$ terraform import mongodbatlas_alert_configuration.test 5d0f1f74cf09a29120e123cd-5d0f1f74cf09a29120e1fscg
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/alert-configurations/)
