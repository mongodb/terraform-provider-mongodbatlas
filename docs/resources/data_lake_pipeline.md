# Resource: mongodbatlas_data_lake_pipeline

`mongodbatlas_data_lake_pipeline` provides a Data Lake Pipeline resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find `group_id` in the official documentation.

## Example Usages


```terraform
resource "mongodbatlas_project" "projectTest" {
  name   = "NAME OF THE PROJECT"
  org_id = "ORGANIZATION ID"
}

resource "mongodbatlas_advanced_cluster" "automated_backup_test" {
  project_id     = var.project_id
  name           = "automated-backup-test"
  cluster_type   = "REPLICASET"
  backup_enabled = true # enable cloud backup snapshots

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "GCP"
      region_name   = "US_EAST_4"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

resource "mongodbatlas_data_lake_pipeline" "pipeline" {
      project_id       = mongodbatlas_project.projectTest.project_id
      name       = "DataLakePipelineName"
      sink {
        type = "DLS"
        partition_fields {
            name = "access"
            order = 0
        }
      }

      source {
        type = "ON_DEMAND_CPS"
        cluster_name = mongodbatlas_advanced_cluster.automated_backup_test.name
        database_name = "sample_airbnb"
        collection_name = "listingsAndReviews"
      }

      transformations {
              field = "test"
              type  = "EXCLUDE"
      }

      transformations {
              field = "test22"
              type  = "EXCLUDE"
      }
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create a data lake pipeline.
* `name` - (Required) Name of the Atlas Data Lake Pipeline.

## Attributes Reference

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


## Import

Data Lake Pipeline can be imported using project ID, name of the data lake and name of the AWS s3 bucket, in the format `project_id`--`name`, e.g.

```
$ terraform import mongodbatlas_data_lake_pipeline.example 1112222b3bf99403840e8934--test-data-lake-pipeline-test
```

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines) Documentation for more information.