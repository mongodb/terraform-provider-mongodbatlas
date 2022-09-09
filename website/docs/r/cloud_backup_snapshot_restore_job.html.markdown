---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_snapshot_restore_job"
sidebar_current: "docs-mongodbatlas-resource-cloud_backup_snapshot_restore_job"
description: |-
    Provides a Cloud Backup Snapshot Restore Job resource.
---

# Resource: mongodbatlas_cloud_backup_snapshot_restore_job

`mongodbatlas_cloud_backup_snapshot_restore_job` provides a resource to create a new restore job from a cloud backup snapshot of a specified cluster. The restore job must define one of three delivery types: 
* **automated:** Atlas automatically restores the snapshot with snapshotId to the Atlas cluster with name targetClusterName in the Atlas project with targetGroupId.

* **download:** Atlas provides a URL to download a .tar.gz of the snapshot with snapshotId. The contents of the archive contain the data files for your Atlas cluster.

* **pointInTime:**  Atlas performs a Continuous Cloud Backup restore.

-> **Important:** If you specify `deliveryType` : `automated` or `deliveryType` : `pointInTime` in your request body to create an automated restore job, Atlas removes all existing data on the target cluster prior to the restore.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **API Key Access List**: This resource requires an Atlas API Access Key List to utilize this feature. This means to manage this resources you must have the IP address or CIDR block that the Terraform connection is coming from added to the Atlas API Key Access List of the Atlas API key you are using. See [Resources that require API Key List](https://www.mongodb.com/docs/atlas/configure-api-access/#use-api-resources-that-require-an-access-list) for details.
## Example Usage

### Example automated delivery type.

```terraform
  resource "mongodbatlas_cluster" "my_cluster" {
    project_id   = "5cf5a45a9ccf6400e60981b6"
    name         = "MyCluster"
    disk_size_gb = 5

  //Provider Settings "block"
    provider_name               = "AWS"
    provider_region_name        = "EU_WEST_2"
    provider_instance_size_name = "M10"
    cloud_backup                = true   // enable cloud backup snapshots
  }

  resource "mongodbatlas_cloud_provider_snapshot" "test" {
    project_id        = mongodbatlas_cluster.my_cluster.project_id
    cluster_name      = mongodbatlas_cluster.my_cluster.name
    description       = "myDescription"
    retention_in_days = 1
  }

  resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
    project_id      = mongodbatlas_cloud_provider_snapshot.test.project_id
    cluster_name    = mongodbatlas_cloud_provider_snapshot.test.cluster_name
    snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
    delivery_type_config   {
      automated           = true
      target_cluster_name = "MyCluster"
      target_project_id   = "5cf5a45a9ccf6400e60981b6"
    }
  }
```

### Example download delivery type.

```terraform
  resource "mongodbatlas_cluster" "my_cluster" {
    project_id   = "5cf5a45a9ccf6400e60981b6"
    name         = "MyCluster"
    disk_size_gb = 5

  //Provider Settings "block"
    provider_name               = "AWS"
    provider_region_name        = "EU_WEST_2"
    provider_instance_size_name = "M10"
    cloud_backup                = true   // enable cloud backup snapshots
  }

  resource "mongodbatlas_cloud_provider_snapshot" "test" {
    project_id        = mongodbatlas_cluster.my_cluster.project_id
    cluster_name      = mongodbatlas_cluster.my_cluster.name
    description       = "myDescription"
    retention_in_days = 1
  }
  
  resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
    project_id      = mongodbatlas_cloud_provider_snapshot.test.project_id
    cluster_name    = mongodbatlas_cloud_provider_snapshot.test.cluster_name
    snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
    delivery_type_config {
      download = true
    }
  }
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster whose snapshot you want to restore.
* `cluster_name` - (Required) The name of the Atlas cluster whose snapshot you want to restore.
* `snapshot_id` - (Required) Unique identifier of the snapshot to restore.
* `delivery_type_config` - (Required) Type of restore job to create. Possible configurations are: **download**, **automated**, or **pointInTime** only one must be set it in ``true``.
* `delivery_type_config.automated` - Set to `true` to use the automated configuration.
* `delivery_type_config.download` - Set to `true` to use the download configuration.
* `delivery_type_config.pointInTime` - Set to `true` to use the pointInTime configuration.
* `delivery_type_config.target_cluster_name` - Name of the target Atlas cluster to which the restore job restores the snapshot. Required for **automated** and **pointInTime**.
* `delivery_type_config.target_project_id` - Name of the target Atlas cluster to which the restore job restores the snapshot. Required for **automated** and **pointInTime**.
* `delivery_type_config.oplog_ts` - Optional setting for **pointInTime** configuration. Timestamp in the number of seconds that have elapsed since the UNIX epoch from which to you want to restore this snapshot. This is the first part of an Oplog timestamp.
* `delivery_type_config.oplog_inc` - Optional setting for **pointInTime** configuration. Oplog operation number from which to you want to restore this snapshot. This is the second part of an Oplog timestamp. Used in conjunction with `oplog_ts`.
* `delivery_type_config.point_in_time_utc_seconds` - Optional setting for **pointInTime** configuration. Timestamp in the number of seconds that have elapsed since the UNIX epoch from which you want to restore this snapshot. Used instead of oplog settings.

### Download
Atlas provides a URL to download a .tar.gz of the snapshot with snapshotId. 

### Automated
Atlas automatically restores the snapshot with snapshotId to the Atlas cluster with name targetClusterName in the Atlas project with targetProjectId. if you want to use automated delivery type, you must to set the arguments for the afformentioned properties.

### Point in time

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshot_restore_job_id` - The unique identifier of the restore job.
* `cancelled` -	Indicates whether the restore job was canceled.
* `created_at` -	UTC ISO 8601 formatted point in time when Atlas created the restore job.
* `delivery_type_config` - Type of restore job to create. Possible values are: automated and download.
* `delivery_url` -	One or more URLs for the compressed snapshot files for manual download. Only visible if deliveryType is download.
* `expired` -	Indicates whether the restore job expired.
* `expires_at` -	UTC ISO 8601 formatted point in time when the restore job expires.
* `finished_at` -	UTC ISO 8601 formatted point in time when the restore job completed.
* `id` -	The Terraform's unique identifier used internally for state management.
* `links` -	One or more links to sub-resources and/or related resources. The relations between URLs are explained in the Web Linking Specification.
* `snapshot_id` -	Unique identifier of the source snapshot ID of the restore job.
* `target_project_id` -	Name of the target Atlas project of the restore job. Only visible if deliveryType is automated.
* `target_cluster_name` -	Name of the target Atlas cluster to which the restore job restores the snapshot. Only visible if deliveryType is automated.
* `timestamp` - Timestamp in ISO 8601 date and time format in UTC when the snapshot associated to snapshotId was taken.
* `oplogTs` - Timestamp in the number of seconds that have elapsed since the UNIX epoch from which to you want to restore this snapshot.
    Three conditions apply to this parameter:
    * Enable Continuous Cloud Backup on your cluster.
    * Specify oplogInc.
    * Specify either oplogTs and oplogInc or pointInTimeUTCSeconds, but not both.
* `oplogInc` - Oplog operation number from which to you want to restore this snapshot. This is the second part of an Oplog timestamp.
    Three conditions apply to this parameter:
    * Enable Continuous Cloud Backup on your cluster.
    * Specify oplogTs.
    * Specify either oplogTs and oplogInc or pointInTimeUTCSeconds, but not both.
* `pointInTimeUTCSeconds` - Timestamp in the number of seconds that have elapsed since the UNIX epoch from which you want to restore this snapshot.
    Two conditions apply to this parameter:
    * Enable Continuous Cloud Backup on your cluster.
    * Specify either pointInTimeUTCSeconds or oplogTs and oplogInc, but not both.

## Import

Cloud Backup Snapshot Restore Job entries can be imported using project project_id, cluster_name and snapshot_id (Unique identifier of the snapshot), in the format `PROJECTID-CLUSTERNAME-JOBID`, e.g.

```
$ terraform import mongodbatlas_cloud_backup_snapshot_restore_job.test 5cf5a45a9ccf6400e60981b6-MyCluster-5d1b654ecf09a24b888f4c79
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/restore/restores/)