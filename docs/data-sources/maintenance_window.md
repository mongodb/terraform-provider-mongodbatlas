---
subcategory: "Projects"
---

# Data Source: mongodbatlas_maintenance_window

`mongodbatlas_maintenance_window` provides a Maintenance Window entry datasource. Gets information regarding the configured maintenance window for a MongoDB Atlas project.

-> **NOTE:** Maintenance window times use the project's configured timezone. To change the timezone, update the Project Time Zone setting in the Atlas Project Settings.

## Examples Usage

```terraform
resource "mongodbatlas_maintenance_window" "test" {
  project_id  = "<your-project-id>"
  day_of_week = 3
  hour_of_day = 4
  auto_defer_once_enabled = true
}

data "mongodbatlas_maintenance_window" "test" {
  project_id = mongodbatlas_maintenance_window.test.id
}
```


```terraform
resource "mongodbatlas_maintenance_window" "test" {
  project_id  = "<your-project-id>"
  start_asap  = true
}

data "mongodbatlas_maintenance_window" "test" {
  project_id = mongodbatlas_maintenance_window.test.id
}
```

## Argument Reference

* `project_id` - The unique identifier of the project for the Maintenance Window, also known as `groupId` in the official documentation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `day_of_week` - Day of the week when you would like the maintenance window to start as a 1-based integer: Su=1, M=2, T=3, W=4, T=5, F=6, Sa=7.
* `hour_of_day` - Hour of the day when you would like the maintenance window to start. This parameter uses the 24-hour clock, where midnight is 0, noon is 12. Uses the project's configured timezone.
* `start_asap` - Flag indicating whether project maintenance has been directed to start immediately. If requested, this field returns true from the time the request was made until the time the maintenance event completes.
* `number_of_deferrals` - Number of times the current maintenance event for this project has been deferred, there can be a maximum of 2 deferrals.
* `auto_defer_once_enabled` - When `true`, enables automatic deferral of all scheduled maintenance for the given project by one week.
* `protected_hours` - (Optional) Defines the time period during which there will be no standard updates to the clusters. See [Protected Hours](#protected-hours).
* `time_zone_id` - Identifier for the current time zone of the maintenance window. This can only be updated via the Project Settings UI.
* `wave_assignment` - Maintenance wave explicitly assigned to this project. Always returned when a value has been set, regardless of the organization's [`wave_assignment_mode`](../data-sources/org_maintenance_settings.md#attributes-reference). When the mode is `ENV_TAG_MAPPING`, the system preserves the stored value but does not uses it for scheduling. Switching back to `MANUAL` restores this value as the effective wave. Returns `0` when no explicit wave has been assigned.
* `effective_wave_assignment` - Read-only maintenance wave Atlas uses when scheduling maintenance for this project. This value can differ from `wave_assignment` in the following scenarios:
  - **`ENV_TAG_MAPPING` mode is active at the organization level.** When the organization's `wave_assignment_mode` is set to `ENV_TAG_MAPPING` (see [`mongodbatlas_org_maintenance_settings`](../resources/org_maintenance_settings.md)), Atlas ignores any explicit `wave_assignment` and derives the effective wave from the project's environment tag. A project can have `wave_assignment = 1` in state while `effective_wave_assignment` returns a different value.
  - **Cross-organization billing (`MAINTENANCE_SEQUENCE_CROSS_ORG`).** When a linked non-paying organization inherits the paying organization's wave assignment mode. If the paying organization switches to `ENV_TAG_MAPPING`, all linked projects follow regardless of any explicit `wave_assignment` set on them.

### Protected Hours
* `start_hour_of_day` - Zero-based integer that represents the beginning hour of the day for the protected hours window.
* `end_hour_of_day` - Zero-based integer that represents the end hour of the day for the protected hours window.

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-maintenance-windows)
