---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_schedule"
sidebar_current: "docs-mongodbatlas-datasource-cloud-backup-schedule"
description: |-
    Provides a Cloud Backup Schedule Datasource.
---

# Data Source: mongodbatlas_cloud_backup_schedule

`mongodbatlas_cloud_backup_schedule` provides a Cloud Backup Schedule datasource. An Atlas Cloud Backup Schedule provides the current cloud backup schedule for the cluster. 

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
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
* `reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the `reference_hour_of_day` that Atlas takes the snapshot.
* `restore_window_days` - Specifies a restore window in days for cloud backup to maintain.
* `id_policy` - Unique identifier of the backup policy.
* `policy_item_hourly` - Hourly policy item
* `policy_item_daily` - Daily policy item
* `policy_item_weekly` - Weekly policy item
* `policy_item_monthly` - Monthly policy item
* `auto_export_enabled` - Flag that indicates whether automatic export of cloud backup snapshots to the AWS bucket is enabled. Value can be one of the following:

    true - enables automatic export of cloud backup snapshots to the AWS bucket
    false - disables automatic export of cloud backup snapshots to the AWS bucket (default)
* `use_org_and_group_names_in_export_prefix` - Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots. To learn more about the metadata files that Atlas uploads, see [Export Cloud Backup Snapshot](https://www.mongodb.com/docs/atlas/backup/cloud-backup/export/#std-label-cloud-provider-snapshot-export).
### Export
* `export_bucket_id` - Unique identifier of the mongodbatlas_cloud_backup_snapshot_export_bucket export_bucket_id value.
* `frequency_type` - Frequency associated with the export snapshot item.
### Policy Item Hourly
*
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - Scope of the backup policy item: days, weeks, or months.
* `retention_value` - Value to associate with `retention_unit`.

### Policy Item Daily
*
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - Scope of the backup policy item: days, weeks, or months.
* `retention_value` - Value to associate with `retention_unit`.

### Policy Item Weekly
*
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - Scope of the backup policy item: days, weeks, or months.
* `retention_value` - Value to associate with `retention_unit`.

### Policy Item Monthly
*
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - Scope of the backup policy item: days, weeks, or months.
* `retention_value` - Value to associate with `retention_unit`.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/get-all-schedules/)