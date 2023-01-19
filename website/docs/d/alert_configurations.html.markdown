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
* `matcher` - Rules to apply when matching an object against this alert configuration
* `metric_threshold_config` - The threshold that causes an alert to be triggered. Required if `event_type_name` : `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`
* `threshold_config` - 	 Threshold that triggers an alert. Required if `event_type_name` is any value other than `OUTSIDE_METRIC_THRESHOLD` or `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`.
* `notifications` - List of notifications to send when an alert condition is detected.
* `output` - Requested output string format for the alert configuration

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/alert-configurations/)
Or refer to the individual resource or data_source documentation on alert configuration.