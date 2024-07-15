---
subcategory: "Deprecated"    
---

**WARNING:** This resource is deprecated, use `mongodbatlas_cloud_backup_snapshot`
**Note:** This resource have now been fully deprecated as part of v1.10.0 release

# Resource: mongodbatlas_cloud_provider_snapshot

`mongodbatlas_cloud_provider_snapshot` provides a resource to take a cloud backup snapshot on demand.
On-demand snapshots happen immediately, unlike scheduled snapshots which occur at regular intervals. If there is already an on-demand snapshot with a status of queued or inProgress, you must wait until Atlas has completed the on-demand snapshot before taking another.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
  resource "mongodbatlas_cluster" "my_cluster" {
    project_id   = "5cf5a45a9ccf6400e60981b6"
    name         = "MyCluster"

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
    timeout           = "10m"
  }
  
  resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
    project_id      = mongodbatlas_cloud_provider_snapshot.test.project_id
    cluster_name    = mongodbatlas_cloud_provider_snapshot.test.cluster_name
    snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
    delivery_type_config {
      download = true
    }
  }
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshots you want to retrieve.
* `description` - (Required) Description of the on-demand snapshot.
* `retention_in_days` - (Required) The number of days that Atlas should retain the on-demand snapshot. Must be at least 1.
* `timeout`- (Optional) The duration of time to wait to finish the on-demand snapshot. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. Default value for the timeout is `10m`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshot_id` - Unique identifier of the snapshot.
* `id` -	Unique identifier used for terraform for internal manages.
* `created_at` - UTC ISO 8601 formatted point in time when Atlas took the snapshot.
* `description` - Description of the snapshot. Only present for on-demand snapshots.
* `expires_at` - UTC ISO 8601 formatted point in time when Atlas will delete the snapshot.
* `master_key_uuid` - Unique ID of the AWS KMS Customer Master Key used to encrypt the snapshot. Only visible for clusters using Encryption at Rest via Customer KMS.
* `mongod_version` - Version of the MongoDB server.
* `snapshot_type` - Specified the type of snapshot. Valid values are onDemand and scheduled.
* `status` - Current status of the snapshot. One of the following values will be returned: queued, inProgress, completed, failed.
* `storage_size_bytes` - Specifies the size of the snapshot in bytes.
* `type` - Specifies the type of cluster: replicaSet or shardedCluster.

## Import

Cloud Backup Snapshot entries can be imported using project project_id, cluster_name and snapshot_id (Unique identifier of the snapshot), in the format `PROJECTID-CLUSTERNAME-SNAPSHOTID`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_snapshot.test 5d0f1f73cf09a29120e173cf-MyClusterTest-5d116d82014b764445b2f9b5
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/backup/backups/)