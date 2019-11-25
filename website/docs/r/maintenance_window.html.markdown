---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: maintenance_window"
sidebar_current: "docs-mongodbatlas-resource-maintenance_window"
description: |-
    Provides an Maintenance Window resource.
---

# mongodbatlas_maintenance_window

`mongodbatlas_maintenance_window` provides a resource to take a scheduled maintenance event for a project up to two times. 
Deferred maintenance events occur during your preferred maintenance window exactly one week after the previously scheduled date and time.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
  resource "mongodbatlas_maintenance_window" "test" {
    project_id  = "<your-project-id>"
    day_of_week = 3
    hour_of_day = 4
  }

```
## Example Usage

```hcl
  resource "mongodbatlas_maintenance_window" "test" {
    project_id  = "<your-project-id>"
    start_asap  = true 
  }
```

## Argument Reference

* `project_id` - The unique identifier of the project for the Maintenance Window.
* `day_of_week` - Day of the week when you would like the maintenance window to start as a 1-based integer: S=1, M=2, T=3, W=4, T=5, F=6, S=7.
* `hour_of_day` - Hour of the day when you would like the maintenance window to start. This parameter uses the 24-hour clock, where midnight is 0, noon is 12.
* `start_asap` - Flag indicating whether project maintenance has been directed to start immediately. If you request that maintenance begin immediately, this field returns true from the time the request was made until the time the maintenance event completes.
* `number_of_deferrals` - Number of times the current maintenance event for this project has been deferred.
* `defer` - Defer maintenance for the given project for one week.

## Import

Maintenance Window entries can be imported using project project_id, in the format `PROJECTID`, e.g.

```
$ terraform import mongodbatlas_maintenance_window.test 5d0f1f73cf09a29120e173cf
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/maintenance-windows/)