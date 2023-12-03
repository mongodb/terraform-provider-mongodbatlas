---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_shared_tier_restore_job"
sidebar_current: "docs-mongodbatlas-datasource-mongodbatlas-shared-tier-restore-jobs"
description: |-
    Provides a Cloud Backup Shared Tier Snapshot Restore Job Datasource.
---


# Data Source: mongodbatlas_shared_tier_restore_job

`mongodbatlas_shared_tier_restore_job` provides a Cloud Backup Snapshot Restore Job data source for Shared Tier Clusters. Gets the cloud backup snapshot restore jobs for the specified shared tier cluster.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.
-> **NOTE:** This data source is only for Shared Tier Clusters (M2 and M5). See [here](https://www.mongodb.com/docs/atlas/reference/free-shared-limitations/) for more details on Shared Tier Cluster Limitations. 

## Example Usage
```terraform
data "mongodbatlas_shared_tier_restore_job" "test" {
  project_id   = "5d0f1f73cf09a29120e173cf"
  cluster_name = "MyClusterTest"
  job_id       = "5d1285acd5ec13b6c2d1726a"
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `job_id` - (Required) Unique 24-hexadecimal digit string that identifies the restore job to return.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` -	Indicates whether the restore job was canceled.
* `target_project_id` -	UTC ISO 8601 formatted point in time when Atlas created the restore job.
* `target_deployment_item_name` - Type of restore job to create. Possible values are: automated and download.
* `snapshot_url` -	Internet address from which you can download the compressed snapshot files. The resource returns this parameter when `deliveryType: DOWNLOAD`.
* `snapshot_id` -	Unique 24-hexadecimal digit string that identifies the snapshot to restore.
* `snapshot_finished_date` -	Date and time when MongoDB Cloud completed writing this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `restore_scheduled_date` -	Date and time when MongoDB Cloud will restore this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `restore_finished_date` -	Date and time when MongoDB Cloud completed writing this snapshot. MongoDB Cloud changes the status of the restore job to `CLOSED`. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `delivery_type` -	Means by which this resource returns the snapshot to the requesting MongoDB Cloud user. Values: `RESTORE`, `DOWNLOAD`.
* `expiration_date` -	Date and time when the download link no longer works. This parameter expresses its value in the ISO 8601 timestamp format in UTC.

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Shared-Tier-Restore-Jobs/operation/getSharedClusterBackupRestoreJob)