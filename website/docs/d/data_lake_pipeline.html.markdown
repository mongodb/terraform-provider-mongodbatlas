---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: data_lake"
sidebar_current: "docs-mongodbatlas-resource-data-lake"
description: |-
    Describe a Data Lake Pipeline.
---

# Data Source: mongodbatlas_data_lake_pipeline

`mongodbatlas_data_lake_pipeline` describe a Data Lake Pipeline.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** All arguments including the password will be stored in the raw state as plain-text. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html)

## Example Usages


```terraform
resource "mongodbatlas_project" "projectTest" {
  name   = "NAME OF THE PROJECT"
  org_id = "ORGANIZATION ID"
}

resource "mongodbatlas_cluster" "automated_backup_test" {
    project_id   = "63f4d4a47baeac59406dc131"
    name         = "automated-backup-test"

    provider_name               = "GCP"
    provider_region_name        = "US_EAST_4"
    provider_instance_size_name = "M30"
    cloud_backup                = true   // enable cloud backup snapshots
    mongo_db_major_version      = "4.4"
    disk_size_gb = "350"
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
        cluster_name = mongodbatlas_cluster.automated_backup_test.name
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

data "mongodbatlas_data_lake_pipeline" "pipelineDataSource" {
  project_id       = mongodbatlas_data_lake_pipeline.pipeline.project_id
  name             = mongodbatlas_data_lake_pipeline.pipeline.name
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
* `sink` - Ingestion destination of a Data Lake Pipeline.
  * `sink.0.type` - Type of ingestion destination of this Data Lake Pipeline.
  * `sink.0.provider` - Target cloud provider for this Data Lake Pipeline.
  * `sink.0.region` - Target cloud provider region for this Data Lake Pipeline. [Supported cloud provider regions](https://www.mongodb.com/docs/datalake/limitations).
  * `sink.0.partition_fields` - Ordered fields used to physically organize data in the destination.
    * `sink.0.partition_fields.#.name` - Human-readable label that identifies the field name used to partition data.
    * `sink.0.partition_fields.#.order` - Sequence in which MongoDB Cloud slices the collection data to create partitions. The resource expresses this sequence starting with zero.
* `source` - Ingestion Source of a Data Lake Pipeline.
  * `source.0.type` - Type of ingestion source of this Data Lake Pipeline.
  * `source.0.cluster_name` - Human-readable name that identifies the cluster.
  * `source.0.collection_name` - Human-readable name that identifies the collection.
  * `source.0.database_name` - Human-readable name that identifies the database.
  * `source.0.project_id` - Unique 24-hexadecimal character string that identifies the project.
  * `source.0.policyItemId` - Unique 24-hexadecimal character string that identifies a policy item.
* `transformations` - Fields to be excluded for this Data Lake Pipeline.
  * `transformations.#.field` - Key in the document.
  * `transformations.#.type` - Type of transformation applied during the export of the namespace in a Data Lake Pipeline.
* `snapshots` - List of backup snapshots that you can use to trigger an on demand pipeline run.
  * `snapshots.#.id` - Unique 24-hexadecimal digit string that identifies the snapshot.
  * `snapshots.#.provider` - Human-readable label that identifies the cloud provider that stores this snapshot.
  * `snapshots.#.created_at` - Date and time when MongoDB Cloud took the snapshot.
  * `snapshots.#.expires_at` - Date and time when MongoDB Cloud deletes the snapshot.
  * `snapshots.#.frequency_type` - Human-readable label that identifies how often this snapshot triggers.
  * `snapshots.#.master_key` - Unique string that identifies the Amazon Web Services (AWS) Key Management Service (KMS) Customer Master Key (CMK) used to encrypt the snapshot.
  * `snapshots.#.mongod_version` - Version of the MongoDB host that this snapshot backs up.
  * `snapshots.#.replica_set_name` - Human-readable label that identifies the replica set from which MongoDB Cloud took this snapshot.
  * `snapshots.#.type` - Human-readable label that categorizes the cluster as a replica set or sharded cluster.
  * `snapshots.#.snapshot_type` - Human-readable label that identifies when this snapshot triggers.
  * `snapshots.#.status` - Human-readable label that indicates the stage of the backup process for this snapshot.
  * `snapshots.#.size` - List of backup snapshots that you can use to trigger an on demand pipeline run.
  * `snapshots.#.copy_region` - List that identifies the regions to which MongoDB Cloud copies the snapshot.
  * `snapshots.#.policies` - List that contains unique identifiers for the policy items.
* `ingestion_schedules` - List of backup schedule policy items that you can use as a Data Lake Pipeline source.
  * `ingestion_schedules.#.id` - Unique 24-hexadecimal digit string that identifies this backup policy item.
  * `ingestion_schedules.#.frequency_type` - Human-readable label that identifies the frequency type associated with the backup policy.
  * `ingestion_schedules.#.frequency_interval` - Number that indicates the frequency interval for a set of snapshots.
  * `ingestion_schedules.#.retention_unit` - Unit of time in which MongoDB Cloud measures snapshot retention.
  * `ingestion_schedules.#.retention_value` - Duration in days, weeks, or months that MongoDB Cloud retains the snapshot. 

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines) Documentation for more information.