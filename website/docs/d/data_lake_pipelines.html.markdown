# Data Source: mongodbatlas_data_lake_pipelines

`mongodbatlas_data_lake_pipelines` describes Data Lake Pipelines.

-> **NOTE:** Groups and projects are synonymous terms. You may find `group_id` in the official documentation.

## Example Usages


```terraform
data "mongodbatlas_data_lake_pipelines" "pipelineDataSource" {
  project_id       = <YOU-PROJECT-ID>
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create a data lake pipeline.


## Attributes Reference
* `results` - A list where each represents a Data Lake Pipeline.

### Data Lake Pipeline

In addition to all arguments above, the following attributes are exported:

* `id` -  Unique 24-hexadecimal digit string that identifies the Data Lake Pipeline.
* `created_date` - Timestamp that indicates when the Data Lake Pipeline was created.
* `last_updated_date` - Timestamp that indicates the last time that the Data Lake Pipeline was updated.
* `state` - State of this Data Lake Pipeline.
* `transformations` - Fields to be excluded for this Data Lake Pipeline.
  * `transformations.#.field` - Key in the document.
  * `transformations.#.type` - Type of transformation applied during the export of the namespace in a Data Lake Pipeline.
* `snapshots` - List of backup snapshots that you can use to trigger an on demand pipeline run.
  * `snapshots.#.id` - Unique 24-hexadecimal digit string that identifies the snapshot.
  * `snapshots.#.provider` - Human-readable label that identifies the cloud provider that stores this snapshot.
  * `snapshots.#.created_at` - Date and time when MongoDB Atlas took the snapshot.
  * `snapshots.#.expires_at` - Date and time when MongoDB Atlas deletes the snapshot.
  * `snapshots.#.frequency_type` - Human-readable label that identifies how often this snapshot triggers.
  * `snapshots.#.master_key` - Unique string that identifies the Amazon Web Services (AWS) Key Management Service (KMS) Customer Master Key (CMK) used to encrypt the snapshot.
  * `snapshots.#.mongod_version` - Version of the MongoDB host that this snapshot backs up.
  * `snapshots.#.replica_set_name` - Human-readable label that identifies the replica set from which MongoDB Atlas took this snapshot.
  * `snapshots.#.type` - Human-readable label that categorizes the cluster as a replica set or sharded cluster.
  * `snapshots.#.snapshot_type` - Human-readable label that identifies when this snapshot triggers.
  * `snapshots.#.status` - Human-readable label that indicates the stage of the backup process for this snapshot.
  * `snapshots.#.size` - List of backup snapshots that you can use to trigger an on demand pipeline run.
  * `snapshots.#.copy_region` - List that identifies the regions to which MongoDB Atlas copies the snapshot.
  * `snapshots.#.policies` - List that contains unique identifiers for the policy items.
* `ingestion_schedules` - List of backup schedule policy items that you can use as a Data Lake Pipeline source.
  * `ingestion_schedules.#.id` - Unique 24-hexadecimal digit string that identifies this backup policy item.
  * `ingestion_schedules.#.frequency_type` - Human-readable label that identifies the frequency type associated with the backup policy.
  * `ingestion_schedules.#.frequency_interval` - Number that indicates the frequency interval for a set of snapshots.
  * `ingestion_schedules.#.retention_unit` - Unit of time in which MongoDB Atlas measures snapshot retention.
  * `ingestion_schedules.#.retention_value` - Duration in days, weeks, or months that MongoDB Atlas retains the snapshot. 

### `sink` - Ingestion destination of a Data Lake Pipeline
  * `type` - Type of ingestion destination of this Data Lake Pipeline.
  * `provider` - Target cloud provider for this Data Lake Pipeline.
  * `region` - Target cloud provider region for this Data Lake Pipeline. [Supported cloud provider regions](https://www.mongodb.com/docs/datalake/limitations).
  * `partition_fields` - Ordered fields used to physically organize data in the destination.
    * `partition_fields.#.field_name` - Human-readable label that identifies the field name used to partition data.
    * `partition_fields.#.order` - Sequence in which MongoDB Atlas slices the collection data to create partitions. The resource expresses this sequence starting with zero.
### `source` - Ingestion Source of a Data Lake Pipeline.
  * `type` - Type of ingestion source of this Data Lake Pipeline.
  * `cluster_name` - Human-readable name that identifies the cluster.
  * `collection_name` - Human-readable name that identifies the collection.
  * `database_name` - Human-readable name that identifies the database.
  * `project_id` - Unique 24-hexadecimal character string that identifies the project.
  * `policyItemId` - Unique 24-hexadecimal character string that identifies a policy item.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines) Documentation for more information.