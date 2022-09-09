---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_schedule"
sidebar_current: "docs-mongodbatlas-resource-cloud-backup-schedule"
description: |-
    Provides a Cloud Backup Schedule resource.
---

# Resource: mongodbatlas_cloud_backup_schedule

`mongodbatlas_cloud_backup_schedule` provides a cloud backup schedule resource. The resource lets you create, read, update and delete a cloud backup schedule.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **API Key Access List**: This resource requires an Atlas API Access Key List to utilize this feature. This means to manage this resources you must have the IP address or CIDR block that the Terraform connection is coming from added to the Atlas API Key Access List of the Atlas API key you are using. See [Resources that require API Key List](https://www.mongodb.com/docs/atlas/configure-api-access/#use-api-resources-that-require-an-access-list) for details.

In the Terraform MongoDB Atlas Provider 1.0.0 we have re-architected the way in which Cloud Backup Policies are manged with Terraform to significantly reduce the complexity. Due to this change we've provided multiple examples below to help express how this new resource functions.


## Example Usage - Create a Cluster with 2 Policies Items

You can create a new cluster with `cloud_backup` enabled and then immediately overwrite the default cloud backup policy that Atlas creates by default at the same time with this example.

```terraform
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // must be enabled in order to use cloud_backup_schedule resource
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4


  // This will now add the desired policy items to the existing mongodbatlas_cloud_backup_schedule resource
  policy_item_hourly {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 1
  }
  policy_item_daily {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 2
  }
}
```

## Example Usage - Create a Cluster with Cloud Backup Enabled but No Policy Items

You can enable `cloud_backup` in the Cluster resource and then use the `cloud_backup_schedule` resource with no policy items to remove the default policy that Atlas creates when you enable Cloud Backup. This allows you to then create a policy when you are ready to via Terraform.

```terraform
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // must be enabled in order to use cloud_backup_schedule resource
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4

}
```

## Example Usage - Add 4 Policies Items To A Cluster With Cloud Backup Previously Enabled but with No Policy Items

If you followed the example to Create a Cluster with Cloud Backup Enabled but No Policy Items and then want to add policy items later to the `mongodbatlas_cloud_backup_schedule` this example shows how.

The cluster already exists with `cloud_backup` enabled
```terraform
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // must be enabled in order to use cloud_backup_schedule resource
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4
  
  // This will now add the desired policy items to the existing mongodbatlas_cloud_backup_schedule resource
  policy_item_hourly {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 1
  }
  policy_item_daily {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 2
  }
  policy_item_weekly {
    frequency_interval = 4
    retention_unit     = "weeks"
    retention_value    = 3
  }
  policy_item_monthly {
    frequency_interval = 5
    retention_unit     = "months"
    retention_value    = 4
  }

}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshot backup policy you want to retrieve.
* `reference_hour_of_day` - (Optional) UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
* `reference_minute_of_hour` - (Optional) UTC Minutes after `reference_hour_of_day` that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.
* `restore_window_days` - (Optional) Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
* `update_snapshots` - (Optional) Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.
* `policy_item_hourly` - (Optional) Hourly policy item
* `policy_item_daily` - (Optional) Daily policy item
* `policy_item_weekly` - (Optional) Weekly policy item
* `policy_item_monthly` - (Optional) Monthly policy item
* `auto_export_enabled` - Flag that indicates whether automatic export of cloud backup snapshots to the AWS bucket is enabled. Value can be one of the following:

    true - enables automatic export of cloud backup snapshots to the AWS bucket
    false - disables automatic export of cloud backup snapshots to the AWS bucket (default)
* `use_org_and_group_names_in_export_prefix` - Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots. To learn more about the metadata files that Atlas uploads, see [Export Cloud Backup Snapshot](https://www.mongodb.com/docs/atlas/backup/cloud-backup/export/#std-label-cloud-provider-snapshot-export).
### Export
* `export_bucket_id` - Unique identifier of the mongodbatlas_cloud_backup_snapshot_export_bucket export_bucket_id value.
* `frequency_type` - Frequency associated with the export snapshot item.

### Policy Item Hourly
* 
* `frequency_interval` - (Required) Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - (Required) Scope of the backup policy item: days, weeks, or months.
* `retention_value` - (Required) Value to associate with `retention_unit`.

### Policy Item Daily
*
* `frequency_interval` - (Required) Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - (Required) Scope of the backup policy item: days, weeks, or months.
* `retention_value` - (Required) Value to associate with `retention_unit`.

### Policy Item Weekly
*
* `frequency_interval` - (Required) Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - (Required) Scope of the backup policy item: days, weeks, or months.
* `retention_value` - (Required) Value to associate with `retention_unit`.

### Policy Item Monthly
*
* `frequency_interval` - (Required) Desired frequency of the new backup policy item specified by `frequency_type`.
* `retention_unit` - (Required) Scope of the backup policy item: days, weeks, or months.
* `retention_value` - (Required) Value to associate with `retention_unit`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - Unique identifier of the Atlas cluster.
* `next_snapshot` - Timestamp in the number of seconds that have elapsed since the UNIX epoch when Atlas takes the next snapshot.
* `id_policy` - Unique identifier of the backup policy.

## Import

Cloud Backup Schedule entries can be imported using project_id and cluster_name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cloud_backup_schedule.test 5d0f1f73cf09a29120e173cf-MyClusterTest
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/)