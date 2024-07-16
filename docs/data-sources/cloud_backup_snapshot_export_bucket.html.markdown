# Data Source: mongodbatlas_cloud_backup_snapshot_export_bucket

`mongodbatlas_cloud_backup_snapshot_export_bucket` datasource allows you to retrieve all the buckets for the specified project.


-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id   = "{PROJECT_ID}"
  iam_role_id = "{IAM_ROLE_ID}"
  bucket_name = "example-bucket"
  cloud_provider = "AWS"
}

data "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id   = "{PROJECT_ID}"
  export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `export_bucket_id` - (Required) Unique identifier of the snapshot export bucket.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `iam_role_id` - Unique identifier of the role that Atlas can use to access the bucket. You must also specify the `bucket_name`.
* `bucket_name` - Name of the bucket that the provided role ID is authorized to access. You must also specify the `iam_role_id`.
* `cloud_provider` - Name of the provider of the cloud service where Atlas can access the S3 bucket. Atlas only supports `AWS`.



For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/create-one-export-bucket/)
