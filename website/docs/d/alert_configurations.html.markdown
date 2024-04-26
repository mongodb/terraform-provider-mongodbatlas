---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: Alert Configurations"
sidebar_current: "docs-mongodbatlas-datasource-alert-configurations"
description: |-
    Describe all Alert Configurations in Project.
---

# Data Source: mongodbatlas_alert_configurations

`mongodbatlas_alert_configurations` describes all Alert Configurations by the provided project_id. The data source requires your Project ID.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```terraform
data "mongodbatlas_alert_configurations" "import" {
  project_id = var.project_id

  output_type = ["resource_hcl", "resource_import"]
}

locals {
  alerts = data.mongodbatlas_alert_configurations.import.results

  outputs = flatten([
    for i, alert in local.alerts :
    alert.output == null ? [] : alert.output
  ])

  output_values = compact([for i, o in local.outputs : o.value])
}

output "alert_output" {
  value = join("\n", local.output_values)
}
```

Refer to the following for a full example on using this data_source as a tool to import all resources:
* [atlas-alert-configurations](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-alert-configurations)

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get the alert configurations.
* `list_options` - (Optional) Arguments that dictate how many and which results are returned by the data source
* `list_options.page_num` - Which page of results to retrieve (default to first page)
* `list_options.items_per_page` - How many alerts to retrieve per page (default 100)
* `list_options.include_count` - Whether to include total count of results in the response (default false)
* `output_type` - (Optional) List of requested string formatted output to be included on each individual result. Options are `resource_hcl` and `resource_import`. Available to make it easy to gather resource statements for existing alert configurations, and corresponding import statements to import said resource state into the statefile.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total count of results
* `results` - A list of alert configurations for the project_id, constrained by the `list_options`.

### Alert Configuration

* `project_id` - The ID of the project where the alert configuration exists
* `alert_configuration_id` - The ID of the alert configuration
* `created` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was created.
* `updated` - Timestamp in ISO 8601 date and time format in UTC when this alert configuration was last updated.
* `enabled` - If set to true, the alert configuration is enabled. If enabled is not exported it is set to false.
* `event_type` - The type of event that will trigger an alert.
* `matcher` - Rules to apply when matching an object against this alert configuration. See [matchers](#matchers).
* `metric_threshold_config` - The threshold that causes an alert to be triggered. Required if `event_type_name` : `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`. See [metric threshold config](#metric-threshold-config).
* `threshold_config` - 	 Threshold that triggers an alert. Required if `event_type_name` is any value other than `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`. See [threshold config](#threshold-config).
* `notifications` - List of notifications to send when an alert condition is detected. See [notifications](#notifications).
* `output` - Requested output string format for the alert configuration

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

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/alert-configurations/)
Or refer to the individual resource or data_source documentation on alert configuration.
