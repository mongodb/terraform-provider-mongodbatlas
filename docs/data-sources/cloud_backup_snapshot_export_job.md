# Data Source: mongodbatlas_cloud_backup_snapshot_export_Job

`mongodbatlas_cloud_backup_snapshot_export_job` datasource allows you to retrieve a snapshot export job for the specified project and cluster.


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

data "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
  project_id   = "{PROJECT_ID}"
  cluster_name = "{CLUSTER_NAME}"
  export_job_id = mongodbatlas_cloud_backup_snapshot_export_job.test.export_job_id
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies the project which contains the Atlas cluster whose snapshot you want to retrieve.
* `cluster_name` - (Required) Name of the Atlas cluster whose export job you want to retrieve.
* `export_job_id` -(Required) Unique identifier of the export job to retrieve.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `snapshot_id` - Unique identifier of the Cloud Backup snapshot to export.
* `export_bucket_id` - Unique identifier of the AWS bucket to export the Cloud Backup snapshot to.
* `custom_data` - Custom data to include in the metadata file named `.complete` that Atlas uploads to the bucket when the export job finishes. Custom data can be specified as key and value pairs.
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

### Custom Data
* `key` - Custom data specified as key in the key and value pair.
* `value` - Value for the key specified using `key`.

### components
* `export_id` - _Returned for sharded clusters only._ Export job details for each replica set in the sharded cluster.
* `replica_set_name` - _Returned for sharded clusters only._ Unique identifier of the export job for the replica set.

### export_status
* `exported_collections` - _Returned for replica set only._ Number of collections that have been exported.
* `total_collections` - _Returned for replica set only._ Total number of collections to export.


For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/get-one-export-job/)
