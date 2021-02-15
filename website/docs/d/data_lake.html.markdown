---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: data_lake"
sidebar_current: "docs-mongodbatlas-datasource-data-lake"
description: |-
    Describes a Data Lake.
---

# mongodbatlas_data_lake

`mongodbatlas_data_lake` describe a Data Lake.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_project" "test" {
  name   = "NAME OF THE PROJECT"
  org_id = "ORGANIZATION ID"
}
resource "mongodbatlas_cloud_provider_access" "test" {
  project_id = mongodbatlas_project.test.id
  provider_name = "AWS"
  iam_assumed_role_arn = "AWS ROLE ID"
}

resource "mongodbatlas_data_lake" "basic_ds" {
  project_id         = mongodbatlas_project.test.id
  name = "DATA LAKE NAME"
  aws{
    role_id = mongodbatlas_cloud_provider_access.test.role_id
    test_s3_bucket = "TEST S3 BUCKET NAME"
  }
}

data "mongodbatlas_data_lake" "test" {
  project_id           = mongodbatlas_data_lake.test.project_id
  name = mongodbatlas_data_lake.test.name
}
```

## Argument Reference

* `name` - (Required) Name of the data lake.
* `project_id` - (Required) The unique ID for the project to create a data lake.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `aws_role_id` - Unique identifier of the role that Data Lake can use to access the data stores.
* `aws_test_s3_bucket` - Name of the S3 data bucket that the provided role ID is authorized to access.
* `data_process_region` - The cloud provider region to which Atlas Data Lake routes client connections for data processing.
    * `data_process_region.cloud_provider` - Name of the cloud service provider. 
    * `data_process_region.region` -Name of the region to which Data Lake routes client connections for data processing.
* `aws_iam_assumed_role_arn` - Amazon Resource Name (ARN) of the IAM Role that Data Lake assumes when accessing S3 Bucket data stores.

  For more information on S3 actions, see [Actions, Resources, and Condition Keys for Amazon S3](https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazons3.html).

* `aws_iam_user_arn` - Amazon Resource Name (ARN) of the user that Data Lake assumes when accessing S3 Bucket data stores.
* `aws_external_id` - Unique identifier associated with the IAM Role that Data Lake assumes when accessing the data stores.
* `hostnames` - The list of hostnames assigned to the Atlas Data Lake. Each string in the array is a hostname assigned to the Atlas Data Lake.
* `state` - Current state of the Atlas Data Lake:
    * `ACTIVE` - The Data Lake is active and verified. You can query the data stores associated with the Atlas Data Lake.
* `storage_databases` - Configuration details for mapping each data store to queryable databases and collections.
    * `storage_databases.#.name` - Name of the database to which Data Lake maps the data contained in the data store.
    * `storage_databases.#.collections` -     Array of objects where each object represents a collection and data sources that map to a [stores](https://docs.mongodb.com/datalake/reference/format/data-lake-configuration#mongodb-datalakeconf-datalakeconf.stores) data store.
        * `storage_databases.#.collections.#.name` - Name of the collection.
            * `storage_databases.#.collections.#.data_sources` -     Array of objects where each object represents a stores data store to map with the collection.
                * `storage_databases.#.collections.#.data_sources.#.store_name` -     Name of a data store to map to the `<collection>`.
                * `storage_databases.#.collections.#.data_sources.#.default_format` - Default format that Data Lake assumes if it encounters a file without an extension while searching the storeName.
                * `storage_databases.#.collections.#.data_sources.#.path` - Controls how Atlas Data Lake searches for and parses files in the storeName before mapping them to the `<collection>`.
    * `storage_databases.#.views` -     Array of objects where each object represents an [aggregation pipeline](https://docs.mongodb.com/manual/core/aggregation-pipeline/#id1) on a collection.
    * `storage_databases.#.views.#.name` - Name of the view.
    * `storage_databases.#.views.#.source` -  Name of the source collection for the view.
    * `storage_databases.#.views.#.pipeline`- Aggregation pipeline stage(s) to apply to the source collection.
* `storage_stores` - Each object in the array represents a data store. Data Lake uses the storage.databases configuration details to map data in each data store to queryable databases and collections.
    * `storage_stores.#.name` - Name of the data store.
    * `storage_stores.#.provider` - Defines where the data is stored.
    * `storage_stores.#.region` - Name of the AWS region in which the S3 bucket is hosted.
    * `storage_stores.#.bucket` - Name of the AWS S3 bucket.
    * `storage_stores.#.prefix` - Prefix Data Lake applies when searching for files in the S3 bucket .
    * `storage_stores.#.delimiter` - The delimiter that separates `storage_databases.#.collections.#.data_sources.#.path` segments in the data store.
    * `storage_stores.#.include_tags` - Determines whether or not to use S3 tags on the files in the given path as additional partition attributes.
    
See [MongoDB Atlas API](https://docs.mongodb.com/datalake/reference/api/dataLakes-get-one-tenant) Documentation for more information.