---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: alert_configuration"
sidebar_current: "docs-mongodbatlas-resource-alert-configuration"
description: |-
    Provides an Alert Configuration resource.
---

# Resource: mongodbatlas_alert_configuration

`mongodbatlas_alert_configuration` provides an Alert Configuration resource to define the conditions that trigger an alert and the methods of notification within a MongoDB Atlas project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

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
    roles         = ["GROUP_CLUSTER_MANAGER"]
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
    field_name = "CLUSTER_NAME"
    operator   = "EQUALS"
    value      = "my-cluster"
  }

  threshold_config {
    operator    = "LESS_THAN"
    threshold   = 1
    units       = "HOURS"
  }
}
```

### Create an alert with two notifications using Email and SMS


```terraform
resource "mongodbatlas_alert_configuration" "test" {
  project_id = "PROJECT ID"
  event_type = "OUTSIDE_METRIC_THRESHOLD"
  enabled    = true

  notification {
    type_name     = "GROUP"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = false
    email_enabled = true
    roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER", "GROUP_DATA_ACCESS_ADMIN"]
  }

  notification {
    type_name     = "ORG"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = true
    email_enabled = false
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
```

### Create third party notification using credentials from existing third party integration


```terraform
data "mongodbatlas_third_party_integration" "test" {
    project_id = "PROJECT ID"
    type = "PAGER_DUTY"
}

resource "mongodbatlas_alert_configuration" "test" {
  project_id = "PROJECT ID"
  enabled    = true
  event_type = "USERS_WITHOUT_MULTI_FACTOR_AUTH"

  notification {
    type_name     = "PAGER_DUTY"
    integration_id = data.mongodbatlas_third_party_integration.test.id
  }
}
```

## Argument Reference

* `project_id` - (Required) The ID of the project where the alert configuration will create.
* `enabled` - It is not required, but If the attribute is omitted, by default will be false, and the configuration would be disabled. You must set true to enable the configuration.
* `event_type` - (Required) The type of event that will trigger an alert.

  -> ***IMPORTANT:*** Event Type has many possible values. Details for both conditional and metric based alerts can be found by selecting the tabs on the [alert config page](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Alert-Configurations/operation/createAlertConfiguration) and checking the latest eventTypeName options.


  -> **NOTE:** If `event_type` is set to `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`, the `metric_threshold_config` field must also be configured.

### Matchers
Rules to apply when matching an object against this alert configuration. Only entities that match all these rules are checked for an alert condition. You can filter using the matchers array only when the eventTypeName specifies an event for a host, replica set, or sharded cluster.

* `field_name` - (Required) Name of the field in the target object to match on.

| Host alerts         | Replica set alerts  |  Sharded cluster alerts |
|:----------           |:-------------       |:------                 |
| `TYPE_NAME`         | `REPLICA_SET_NAME`  | `CLUSTER_NAME`          |
| `HOSTNAME`          | `SHARD_NAME`        | `SHARD_NAME`            |
| `PORT`              | `CLUSTER_NAME`      |                         |
| `HOSTNAME_AND_PORT` |                     |                         |
| `REPLICA_SET_NAME`  |                     |                         |



All other types of alerts do not support matchers.

* `operator` - (Required) The operator to test the field’s value.
  Accepted values are:
    - `EQUALS`
    - `NOT_EQUALS`
    - `CONTAINS`
    - `NOT_CONTAINS`
    - `STARTS_WITH`
    - `ENDS_WITH`
    - `REGEX`

* `value` - (Required) Value to test with the specified operator. If `field_name` is set to TYPE_NAME, you can match on the following values:
    - `PRIMARY`
    - `SECONDARY`
    - `STANDALONE`
    - `CONFIG`
    - `MONGOS`

### Metric Threshold Config (`metric_threshold_config`)
The threshold that causes an alert to be triggered. Required if `event_type_name` : `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`

* `metric_name` - (Required) Name of the metric to check. The full list being quite large, please refer to atlas docs [here for general metrics](https://docs.atlas.mongodb.com/reference/alert-host-metrics/#measurement-types) and [here for serverless metrics](https://www.mongodb.com/docs/atlas/reference/api/alert-configurations-create-config/#serverless-measurements)
* `operator` - The operator to apply when checking the current metric value against the threshold value.
  Accepted values are:
    - `GREATER_THAN`
    - `LESS_THAN`

* `threshold` - Threshold value outside of which an alert will be triggered.
* `units` - The units for the threshold value. Depends on the type of metric.
  Refer to the [MongoDB API Alert Configuration documentation](https://www.mongodb.com/docs/atlas/reference/api/alert-configurations-get-config/#request-body-parameters) for a list of accepted values.

* `mode` - This must be set to AVERAGE. Atlas computes the current metric value as an average.

### Threshold Config (`threshold_config`)
* `operator` - The operator to apply when checking the current metric value against the threshold value.
  Accepted values are:
    - `GREATER_THAN`
    - `LESS_THAN`

* `threshold` - Threshold value outside of which an alert will be triggered.
* `units` - The units for the threshold value. Depends on the type of metric.
  Refer to the [MongoDB API Alert Configuration documentation](https://www.mongodb.com/docs/atlas/reference/api/alert-configurations-get-config/#request-body-parameters) for a list of accepted values.

### Notifications
List of notifications to send when an alert condition is detected.

* `api_token` - Slack API token. Required for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
* `channel_name` - Slack channel name. Required for the SLACK notifications type.
* `datadog_api_key` - Datadog API Key. Found in the Datadog dashboard. Required for the DATADOG notifications type.
* `datadog_region` - Region that indicates which API URL to use. See the `datadogRegion` field in the `notifications` request parameter of [MongoDB API Alert Configuration documentation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Alert-Configurations/operation/createAlertConfiguration) for more details. The default Datadog region is US.
* `delay_min` - Number of minutes to wait after an alert condition is detected before sending out the first notification.
* `email_address` - Email address to which alert notifications are sent. Required for the EMAIL notifications type.
* `email_enabled` - Flag indicating email notifications should be sent. This flag is only valid if `type_name` is set to `ORG`, `GROUP`, or `USER`.
* `interval_min` - Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved. The minimum value is 5. **NOTE** `PAGER_DUTY`, `VICTOR_OPS`, and `OPS_GENIE` notifications do not return this value. The notification interval must be configured and managed within each external service.
* `mobile_number` - Mobile number to which alert notifications are sent. Required for the SMS notifications type.
* `ops_genie_api_key` - Opsgenie API Key. Required for the `OPS_GENIE` notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
* `ops_genie_region` - Region that indicates which API URL to use. Accepted regions are: `US` ,`EU`. The default Opsgenie region is US.
* `service_key` - PagerDuty service key. Required for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
* `sms_enabled` - Flag indicating if text message notifications should be sent to this user's mobile phone. This flag is only valid if `type_name` is set to `ORG`, `GROUP`, or `USER`.
* `team_id` - Unique identifier of a team.
* `team_name` - Label for the team that receives this notification.
* `type_name` - (Required) Type of alert notification.
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
* `webhook_secret` - Optional authentication secret for the `WEBHOOK` notifications type.
* `webhook_url` - Target URL  for the `WEBHOOK` notifications type.
* `microsoft_teams_webhook_url` - Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. Required if `type_name` is `MICROSOFT_TEAMS`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.
* `roles` - Optional. One or more roles that receive the configured alert. If you include this field, Atlas sends alerts only to users assigned the roles you specify in the array. If you omit this field, Atlas sends alerts to users assigned any role. This parameter is only valid if `type_name` is set to `ORG`, `GROUP`, or `USER`.
  Accepted values are:

    | Project roles                   | Organization roles  |
    |:----------                      |:-----------         |
    | `GROUP_CLUSTER_MANAGER`         | `ORG_OWNER`         |
    | `GROUP_DATA_ACCESS_ADMIN`       | `ORG_MEMBER`        |
    | `GROUP_DATA_ACCESS_READ_ONLY`   | `ORG_GROUP_CREATOR` |
    | `GROUP_DATA_ACCESS_READ_WRITE`  | `ORG_BILLING_ADMIN` |
    | `GROUP_OWNER`                   | `ORG_READ_ONLY`     |
    | `GROUP_READ_ONLY`               |                     |

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
terraform import mongodbatlas_alert_configuration.test 5d0f1f74cf09a29120e123cd-5d0f1f74cf09a29120e1fscg
```

**NOTE**: Third-party notifications will not contain their respective credentials as these are sensitive attributes. If you wish to perform updates on these notifications without providing the original credentials, the corresponding `notifier_id` attribute must be provided instead.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/alert-configurations/)
