---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_snapshot_restore_jobs"
sidebar_current: "docs-mongodbatlas-datasource-cloud_backup_snapshot_restore_jobs"
description: |-
    Provides a Cloud Backup Snapshot Restore Jobs Datasource.
---

# Data Source: mongodbatlas_cloud_backup_snapshot_restore_jobs

`mongodbatlas_cloud_backup_snapshot_restore_jobs` provides a Cloud Backup Snapshot Restore Jobs datasource. Gets all the cloud backup snapshot restore jobs for the specified cluster.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage
First create a snapshot of the desired cluster. Then request that snapshot be restored in an automated fashion to the designated cluster and project.

```terraform
resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id          = "5cf5a45a9ccf6400e60981b6"
  cluster_name      = "MyCluster"
  description       = "MyDescription"
  retention_in_days = 1
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id     = "5cf5a45a9ccf6400e60981b6"
  cluster_name = "MyCluster"
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id
  delivery_type_config {
    automated = true
    target_cluster_name = "MyCluster"
    target_project_id     = "5cf5a45a9ccf6400e60981b6"
  }
}

data "mongodbatlas_cloud_backup_snapshot_restore_jobs" "test" {
  project_id     = mongodbatlas_cloud_backup_snapshot_restore_job.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot_restore_job.test.cluster_name
  page_num = 1
  items_per_page = 5
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster for which you want to retrieve restore jobs.
* `page_num` - (Optional)  	The page to return. Defaults to `1`.
* `items_per_page` - (Optional) Number of items to return per page, up to a maximum of 500. Defaults to `100`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `results` - Includes cloudProviderSnapshotRestoreJob object for each item detailed in the results array section.
* `totalCount` - Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.

### CloudProviderSnapshotRestoreJob

* `cancelled` -	Indicates whether the restore job was canceled.
* `created_at` -	UTC ISO 8601 formatted point in time when Atlas created the restore job.
* `delivery_type` - Type of restore job to create. Possible values are: automated and download.
* `delivery_url` -	One or more URLs for the compressed snapshot files for manual download. Only visible if deliveryType is download.
* `expired` -	Indicates whether the restore job expired.
* `expires_at` -	UTC ISO 8601 formatted point in time when the restore job expires.
* `finished_at` -	UTC ISO 8601 formatted point in time when the restore job completed.
* `id` -	The unique identifier of the restore job.
* `snapshot_id` -	Unique identifier of the source snapshot ID of the restore job.
* `target_project_id` -	Name of the target Atlas project of the restore job. Only visible if deliveryType is automated.
* `target_cluster_name` -	Name of the target Atlas cluster to which the restore job restores the snapshot. Only visible if deliveryType is automated.
* `timestamp` - Timestamp in ISO 8601 date and time format in UTC when the snapshot associated to snapshotId was taken.
* `oplogTs` - Timestamp in the number of seconds that have elapsed since the UNIX epoch.
* `oplogInc` - Oplog operation number from which to you want to restore this snapshot. 
* `pointInTimeUTCSeconds` - Timestamp in the number of seconds that have elapsed since the UNIX epoch.


For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/restore/get-all-restore-jobs/)