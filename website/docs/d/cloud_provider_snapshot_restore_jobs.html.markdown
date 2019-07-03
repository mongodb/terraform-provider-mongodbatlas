---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_provider_snapshot_restore_jobs"
sidebar_current: "docs-mongodbatlas-datasource-cloud_provider_snapshot_restore_jobs"
description: |-
    Provides an Cloud Provider Snapshot Restore Jobs Datasource.
---

# mongodbatlas_cloud_provider_snapshot_restore_jobs

`mongodbatlas_cloud_provider_snapshot_restore_jobs` provides an Cloud Provider Snapshot Restore Jobs entry datasource. Get all cloud provider snapshot restore jobs for the specified cluster.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_cloud_provider_snapshot" "test" {
  group_id          = "5cf5a45a9ccf6400e60981b6"
  cluster_name      = "MyCluster"
  description       = "MyDescription"
  retention_in_days = 1
}

resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
  group_id     = "5cf5a45a9ccf6400e60981b6"
  cluster_name = "MyCluster"
  snapshot_id  = "${mongodbatlas_cloud_provider_snapshot.test.id}"
  delivery_type = {
    automated = true
    target_cluster_name = "MyCluster"
    target_group_id     = "5cf5a45a9ccf6400e60981b6"
  }
}

data "mongodbatlas_cloud_provider_snapshot_restore_jobs" "test" {
  group_id     = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.group_id}"
  cluster_name = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.cluster_name}"
}
```

## Argument Reference

* `group_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster for which you want to retrieve restore jobs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `results` - Includes cloudProviderSnapshotRestoreJob object for each item detailed in the results array section.
* `totalCount` - Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.

### CloudProviderSnapshotRestoreJob

* `cancelled` -	Indicates whether the restore job was canceled.
* `createdAt` -	UTC ISO 8601 formatted point in time when Atlas created the restore job.
* `deliveryType` - Type of restore job to create. Possible values are: automated and download.
* `deliveryUrl` -	One or more URLs for the compressed snapshot files for manual download. Only visible if deliveryType is download.
* `expired` -	Indicates whether the restore job expired.
* `expiresAt` -	UTC ISO 8601 formatted point in time when the restore job expires.
* `finishedAt` -	UTC ISO 8601 formatted point in time when the restore job completed.
* `id` -	The unique identifier of the restore job.
* `snapshotId` -	Unique identifier of the source snapshot ID of the restore job.
* `targetGroupId` -	Name of the target Atlas project of the restore job. Only visible if deliveryType is automated.
* `targetClusterName` -	Name of the target Atlas cluster to which the restore job restores the snapshot. Only visible if deliveryType is automated.
* `timestamp` - Timestamp in ISO 8601 date and time format in UTC when the snapshot associated to snapshotId was taken.


For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs-get-all/)