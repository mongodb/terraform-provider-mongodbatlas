---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: alert_configuration"
sidebar_current: "docs-mongodbatlas-datasource-alert-configuration"
description: |-
    Describes a Alert Configuration.
---

# mongodbatlas_alert_configuration

`mongodbatlas_alert_configuration` describes an Alert Configuration.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

```terraform
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
  }

  matcher {
    field_name = "HOSTNAME_AND_PORT"
    operator   = "EQUALS"
    value      = "SECONDARY"
  }

  metric_threshold_config {
    metric_name = "ASSERT_REGULAR"
    operator    = "LESS_THAN"
    threshold   = 99.0
    units       = "RAW"
    mode        = "AVERAGE"
  }
}


data "mongodbatlas_alert_configuration" "test" {
	project_id             = mongodbatlas_alert_configuration.test.project_id
	alert_configuration_id = mongodbatlas_alert_configuration.test.alert_configuration_id
}
```

-> **NOTE:** In order to allow for a fast pace of change to alert variables some validations have been removed from this resource in order to unblock alert creation. Impacted areas have links to the MongoDB Atlas API documentation so always check it for the most current information: https://docs.atlas.mongodb.com/reference/api/alert-configurations-create-config/

```terraform
resource "mongodbatlas_alert_configuration" "test" {
  project_id = "<PROJECT-ID>"
  event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
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

  threshold_config {
    operator    = "LESS_THAN"
    threshold   = 1
    units       = "HOURS"
  }
}

data "mongodbatlas_alert_configuration" "test" {
	project_id             = mongodbatlas_alert_configuration.test.project_id
	alert_configuration_id = mongodbatlas_alert_configuration.test.alert_configuration_id
}
```

## Argument Reference

* `project_id` - (Required) The ID of the project where the alert configuration will create.
* `alert_configuration_id` - (Required) Unique identifier for the alert configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_id` - Unique identifier of the project that owns this alert configuration.
* `created` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was created.
* `updated` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was last updated.
* `enabled` - If set to true, the alert configuration is enabled. If enabled is not exported it is set to false.
* `event_type` - The type of event that will trigger an alert.

  -> ***IMPORTANT:*** Event Type has many possible values. All current options at available at https://docs.atlas.mongodb.com/reference/api/alert-configurations-create-config/ Details for both conditional and metric based alerts can be found by selecting the tabs on the [alert config page](https://docs.atlas.mongodb.com/reference/api/alert-configurations-create-config/) and checking the latest eventTypeName options.

  -> **NOTE:** If `event_type` is set to OUTSIDE_METRIC_THRESHOLD, the metricThreshold field must also be set.

### Matchers
Rules to apply when matching an object against this alert configuration. Only entities that match all these rules are checked for an alert condition. You can filter using the matchers array only when the eventTypeName specifies an event for a host, replica set, or sharded cluster.

* `field_name` - Name of the field in the target object to match on.

| Host alerts         | Replica set alerts  |  Sharded cluster alerts |
|:--------------      |:-------------       |:------                 |
| `TYPE_NAME`         | `REPLICA_SET_NAME`  | `CLUSTER_NAME`          |
| `HOSTNAME`          | `SHARD_NAME`        | `SHARD_NAME`            |
| `PORT`              | `CLUSTER_NAME`      |                         |
| `HOSTNAME_AND_PORT` |                     |                         |
| `REPLICA_SET_NAME`  |                     |                         |


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

### Metric Threshold Config
The threshold that causes an alert to be triggered. Required if `event_type_name` : "OUTSIDE_METRIC_THRESHOLD".

* `metric_name` - Name of the metric to check. The full list of current options is available [here](https://docs.atlas.mongodb.com/reference/alert-host-metrics/#measurement-types)

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

### Threshold Config
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
* `team_name` - Label for the team that receives this notification.
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

* `Roles` - The following roles grant privileges within a project.

    | Project roles                   | Organization roles  |
    |:----------                      |:-----------         |
    | `GROUP_CHARTS_ADMIN`            | `ORG_OWNER`         |
    | `GROUP_CLUSTER_MANAGER`         | `ORG_MEMBER`        |
    | `GROUP_DATA_ACCESS_ADMIN`       | `ORG_GROUP_CREATOR` |
    | `GROUP_DATA_ACCESS_READ_ONLY`   | `ORG_BILLING_ADMIN` |
    | `GROUP_DATA_ACCESS_READ_WRITE`  | `ORG_READ_ONLY`     |
    | `GROUP_OWNER`                   |                     |
    | `GROUP_READ_ONLY`               |                     |

See detailed information for arguments and attributes: [MongoDB API Alert Configuration](https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-config/)