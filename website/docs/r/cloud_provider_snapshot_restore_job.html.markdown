---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_provider_snapshot_restore_job"
sidebar_current: "docs-mongodbatlas-resource-cloud_provider_snapshot_restore_job"
description: |-
    Provides a Cloud Provider Snapshot Restore Job resource.
---

# mongodbatlas_cloud_provider_snapshot_restore_job

`mongodbatlas_cloud_provider_snapshot_restore_job` provides a resource to create a new restore job from a cloud provider snapshot of a specified cluster. The restore job can be one of two types: 
* **automated:** Atlas automatically restores the snapshot with snapshotId to the Atlas cluster with name targetClusterName in the Atlas project with targetGroupId.

* **download:** Atlas provides a URL to download a .tar.gz of the snapshot with snapshotId. The contents of the archive contain the data files for your Atlas cluster.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

### Example automated delivery type.

```hcl
  resource "mongodbatlas_cluster" "my_cluster" {
    project_id   = "5cf5a45a9ccf6400e60981b6"
    name         = "MyCluster"
    disk_size_gb = 5

  //Provider Settings "block"
    provider_name               = "AWS"
    provider_region_name        = "EU_WEST_2"
    provider_instance_size_name = "M10"
    provider_backup_enabled     = true   // enable cloud provider snapshots
    provider_disk_iops          = 100
    provider_encrypt_ebs_volume = false
  }

  resource "mongodbatlas_cloud_provider_snapshot" "test" {
    project_id        = mongodbatlas_cluster.my_cluster.project_id
    cluster_name      = mongodbatlas_cluster.my_cluster.name
    description       = "myDescription"
    retention_in_days = 1
  }

  resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
    project_id      = mongodbatlas_cloud_provider_snapshot.test.project_id
    cluster_name    = mongodbatlas_cloud_provider_snapshot.test.cluster_name
    snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
    delivery_type   = {
      automated           = true
      target_cluster_name = "MyCluster"
      target_project_id   = "5cf5a45a9ccf6400e60981b6"
    }
    depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
  }
```

### Example download delivery type.

```hcl
  resource "mongodbatlas_cluster" "my_cluster" {
    project_id   = "5cf5a45a9ccf6400e60981b6"
    name         = "MyCluster"
    disk_size_gb = 5

  //Provider Settings "block"
    provider_name               = "AWS"
    provider_region_name        = "EU_WEST_2"
    provider_instance_size_name = "M10"
    provider_backup_enabled     = true   // enable cloud provider snapshots
    provider_disk_iops          = 100
    provider_encrypt_ebs_volume = false
  }

  resource "mongodbatlas_cloud_provider_snapshot" "test" {
    project_id        = mongodbatlas_cluster.my_cluster.project_id
    cluster_name      = mongodbatlas_cluster.my_cluster.name
    description       = "myDescription"
    retention_in_days = 1
  }
  
  resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
    project_id      = mongodbatlas_cloud_provider_snapshot.test.project_id
    cluster_name    = mongodbatlas_cloud_provider_snapshot.test.cluster_name
    snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
    delivery_type = {
      download = true
    }
  }
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster whose snapshot you want to restore.
* `cluster_name` - (Required) The name of the Atlas cluster whose snapshot you want to restore.
* `snapshot_id` - (Required) Unique identifier of the snapshot to restore.
* `delivery_type` - (Required) Type of restore job to create. Possible values are: **download** or **automated**, only one must be set it in ``true``.

### Download
Atlas provides a URL to download a .tar.gz of the snapshot with snapshotId. 

### Automated
Atlas automatically restores the snapshot with snapshotId to the Atlas cluster with name targetClusterName in the Atlas project with targetGroupId. if you want to use automated delivery type, you must to set the following arguments:

* `target_cluster_name` - (Required) 	Name of the target Atlas cluster to which the restore job restores the snapshot. Only required if deliveryType is automated.
* `target_group_id` - (Required) 	Unique ID of the target Atlas project for the specified targetClusterName. Only required if deliveryType is automated.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshot_restore_job_id` - The unique identifier of the restore job.
* `cancelled` -	Indicates whether the restore job was canceled.
* `created_at` -	UTC ISO 8601 formatted point in time when Atlas created the restore job.
* `delivery_type` - Type of restore job to create. Possible values are: automated and download.
* `delivery_url` -	One or more URLs for the compressed snapshot files for manual download. Only visible if deliveryType is download.
* `expired` -	Indicates whether the restore job expired.
* `expires_at` -	UTC ISO 8601 formatted point in time when the restore job expires.
* `finished_at` -	UTC ISO 8601 formatted point in time when the restore job completed.
* `id` -	Unique identifier used for terraform for internal manages.
* `links` -	One or more links to sub-resources and/or related resources. The relations between URLs are explained in the Web Linking Specification.
* `snapshot_id` -	Unique identifier of the source snapshot ID of the restore job.
* `target_group_id` -	Name of the target Atlas project of the restore job. Only visible if deliveryType is automated.
* `target_cluster_name` -	Name of the target Atlas cluster to which the restore job restores the snapshot. Only visible if deliveryType is automated.
* `timestamp` - Timestamp in ISO 8601 date and time format in UTC when the snapshot associated to snapshotId was taken.

## Import

Cloud Provider Snapshot Restore Job entries can be imported using project project_id, cluster_name and snapshot_id (Unique identifier of the snapshot), in the format `PROJECTID-CLUSTERNAME-JOBID`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_snapshot_restore_job.test 5cf5a45a9ccf6400e60981b6-MyCluster-5d1b654ecf09a24b888f4c79
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs/)