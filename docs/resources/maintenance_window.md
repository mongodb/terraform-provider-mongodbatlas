---
subcategory: "Maintenance Windows"
---

# Resource: mongodbatlas_maintenance_window

`mongodbatlas_maintenance_window` provides a resource to schedule the maintenance window for your MongoDB Atlas Project and/or set to defer a scheduled maintenance up to two times. Please refer to [Maintenance Windows](https://www.mongodb.com/docs/atlas/tutorial/cluster-maintenance-window/#configure-maintenance-window) documentation for more details.

-> **NOTE:** Only a single maintenance window resource can be defined per project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **NOTE:** Maintenance window times use the project's configured timezone. To change the timezone, update the Project Time Zone setting in the Atlas Project Settings.

## Maintenance Window Considerations:
- Urgent Maintenance Activities Cannot Wait: Urgent maintenance activities such as security patches cannot wait for your chosen window. Atlas will start those maintenance activities when needed.

Once maintenance is scheduled for your cluster, you cannot change your maintenance window until the current maintenance efforts have completed.
- Maintenance Requires Replica Set Elections: Atlas performs maintenance the same way as the manual maintenance procedure. This requires at least one replica set election during the maintenance window per replica set.
- Maintenance Starts As Close to the Hour As Possible: Maintenance always begins as close to the scheduled hour as possible, but in-progress cluster updates or expected system issues could delay the start time.


## Example Usage

```terraform
  resource "mongodbatlas_maintenance_window" "test" {
    project_id  = "<your-project-id>"
    day_of_week = 3
    hour_of_day = 4

    protected_hours {
    start_hour_of_day = 9
    end_hour_of_day   = 17
    }
  }

```

```terraform
  resource "mongodbatlas_maintenance_window" "test" {
    project_id = "<your-project-id>"
    defer      = true
  }
```

## Argument Reference

* `project_id` - The unique identifier of the project for the Maintenance Window.
* `day_of_week` - (Required) Day of the week when you would like the maintenance window to start as a 1-based integer: Su=1, M=2, T=3, W=4, T=5, F=6, Sa=7.
* `hour_of_day` - Hour of the day when you would like the maintenance window to start. This parameter uses the 24-hour clock, where midnight is 0, noon is 12. Uses the project's configured timezone. Defaults to 0.
* `defer` - Defer the next scheduled maintenance for the given project for one week.
* `auto_defer` - Defer any scheduled maintenance for the given project for one week.
* `auto_defer_once_enabled` - Flag that indicates whether you want to defer all maintenance windows one week they would be triggered.
* `protected_hours` - (Optional) Defines the time period during which there will be no standard updates to the clusters. See [Protected Hours](#protected-hours).

### Protected Hours
* `start_hour_of_day` - Zero-based integer that represents the beginning hour of the day for the protected hours window.
- `end_hour_of_day` - Zero-based integer that represents the end hour of the day for the protected hours window.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `number_of_deferrals` - Number of times the current maintenance event for this project has been deferred, there can be a maximum of 2 deferrals.
* `time_zone_id` - Identifier for the current time zone of the maintenance window. This can only be updated via the Project Settings UI.
* `start_asap` - Flag indicating whether project maintenance has been directed to start immediately. If requested, this field returns true from the time the request was made until the time the maintenance event completes.

-> **NOTE:** The `start_asap` attribute can only be enabled via API.

## Import

Maintenance Window entries can be imported using project project_id, in the format `PROJECTID`, e.g.

```
$ terraform import mongodbatlas_maintenance_window.test 5d0f1f73cf09a29120e173cf
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/maintenance-windows/)
