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
* `policy_item_yearly` - Yearly policy item
* `auto_export_enabled` - Flag that indicates whether automatic export of cloud backup snapshots to the AWS bucket is enabled. Value can be one of the following:
    * true - enables automatic export of cloud backup snapshots to the AWS bucket
    * false - disables automatic export of cloud backup snapshots to the AWS bucket (default)
* `use_org_and_group_names_in_export_prefix` - Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots. To learn more about the metadata files that Atlas uploads, see [Export Cloud Backup Snapshot](https://www.mongodb.com/docs/atlas/backup/cloud-backup/export/#std-label-cloud-provider-snapshot-export).
### Export
* `export_bucket_id` - Unique identifier of the mongodbatlas_cloud_backup_snapshot_export_bucket export_bucket_id value.
* `frequency_type` - Frequency associated with the export snapshot item.

### Policy Item Hourly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For hourly policies, the frequency type is defined as `hourly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (hourly in this case). The supported values for hourly policies are `1`, `2`, `4`, `6`, `8` or `12` hours. Note that `12` hours is the only accepted value for NVMe clusters.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`.

### Policy Item Daily
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For daily policies, the frequency type is defined as `daily`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (daily in this case). The only supported value for daily policies is `1` day.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`.  Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the hourly policy item specifies a retention of two days, the daily retention policy must specify two days or greater.

### Policy Item Weekly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For weekly policies, the frequency type is defined as `weekly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (weekly in this case). The supported values for weekly policies are `1` through `7`, where `1` represents Monday and `7` represents Sunday.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Weekly policy must have retention of at least 7 days or 1 week. Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the daily policy item specifies a retention of two weeks, the weekly retention policy must specify two weeks or greater.

### Policy Item Monthly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For monthly policies, the frequency type is defined as `monthly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (monthly in this case). The supported values for weekly policies are 
  * `1` through `28` where the number represents the day of the month i.e. `1` is the first of the month and `5` is the fifth day of the month.
  * `40` represents the last day of the month (depending on the month).
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Monthly policy must have retention days of at least 31 days or 5 weeks or 1 month. Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the weekly policy item specifies a retention of two weeks, the montly retention policy must specify two weeks or greater.

### Policy Item Yearly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For yearly policies, the frequency type is defined as `yearly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (yearly in this case). The supported values for yearly policies are 
  * `1` through `12` the first day of the month where the number represents the month, i.e. `1` is January and `12` is December.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Yearly policy must have retention of at least 1 year.

### Snapshot Distribution
* `cloud_provider` - Human-readable label that identifies the cloud provider that stores the snapshot copy. i.e. "AWS" "AZURE" "GCP"
* `frequencies` - List that describes which types of snapshots to copy. i.e. "HOURLY" "DAILY" "WEEKLY" "MONTHLY" "YEARLY" "ON_DEMAND"
* `region_name` - Target region to copy snapshots belonging to replicationSpecId to. Please supply the 'Atlas Region' which can be found under https://www.mongodb.com/docs/atlas/reference/cloud-providers/ 'regions' link
* `replication_spec_id` - Unique 24-hexadecimal digit string that identifies the replication object for a zone in a cluster. For global clusters, there can be multiple zones to choose from. For sharded clusters and replica set clusters, there is only one zone in the cluster. To find the Replication Spec Id, consult the replicationSpecs array returned from [Return One Multi-Cloud Cluster in One Project](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/getCluster).
* `should_copy_oplogs` - Flag that indicates whether to copy the oplogs to the target region. You can use the oplogs to perform point-in-time restores.

**Note** The parameter deleteCopiedBackups is not supported in terraform please leverage Atlas Admin API or AtlasCLI instead to manage the lifecycle of backup snaphot copies.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/get-all-schedules/)
