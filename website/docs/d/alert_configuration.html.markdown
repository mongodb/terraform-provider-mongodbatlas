---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: alert_configuration"
sidebar_current: "docs-mongodbatlas-datasource-alert-configuration"
description: |-
    Describes a Alert Configuration.
---

# Data Source: mongodbatlas_alert_configuration

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
    roles         = ["GROUP_CLUSTER_MANAGER"]
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

Utilize data_source to generate resource hcl and import statement. Useful if you have a specific alert_configuration_id and are looking to manage it as is in state. To import all alerts, refer to the documentation on [data_source_mongodbatlas_alert_configurations](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/alert_configurations)
```
data "mongodbatlas_alert_configuration" "test" {
    project_id             = var.project_id
    alert_configuration_id = var.alert_configuration_id

    output {
        type = "resource_hcl"
        label = "test"
    }

    output {
        type = "resource_import"
        label = "test"
    }
}
```

## Argument Reference

* `project_id` - (Required) The ID of the project where the alert configuration will create.
* `alert_configuration_id` - (Required) Unique identifier for the alert configuration.
* `output` - (Optional) List of formatted output requested for this alert configuration
* `output.#.type` - (Required) If the output is requested, you must specify its type. The format is computed as `output.#.value`, the following are the supported types:
- `resource_hcl`: This string is used to define the resource as it exists in MongoDB Atlas.
- `resource_import`: This string is used to import the existing resource into the state file.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `created` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was created.
* `updated` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was last updated.
* `enabled` - If set to true, the alert configuration is enabled. If enabled is not exported it is set to false.
* `event_type` - The type of event that will trigger an alert.
* `matcher` - Rules to apply when matching an object against this alert configuration. See [matchers](#matchers).
* `metric_threshold_config` - The threshold that causes an alert to be triggered. Required if `event_type_name` : `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`. See [metric threshold config](#metric-threshold-config).
* `threshold_config` - 	 Threshold that triggers an alert. Required if `event_type_name` is any value other than `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`. See [threshold config](#threshold-config).
* `notifications` - List of notifications to send when an alert condition is detected. See [notifications](#notifications).

  -> ***IMPORTANT:*** Event Type has many possible values. Details for both conditional and metric based alerts can be found by selecting the tabs on the [alert config page](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Alert-Configurations/operation/createAlertConfiguration) and checking the latest eventTypeName options.

  -> **NOTE:** If `event_type` is set to `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`, the `metric_threshold_config` field must also be configured.

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
The threshold that causes an alert to be triggered. Required if `event_type_name` : `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`.

* `metric_name` - Name of the metric to check. The full list being quite large, please refer to atlas docs [here for general metrics](https://docs.atlas.mongodb.com/reference/alert-host-metrics/#measurement-types) and [here for serverless metrics](https://www.mongodb.com/docs/atlas/reference/api/alert-configurations-create-config/#serverless-measurements)

* `operator` - The operator to apply when checking the current metric value against the threshold value.
  Accepted values are:
    - `GREATER_THAN`
    - `LESS_THAN`

* `threshold` - Threshold value outside of which an alert will be triggered.
* `units` - The units for the threshold value. Depends on the type of metric.
  Refer to the [MongoDB API Alert Configuration documentation](https://www.mongodb.com/docs/atlas/reference/api/alert-configurations-get-config/#request-body-parameters) for a list of accepted values.
* `mode` - This must be set to AVERAGE. Atlas computes the current metric value as an average.

### Threshold Config
* `operator` - The operator to apply when checking the current metric value against the threshold value.
  Accepted values are:
    - `GREATER_THAN`
    - `LESS_THAN`

* `threshold` - Threshold value outside of which an alert will be triggered.
* `units` - The units for the threshold value. Depends on the type of metric.
  Refer to the [MongoDB API Alert Configuration documentation](https://www.mongodb.com/docs/atlas/reference/api/alert-configurations-get-config/#request-body-parameters) for a list of accepted values.

### Notifications
Notifications to send when an alert condition is detected.

* `api_token` - Slack API token. Required for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
* `channel_name` - Slack channel name. Required for the SLACK notifications type.
* `datadog_api_key` - Datadog API Key. Found in the Datadog dashboard. Required for the DATADOG notifications type.
* `datadog_region` - Region that indicates which API URL to use. See the `datadogRegion` field in the `notifications` request parameter of [MongoDB API Alert Configuration documentation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Alert-Configurations/operation/createAlertConfiguration) for more details. The default Datadog region is US.
* `delay_min` - Number of minutes to wait after an alert condition is detected before sending out the first notification.
* `email_address` - Email address to which alert notifications are sent. Required for the EMAIL notifications type.
* `email_enabled` - Flag indicating email notifications should be sent. Atlas returns this value if `type_name` is set  to `ORG`, `GROUP`, or `USER`.
* `interval_min` - Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved. The minimum value is 5.
* `mobile_number` - Mobile number to which alert notifications are sent. Required for the SMS notifications type.
* `ops_genie_api_key` - Opsgenie API Key. Required for the `OPS_GENIE` notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
* `ops_genie_region` - Region that indicates which API URL to use. Accepted regions are: `US` ,`EU`. The default Opsgenie region is US.
* `service_key` - PagerDuty service key. Required for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
* `sms_enabled` - Flag indicating text notifications should be sent. Atlas returns this value if `type_name` is set to `ORG`, `GROUP`, or `USER`.
* `team_id` - Unique identifier of a team.
* `team_name` - Label for the team that receives this notification.
* `type_name` - Type of alert notification.
  Accepted values are:
    - `DATADOG`
    - `EMAIL`
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
    - `MICROSOFT_TEAMS`

* `integration_id` - The ID of the associated integration, the credentials of which to use for requests.
* `notifier_id` - The notifier ID is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
* `username` - Name of the Atlas user to which to send notifications. Only a user in the project that owns the alert configuration is allowed here. Required for the `USER` notifications type.
* `victor_ops_api_key` - VictorOps API key. Required for the `VICTOR_OPS` notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
* `victor_ops_routing_key` - VictorOps routing key. Optional for the `VICTOR_OPS` notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
* `webhook_secret` - Authentication secret for the `WEBHOOK` notifications type.
* `webhook_url` - Target URL  for the `WEBHOOK` notifications type.
* `microsoft_teams_webhook_url` - Microsoft Teams channel incoming webhook URL. Required for the `MICROSOFT_TEAMS` notifications type.
* `roles` - Atlas role in current Project or Organization. Atlas returns this value if you set `type_name` to `ORG` or `GROUP`.

See detailed information for arguments and attributes: [MongoDB API Alert Configuration](https://docs.atlas.mongodb.com/reference/api/alert-configurations-get-config/)
