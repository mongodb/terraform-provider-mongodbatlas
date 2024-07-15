# Resource: mongodbatlas_cloud_backup_snapshot_export_bucket

`mongodbatlas_cloud_backup_snapshot_export_bucket` allows you to create an export snapshot bucket for the specified project. 


-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id   = "{PROJECT_ID}"
  iam_role_id = "{IAM_ROLE_ID}"
  bucket_name = "example-bucket"
  cloud_provider = "AWS"
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `iam_role_id` - (Required) Unique identifier of the role that Atlas can use to access the bucket. You must also specify the `bucket_name`.
* `bucket_name` - (Required) Name of the bucket that the provided role ID is authorized to access. You must also specify the `iam_role_id`.
* `cloud_provider` - (Required) Name of the provider of the cloud service where Atlas can access the S3 bucket. Atlas only supports `AWS`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `export_bucket_id` -	Unique identifier of the snapshot export bucket.

## Import

Cloud Backup Snapshot Export Backup entries can be imported using project project_id, and bucket_id (Unique identifier of the snapshot export bucket), in the format `PROJECTID-BUCKETID`, e.g.

```
$ terraform import mongodbatlas_cloud_backup_snapshot_export_bucket.test 5d0f1f73cf09a29120e173cf-5d116d82014b764445b2f9b5
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/create-one-export-bucket/)
