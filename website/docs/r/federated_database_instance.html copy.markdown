---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: federated_database_instance"
sidebar_current: "docs-mongodbatlas-resource-federated-database-instance"
description: |-
    Provides a Federated Database Instance resource.
---

# Resource: mongodbatlas_federated_database_instance

`mongodbatlas_federated_database_instance` provides a Federated Database Instance resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** All arguments including the password will be stored in the raw state as plain-text. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html)

## Example Usages


```terraform
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

resource "mongodbatlas_federated_database_instance" "test" {
  project_id         = "PROJECT ID"
  name = "NAME OF THE FEDERATED DATABASE INSTANCE"
  aws {
    role_id = "AWS ROLE ID"
    test_s3_bucket = "S3 BUCKET NAME"
  }
  storage_databases {
    name = "VirtualDatabase0"
    collections {
      name = "NAME OF THE COLLECTION"
      data_sources {
          collection = "COLLECTION IN THE CLUSTER"
          database = "DB IN THE CLUSTER"
          store_name =  "CLUSTER NAME"
      }
      data_sources {
          store_name = "S3 BUCKET NAME"
          path = "S3 BUCKET PATH"
      }
    }
  }

  storage_stores {
    name = "STORE 1 NAME"
    cluster_name = "CLUSTER NAME"
    project_id = "PROJECT ID"
    provider = "atlas"
    read_preference {
      mode = "secondary"
    }
  }

  storage_stores {
    bucket = "STORE 2 NAME"
    delimiter = "/"
    name = "S3 BUCKET NAME"
    prefix = "S3 BUCKET PREFIX"
    provider = "s3"
    region = "AWS REGION"
  }
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create a data lake.
* `name` - (Required) Name of the Atlas Federated Database Instance.
* `aws` - (Required) AWS provider of the cloud service where Data Lake can access the S3 Bucket.
  * `aws.0.role_id` - (Required) Unique identifier of the role that the Federated Instance can use to access the data stores. If necessary, use the Atlas [UI](https://docs.atlas.mongodb.com/security/manage-iam-roles/) or [API](https://docs.atlas.mongodb.com/reference/api/cloud-provider-access-get-roles/) to retrieve the role ID. You must also specify the `aws.0.test_s3_bucket`.
  * `aws.0.test_s3_bucket` - (Required) Name of the S3 data bucket that the provided role ID is authorized to access. You must also specify the `aws.0.role_id`.
* `data_process_region` - (Optional) The cloud provider region to which the Federated Instance routes client connections for data processing. Set to `null` to route client connections to the region nearest to the client based on DNS resolution.
  * `data_process_region.0.cloud_provider` - (Required) Name of the cloud service provider. Atlas Data Lake only supports AWS.
  * `data_process_region.0.region` - (Required). Name of the region to which the Federanted Instnace routes client connections for data processing. Atlas Federated Database only supports the following regions:
    * `SYDNEY_AUS` (ap-southeast-2)
    * `FRANKFURT_DEU` (eu-central-1)
    * `DUBLIN_IRL` (eu-west-1)
    * `LONDON_GBR` (eu-west-2)
    * `VIRGINIA_USA` (us-east-1)
    * `OREGON_USA` (us-west-2)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `aws.0.iam_assumed_role_arn` - Amazon Resource Name (ARN) of the IAM Role that Data Lake assumes when accessing S3 Bucket data stores. The IAM Role must support the following actions against each S3 bucket:
  * `s3:GetObject`
  * `s3:ListBucket`
  * `s3:GetObjectVersion` 
    
  For more information on S3 actions, see [Actions, Resources, and Condition Keys for Amazon S3](https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazons3.html).

* `aws.0.iam_user_arn` - Amazon Resource Name (ARN) of the user that Federated Database Instance assumes when accessing S3 Bucket data stores.
* `aws.0.external_id` - Unique identifier associated with the IAM Role that Data Lake assumes when accessing the data stores.
* `storage_databases` - Configuration details for mapping each data store to queryable databases and collections. For complete documentation on this object and its nested fields, see [databases](https://docs.mongodb.com/datalake/reference/format/data-lake-configuration#std-label-datalake-databases-reference). An empty object indicates that the Data Lake has no mapping configuration for any data store.
  * `storage_databases.#.name` - Name of the database to which Federated Database maps the data contained in the data store.
  * `storage_databases.#.collections` -     Array of objects where each object represents a collection and data sources that map to a [stores](https://docs.mongodb.com/datalake/reference/format/data-lake-configuration#mongodb-datalakeconf-datalakeconf.stores) data store.
    * `storage_databases.#.collections.#.name` - Name of the collection.
      * `storage_databases.#.collections.#.data_sources` -     Array of objects where each object represents a stores data store to map with the collection.
        * `storage_databases.#.collections.#.data_sources.#.store_name` -     Name of a data store to map to the `<collection>`. Must match the name of an object in the stores array.
        * `storage_databases.#.collections.#.data_sources.#.default_format` - Default format that Federated Database assumes if it encounters a file without an extension while searching the storeName. 
        * `storage_databases.#.collections.#.data_sources.#.path` - Controls how Atlas Federated Database searches for and parses files in the storeName before mapping them to the `<collection>`.
  * `storage_databases.#.views` -     Array of objects where each object represents an [aggregation pipeline](https://docs.mongodb.com/manual/core/aggregation-pipeline/#id1) on a collection. To learn more about views, see [Views](https://docs.mongodb.com/manual/core/views/).
  * `storage_databases.#.views.#.name` - Name of the view.
  * `storage_databases.#.views.#.source` -  Name of the source collection for the view.
  * `storage_databases.#.views.#.pipeline`- Aggregation pipeline stage(s) to apply to the source collection.
* `storage_stores` - Each object in the array represents a data store. Federated Database uses the storage.databases configuration details to map data in each data store to queryable databases and collections. For complete documentation on this object and its nested fields, see [stores](https://docs.mongodb.com/datalake/reference/format/data-lake-configuration#std-label-datalake-stores-reference). An empty object indicates that the Federated Database Instance has no configured data stores.
  * `storage_stores.#.name` - Name of the data store.
  * `storage_stores.#.provider` - Defines where the data is stored.
  * `storage_stores.#.region` - Name of the AWS region in which the S3 bucket is hosted.
  * `storage_stores.#.bucket` - Name of the AWS S3 bucket.
  * `storage_stores.#.prefix` - Prefix Data Lake applies when searching for files in the S3 bucket .
  * `storage_stores.#.delimiter` - The delimiter that separates `storage_databases.#.collections.#.data_sources.#.path` segments in the data store.
  * `storage_stores.#.include_tags` - Determines whether or not to use S3 tags on the files in the given path as additional partition attributes.

## Import

The Federated Database Instance can be imported using project ID, name of the instance and name of the AWS s3 bucket, in the format `project_id`--`name`--`aws_test_s3_bucket`, e.g.

```
$ terraform import mongodbatlas_federated_database_instance.example 1112222b3bf99403840e8934--test--s3-test
```

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation) Documentation for more information.