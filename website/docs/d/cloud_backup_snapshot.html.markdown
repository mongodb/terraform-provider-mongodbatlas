---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_snapshot"
sidebar_current: "docs-mongodbatlas-datasource-cloud_backup_snapshot"
description: |-
    Provides a Cloud Backup Snapshot Datasource.
---

# Data Source: mongodbatlas_cloud_backup_snapshot

`mongodbatlas_cloud_backup_snapshot` provides an Cloud Backup Snapshot datasource. Atlas Cloud Backup Snapshots provide localized backup storage using the native snapshot functionality of the cluster’s cloud service.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_cloud_backup_snapshot" "test" {
  group_id          = "5d0f1f73cf09a29120e173cf"
  cluster_name      = "MyClusterTest"
  description       = "SomeDescription"
  retention_in_days = 1
}

data "mongodbatlas_cloud_backup_snapshot" "test" {
  snapshot_id  = "5d1285acd5ec13b6c2d1726a"
  group_id     = mongodbatlas_cloud_backup_snapshot.test.group_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
}
```

## Argument Reference

* `snapshot_id` - (Required) The unique identifier of the snapshot you want to retrieve.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshot you want to retrieve.
* `group_id` - (Required) The unique identifier of the project for the Atlas cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier of the snapshot.
* `created_at` - UTC ISO 8601 formatted point in time when Atlas took the snapshot.
* `expires_at` - UTC ISO 8601 formatted point in time when Atlas will delete the snapshot.
* `description` - UDescription of the snapshot. Only present for on-demand snapshots.
* `master_key_uuid` - Unique ID of the AWS KMS Customer Master Key used to encrypt the snapshot. Only visible for clusters using Encryption at Rest via Customer KMS.
* `mongod_version` - Version of the MongoDB server.
* `snapshot_type` - Specified the type of snapshot. Valid values are onDemand and scheduled.
* `status` - Current status of the snapshot. One of the following values: queued, inProgress, completed, failed.
* `storage_size_bytes` - Specifies the size of the snapshot in bytes.
* `type` - Specifies the type of cluster: replicaSet or shardedCluster.
* `cloud_provider` - Cloud provider that stores this snapshot.
* `members` - Block of List of snapshots and the cloud provider where the snapshots are stored. See below
* `replica_set_name` - Label given to the replica set from which Atlas took this snapshot.
* `snapshot_ids` - Unique identifiers of the snapshots created for the shards and config server for a sharded cluster. 

### members

* `cloud_provider` - Cloud provider that stores this snapshot.
* `id` - Unique identifier for the sharded cluster snapshot.
* `replica_set_name` - Label given to a shard or config server from which Atlas took this snapshot.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/backup/get-one-backup/)