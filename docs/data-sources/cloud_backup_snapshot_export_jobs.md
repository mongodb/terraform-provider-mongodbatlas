# Data Source: mongodbatlas_cloud_backup_snapshot_export_jobs

`mongodbatlas_cloud_backup_snapshot_export_jobs` datasource allows you to retrieve all the buckets for the specified project.


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

data "mongodbatlas_cloud_backup_snapshot_export_jobs" "test" {
  project_id   = "{PROJECT_ID}"
  cluster_name = "{CLUSTER_NAME}"
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) Name of the Atlas cluster whose export job you want to retrieve.
* `page_num` - (Optional)  	The page to return. Defaults to `1`.
* `items_per_page` - (Optional) Number of items to return per page, up to a maximum of 500. Defaults to `100`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `links` - One or more links to sub-resources and/or related resources.
* `results` - Includes CloudProviderSnapshotExportJob object for each item detailed in the results array section.
* `totalCount` - Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.


### CloudProviderSnapshotExportJob
* `project_id` - The unique identifier of the project for the Atlas cluster.
* `export_job_id` -	Unique identifier of the S3 bucket.
* `snapshot_id` - Unique identifier of the Cloud Backup snapshot to export.
* `export_bucket_id` - Unique identifier of the AWS bucket to export the Cloud Backup snapshot to.
* `custom_data` - Custom data to include in the metadata file named `.complete` that Atlas uploads to the bucket when the export job finishes. Custom data can be specified as key and value pairs.
* `components` - _Returned for sharded clusters only._ Export job details for each replica set in the sharded cluster.
* `created_at` - Timestamp in ISO 8601 date and time format in UTC when the export job was created.
* `err_msg` - Error message, only if the export job failed. **Note:** This attribute is deprecated as it is not being used.
* `export_status` - _Returned for replica set only._ Status of the export job.
* `finished_at` - Timestamp in ISO 8601 date and time format in UTC when the export job completes.
* `export_job_id` - Unique identifier of the export job.
* `prefix ` - Full path on the cloud provider bucket to the folder where the snapshot is exported. The path is in the following format:`/exported_snapshots/{ORG-NAME}/{PROJECT-NAME}/{CLUSTER-NAME}/{SNAPSHOT-INITIATION-DATE}/{TIMESTAMP}`
* `state` - Status of the export job. Value can be one of the following:
    * `Queued` - indicates that the export job is queued
    * `InProgress` - indicates that the snapshot is being exported
    * `Successful` - indicates that the export job has completed successfully
    * `Failed` - indicates that the export job has failed

#### Custom Data
* `key` - Custom data specified as key in the key and value pair.
* `value` - Value for the key specified using `key`.

#### components
* `export_id` - _Returned for sharded clusters only._ Export job details for each replica set in the sharded cluster.
* `replica_set_name` - _Returned for sharded clusters only._ Unique identifier of the export job for the replica set.

#### export_status
* `exported_collections` - _Returned for replica set only._ Number of collections that have been exported.
* `total_collections` - _Returned for replica set only._ Total number of collections to export.




For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/get-all-export-jobs/)
