---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_schedule"
sidebar_current: "docs-mongodbatlas-datasource-cloud-backup-schedule"
description: |-
    Provides a Cloud Backup Schedule Datasource.
---

# mongodbatlas_cloud_backup_schedule

`mongodbatlas_cloud_backup_schedule` provides a Cloud Backup Schedule datasource. An Atlas Cloud Backup Schedule provides the current cloud backup schedule for the cluster. 

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // enable cloud backup snapshots
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4
}

data "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cloud_backup_schedule.test.project_id
  cluster_name = mongodbatlas_cloud_backup_schedule.test.cluster_name
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshots backup policy you want to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - Unique identifier of the Atlas cluster.
* `next_snapshot` - UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
* `reference_hour_of_day` - UTC Hour of day between 0 and 23 representing which hour of the day that Atlas takes a snapshot.
* `reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the referenceHourOfDay that Atlas takes the snapshot.
* `restore_window_days` - Specifies a restore window in days for cloud backup to maintain.

### Policies
* `policies` - A list of policy definitions for the cluster.
* `policies.0.id` - Unique identifier of the backup policy.

#### Policy Item
* `policies.0.policy_item` - A list of specifications for a policy.
* `policies.0.policy_item.#.id` - Unique identifier for this policy item.
* `policies.0.policy_item.#.frequency_interval` - The frequency interval for a set of snapshots.
* `policies.0.policy_item.#.frequency_type` - A type of frequency (hourly, daily, weekly, monthly).
* `policies.0.policy_item.#.retention_unit` - The unit of time in which snapshot retention is measured (days, weeks, months).
* `policies.0.policy_item.#.retention_value` - The number of days, weeks, or months the snapshot is retained.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/get-all-schedules/)