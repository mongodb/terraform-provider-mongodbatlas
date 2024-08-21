# Data Source: mongodbatlas_federated_database_instances

`mongodbatlas_federated_database_instancess` provides a Federated Database Instance data source.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** All arguments including the password will be stored in the raw state as plain text. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html)

## Example Usages


```terraform
data "mongodbatlas_federated_database_instances" "test" {
  project_id         = "PROJECT ID"
  name               = "TENANT NAME OF THE FEDERATED DATABASE INSTANCE"
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create a Federated Database Instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `results` - Includes Federated Database instance for each item detailed in the results array section.

### Federated Database Instance

* `id` - The Terraform's unique identifier used internally for state management.
* `hostnames` - The list of hostnames assigned to the Federated Database Instance. Each string in the array is a hostname assigned to the Federated Database Instance.
* `state` - Current state of the Federated Database Instance:
  * `ACTIVE` - The Federated Database Instance is active and verified. You can query the data stores associated with the Federated Database Instance.
  * `DELETED` - The Federated Database Instance was deleted.
* `storage_databases` - Configuration details for mapping each data store to queryable databases and collections. For complete documentation on this object and its nested fields, see [databases](https://docs.mongodb.com/datalake/reference/format/data-lake-configuration#std-label-datalake-databases-reference). An empty object indicates that the Federated Database Instance has no mapping configuration for any data store.
  * `storage_databases.#.name` - Name of the database to which the Federated Database Instance maps the data contained in the data store.
  * `storage_databases.#.collections` -     Array of objects where each object represents a collection and data sources that map to a [stores](https://docs.mongodb.com/datalake/reference/format/data-lake-configuration#mongodb-datalakeconf-datalakeconf.stores) data store.
    * `storage_databases.#.collections.#.name` - Name of the collection.
      * `storage_databases.#.collections.#.data_sources` -     Array of objects where each object represents a stores data store to map with the collection.
        * `storage_databases.#.collections.#.data_sources.#.store_name` -     Name of a data store to map to the `<collection>`. Must match the name of an object in the stores array.
        * `storage_databases.#.collections.#.data_sources.#.dataset_name` -     Human-readable label that identifies the dataset that Atlas generates for an ingestion pipeline run or Online Archive.
        * `storage_databases.#.collections.#.data_sources.#.default_format` - Default format that Federated Database assumes if it encounters a file without an extension while searching the storeName. 
        * `storage_databases.#.collections.#.data_sources.#.path` - File path that controls how MongoDB Cloud searches for and parses files in the storeName before mapping them to a collection. Specify / to capture all files and folders from the prefix path.
        * `storage_databases.#.collections.#.data_sources.#.database` - Human-readable label that identifies the database, which contains the collection in the cluster.
        * `storage_databases.#.collections.#.data_sources.#.allow_insecure` - Flag that validates the scheme in the specified URLs. If true, allows insecure HTTP scheme, doesn't verify the server's certificate chain and hostname, and accepts any certificate with any hostname presented by the server. If false, allows secure HTTPS scheme only.
        * `storage_databases.#.collections.#.data_sources.#.database_regex` - Regex pattern to use for creating the wildcard database.
        * `storage_databases.#.collections.#.data_sources.#.collection` - Human-readable label that identifies the collection in the database.
        * `storage_databases.#.collections.#.data_sources.#.collection_regex` - Regex pattern to use for creating the wildcard (*) collection.
        * `storage_databases.#.collections.#.data_sources.#.provenance_field_name` - Name for the field that includes the provenance of the documents in the results.
        * `storage_databases.#.collections.#.data_sources.#.storeName` - Human-readable label that identifies the data store that MongoDB Cloud maps to the collection.
        * `storage_databases.#.collections.#.data_sources.#.urls` - URLs of the publicly accessible data files. You can't specify URLs that require authentication.
  * `storage_databases.#.views` -     Array of objects where each object represents an [aggregation pipeline](https://docs.mongodb.com/manual/core/aggregation-pipeline/#id1) on a collection. To learn more about views, see [Views](https://docs.mongodb.com/manual/core/views/).
    * `storage_databases.#.views.#.name` - Name of the view.
    * `storage_databases.#.views.#.source` -  Name of the source collection for the view.
    * `storage_databases.#.views.#.pipeline`- Aggregation pipeline stage(s) to apply to the source collection.
* `storage_stores` - Each object in the array represents a data store. Federated Database uses the storage.databases configuration details to map data in each data store to queryable databases and collections. For complete documentation on this object and its nested fields, see [stores](https://docs.mongodb.com/datalake/reference/format/data-lake-configuration#std-label-datalake-stores-reference). An empty object indicates that the Federated Database Instance has no configured data stores.
  * `storage_stores.#.name` - Name of the data store.
  * `storage_stores.#.provider` - Defines where the data is stored.
  * `storage_stores.#.region` - Name of the AWS region in which the S3 bucket is hosted.
  * `storage_stores.#.bucket` - Name of the AWS S3 bucket.
  * `storage_stores.#.prefix` - Prefix the Federated Database Instance applies when searching for files in the S3 bucket.
  * `storage_stores.#.delimiter` - The delimiter that separates `storage_databases.#.collections.#.data_sources.#.path` segments in the data store.
  * `storage_stores.#.include_tags` - Determines whether or not to use S3 tags on the files in the given path as additional partition attributes.
  * `storage_stores.#.cluster_name` - Human-readable label of the MongoDB Cloud cluster on which the store is based.
  * `storage_stores.#.cluster_id` - ID of the Cluster the Online Archive belongs to.
  * `storage_stores.#.allow_insecure` - Flag that validates the scheme in the specified URLs.
  * `storage_stores.#.public` - Flag that indicates whether the bucket is public.
  * `storage_stores.#.default_format` - Default format that Data Lake assumes if it encounters a file without an extension while searching the storeName.
  * `storage_stores.#.urls` - Comma-separated list of publicly accessible HTTP URLs where data is stored.
  * `storage_stores.#.read_preference` - MongoDB Cloud cluster read preference, which describes how to route read requests to the cluster.
    * `storage_stores.#.read_preference.maxStalenessSeconds` - Maximum replication lag, or staleness, for reads from secondaries.
    * `storage_stores.#.read_preference.mode` - Read preference mode that specifies to which replica set member to route the read requests.
    * `storage_stores.#.read_preference.tag_sets` - List that contains tag sets or tag specification documents.
      * `storage_stores.#.read_preference.tags` - List of all tags within a tag set
        * `storage_stores.#.read_preference.tags.name` - Human-readable label of the tag.
        * `storage_stores.#.read_preference.tags.value` - Value of the tag.

### `cloud_provider_config` - Cloud provider linked to this data federated instance.
#### `aws` - AWS provider of the cloud service where the Federated Database Instance can access the S3 Bucket.
* `iam_assumed_role_arn` - Amazon Resource Name (ARN) of the IAM Role that the Federated Database Instance assumes when accessing S3 Bucket data stores. The IAM Role must support the following actions against each S3 bucket:
  * `s3:GetObject`
  * `s3:ListBucket`
  * `s3:GetObjectVersion` 
    
  For more information on S3 actions, see [Actions, Resources, and Condition Keys for Amazon S3](https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazons3.html).

* `iam_user_arn` - Amazon Resource Name (ARN) of the user that the Federated Database Instance assumes when accessing S3 Bucket data stores.
* `external_id` - Unique identifier associated with the IAM Role that the Federated Database Instance assumes when accessing the data stores.
* `role_id` - Unique identifier of the role that the data lake can use to access the data stores.
#### `data_process_region` - The cloud provider region to which the Federated Instance routes client connections for data processing.
* `cloud_provider` -  Name of the cloud service provider. Atlas Federated Database only supports AWS.
* `region` - Name of the region to which the Federanted Instnace routes client connections for data processing. 

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation) Documentation for more information.