# Data Source: mongodbatlas_cloud_backup_snapshot_export_buckets

`mongodbatlas_cloud_backup_snapshot_export_buckets` datasource allows you to retrieve all the buckets for the specified project.


-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id   = "{PROJECT_ID}"
  iam_role_id = "{IAM_ROLE_ID}"
  bucket_name = "example-bucket"
  cloud_provider = "AWS"
}

data "mongodbatlas_cloud_backup_snapshot_export_buckets" "test" {
  project_id   = "{PROJECT_ID}"
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `page_num` - (Optional)  	The page to return. Defaults to `1`.
* `items_per_page` - (Optional) Number of items to return per page, up to a maximum of 500. Defaults to `100`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `links` - One or more links to sub-resources and/or related resources.
* `results` - Includes CloudProviderSnapshotExportBucket object for each item detailed in the results array section.
* `totalCount` - Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.


### CloudProviderSnapshotExportBucket
* `project_id` - The unique identifier of the project for the Atlas cluster.
* `export_bucket_id` -	Unique identifier of the snapshot bucket id.
* `iam_role_id` - Unique identifier of the role that Atlas can use to access the bucket.
* `bucket_name` - Name of the bucket that the provided role ID is authorized to access.
* `cloud_provider` - Name of the provider of the cloud service where Atlas can access the S3 bucket.
* `role_id` - Unique identifier of the Azure Service Principal that Atlas can use to access the Azure Blob Storage Container.
* `service_url` - URL that identifies the blob Endpoint of the Azure Blob Storage Account.
* `tenant_id` - UUID that identifies the Azure Active Directory Tenant ID.


For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/create-one-export-bucket/)
