---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_shared_tier_snapshot"
sidebar_current: "docs-mongodbatlas-datasource-cloud_provider_shared_tier_snapshot"
description: |-
    Provides an Cloud Backup Snapshot Datasource for Shared Tier Clusters.
---

# Data Source: mongodbatlas_shared_tier_snapshot

`mongodbatlas_shared_tier_snapshot` provides an Cloud Backup Snapshot data source for Shared Tier Clusters. Atlas Cloud Backup Snapshots provide localized backup storage using the native snapshot functionality of the cluster’s cloud service.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.
-> **NOTE:** This data source is only for Shared Tier Clusters (M2 and M5). See [here](https://www.mongodb.com/docs/atlas/reference/free-shared-limitations/) for more details on Shared Tier Cluster Limitations. 


## Example Usage

```terraform
data "mongodbatlas_shared_tier_snapshot" "test" {
  project_id          = "5d0f1f73cf09a29120e173cf"
  cluster_name      = "MyClusterTest"
  snapshot_id       = "5d1285acd5ec13b6c2d1726a"
}
```

## Argument Reference

* `snapshot_id` - (Required) Unique 24-hexadecimal digit string that identifies the desired snapshot.
* `cluster_name` - (Required) Human-readable label that identifies the cluster.
* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project..

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Phase of the workflow for this snapshot at the time this resource made this request. Values: `PENDING` `QUEUED` `RUNNING` `FAILED` `COMPLETED`.
* `mongo_db_version` - MongoDB host version that the snapshot runs.
* `expiration` - Date and time when the download link no longer works. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `start_time` - Date and time when MongoDB Cloud began taking the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `finish_time` - Date and time when MongoDB Cloud completed writing this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `scheduled_time` - 	Date and time when MongoDB Cloud will take the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Shared-Tier-Snapshots/operation/getSharedClusterBackup)