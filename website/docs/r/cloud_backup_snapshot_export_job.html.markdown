---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_snapshot_export_job"
sidebar_current: "docs-mongodbatlas-resource-cloud_backup_snapshot_export_job"
description: |-
    Provides a Cloud Backup Snapshot Export Job resource.
---

# mongodbatlas_cloud_backup_snapshot_export_job
`mongodbatlas_cloud_backup_snapshot_export_job` resource allows you to create a cloud backup snapshot export job for the specified project. 


-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform

resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id = "{PROJECT_ID}"
  iam_role_id    = "{IAM_ROLE_ID}"
  bucket_name    = "example_bucket"
  cloud_provider = "AWS"
}

resource "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
  project_id   = "{PROJECT_ID}"
  cluster_name = "{CLUSTER_NAME}"
  snapshot_id = "{SNAPSHOT_ID}"
  export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id
  
  custom_data {
    key   = "exported by"
    value = "myName"
  }
}

```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies the project which contains the Atlas cluster whose snapshot you want to export.
* `cluster_name` - (Required) Name of the Atlas cluster whose snapshot you want to export.
* `snapshot_id` - (Required) Unique identifier of the Cloud Backup snapshot to export. If necessary, use the [Get All Cloud Backups](https://docs.atlas.mongodb.com/reference/api/cloud-backup/backup/get-all-backups/) API to retrieve the list of snapshot IDs for a cluster or use the data source [mongodbatlas_cloud_cloud_backup_snapshots](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/cloud_backup_snapshots)
* `export_bucket_id` - (Required) Unique identifier of the AWS bucket to export the Cloud Backup snapshot to. If necessary, use the [Get All Snapshot Export Buckets](https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/get-all-export-buckets/) API to retrieve the IDs of all available export buckets for a project or use the data source [mongodbatlas_cloud_backup_snapshot_export_buckets](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/backup_snapshot_export_buckets)
* `custom_data` - (Optional) Custom data to include in the metadata file named `.complete` that Atlas uploads to the bucket when the export job finishes. Custom data can be specified as key and value pairs.

### Custom Data
* `key` - (Required) Required if you want to include custom data using `custom_data` in the metadata file uploaded to the bucket. Key to include in the metadata file that Atlas uploads to the bucket when the export job finishes.
* `value` - (Required) Required if you specify `key`.



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `components` - _Returned for sharded clusters only._ Export job details for each replica set in the sharded cluster.
* `created_at` - Timestamp in ISO 8601 date and time format in UTC when the export job was created.
* `err_msg` - Error message, only if the export job failed.
* `export_status` - _Returned for replica set only._ Status of the export job.
* `finished_at` - Timestamp in ISO 8601 date and time format in UTC when the export job completes.
* `export_job_id` - Unique identifier of the export job.
* `prefix ` - Full path on the cloud provider bucket to the folder where the snapshot is exported. The path is in the following format:`/exported_snapshots/{ORG-NAME}/{PROJECT-NAME}/{CLUSTER-NAME}/{SNAPSHOT-INITIATION-DATE}/{TIMESTAMP}`
* `state` - Status of the export job. Value can be one of the following:
    * `Queued` - indicates that the export job is queued
    * `InProgress` - indicates that the snapshot is being exported
    * `Successful` - indicates that the export job has completed successfully
    * `Failed` - indicates that the export job has failed

### components
* `export_id` - _Returned for sharded clusters only._ Export job details for each replica set in the sharded cluster.
* `replica_set_name` - _Returned for sharded clusters only._ Unique identifier of the export job for the replica set.

### export_status
* `exported_collections` - _Returned for replica set only._ Number of collections that have been exported.
* `total_collections` - _Returned for replica set only._ Total number of collections to export.

## Import

Cloud Backup Snapshot Export Backup entries can be imported using project project_id, cluster_name and export_job_id (Unique identifier of the snapshot export job), in the format `PROJECTID-CLUSTERNAME-EXPORTJOBID`, e.g.

```
$ terraform import mongodbatlas_cloud_backup_snapshot_export_job.test 5d0f1f73cf09a29120e173cf-5d116d82014b764445b2f9b5-5d116d82014b764445b2f9b5
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/create-one-export-job/)
