---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_provider_snapshot"
sidebar_current: "docs-mongodbatlas-resource-cloud_provider_snapshot"
description: |-
    Provides an Cloud Provider Snapshot resource.
---

# mongodbatlas_cloud_provider_snapshot

`mongodbatlas_cloud_provider_snapshot` provides an Cloud Provider Snapshot entry resource. Atlas Cloud Provider Snapshots provide localized backup storage using the native snapshot functionality of the cluster’s cloud service provider.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_cloud_provider_snapshot" "test" {
  group_id          = "<GROUP-ID>"
  cluster_name      = "MyClusterName"
  description       = "SomeDescription"
  retention_in_days = 1
}
```

## Argument Reference

* `group_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshots you want to retrieve.
* `description` - (Required) Description of the on-demand snapshot.
* `retention_in_days` - (Required) The number of days that Atlas should retain the on-demand snapshot. Must be at least 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier of the snapshot.
* `created_at` - UTC ISO 8601 formatted point in time when Atlas took the snapshot.
* `expires_at` - UTC ISO 8601 formatted point in time when Atlas will delete the snapshot.
* `master_key_uuid` - Unique ID of the AWS KMS Customer Master Key used to encrypt the snapshot. Only visible for clusters using Encryption at Rest via Customer KMS.
* `mongod_version` - Version of the MongoDB server.
* `snapshot_type` - Specified the type of snapshot. Valid values are onDemand and scheduled.
* `status` - Current status of the snapshot. One of the following values: queued, inProgress, completed, failed.
* `storage_size_bytes` - Specifies the size of the snapshot in bytes.
* `type` - Specifies the type of cluster: replicaSet or shardedCluster.
## Import

Cloud Provider Snapshot entries can be imported using project group_id, cluster_name and snapshot_id (Unique identifier of the snapshot), in the format `GROUPID-CLUSTERNAME-SNAPSHOTID`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_snapshot.test 5d0f1f73cf09a29120e173cf-MyClusterTest-5d116d82014b764445b2f9b5
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot/)