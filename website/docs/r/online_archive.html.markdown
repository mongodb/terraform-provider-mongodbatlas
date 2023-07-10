---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_online_archive"
sidebar_current: "docs-mongodbatlas-resource-online-archive"
description: |-
    Provides a Online Archive resource for creation, update, and delete
---

# Resource: mongodbatlas_online_archive

`mongodbatlas_online_archive` resource provides access to create, edit, pause and resume an online archive for a collection. 

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** The collection must exists before performing an online archive.

~> **IMPORTANT:** There are fields that are immutable after creation, i.e if `date_field` value does not exist in the collection, the online archive state will be pending forever, and this field cannot be updated, that means a destroy is required, known error `ONLINE_ARCHIVE_CANNOT_MODIFY_FIELD`

## Example Usages
```terraform
resource "mongodbatlas_online_archive" "test" {
    project_id   = var.project_id
    cluster_name = var.cluster_name
    coll_name    = var.collection_name
    db_name      = var.database_name

    partition_fields {
        field_name = "firstName"
        order = 0
    }

    partition_fields {
        field_name = "lastName"
        order = 1
    }

    criteria {
        type = "DATE"
        date_field = "created"
        expire_after_days = 5
    }

    schedule {
        type = "DAILY"
        end_hour = 1
        end_minute = 1
        start_hour = 1
        start_minute = 1
    }
}
```

For custom criteria example

```terraform
resource "mongodbatlas_online_archive" "test" {
    project_id   = var.project_id
    cluster_name = var.cluster_name
    coll_name    = var.collection_name
    db_name      = var.database_name

    partition_fields {
        field_name = "firstName"
        order      = 0 
    }

    partitions_fields {
        field_name = "secondName"
        order      = 1 
    }

    criteria {
        type  = "CUSTOM"
        query =  "{ \"department\": \"engineering\" }"
    }

}

```

## Argument Reference
* `project_id`       -  (Required) The unique ID for the project
* `cluster_name`     -  (Required) Name of the cluster that contains the collection.
* `db_name`          -  (Required) Name of the database that contains the collection.
* `coll_name`        -  (Required) Name of the collection.
* `collection_type`  -  Classification of MongoDB database collection that you want to return, "TIMESERIES" or "STANDARD". Default is "STANDARD". 
* `criteria`         -  (Required) Criteria to use for archiving data.
* `partition_fields` -  (Recommended) Fields to use to partition data. You can specify up to two frequently queried fields to use for partitioning data. Note that queries that don’t contain the specified fields will require a full collection scan of all archived documents, which will take longer and increase your costs. To learn more about how partition improves query performance, see [Data Structure in S3](https://docs.mongodb.com/datalake/admin/optimize-query-performance/#data-structure-in-s3). The value of a partition field can be up to a maximum of 700 characters. Documents with values exceeding 700 characters are not archived.
* `paused`           - (Optional) State of the online archive. This is required for pausing an active or resume a paused online archive. The resume request will fail if the collection has another active online archive.

### Criteria

There are two types of criteria, `DATE` to select documents for archiving based on a date and
`CUSTOM` to select documents for archiving based on a custom JSON query.

* `type`          - Type of criteria (DATE, CUSTOM)

The following fields are required for criteria type `DATE`

* `date_field`   - Indexed database parameter that stores the date that determines when data moves to the online archive. MongoDB Cloud archives the data when the current date exceeds the date in this database parameter plus the number of days specified through the expireAfterDays parameter.
* `date_format`   - Syntax used to write the date after which data moves to the online archive. Date can be expressed as ISO 8601 or Epoch timestamps. The Epoch timestamp can be expressed as nanoseconds, milliseconds, or seconds. You must set `type` to `DATE` if `collectionType` is `TIMESERIES`. Valid values:  ISODATE (default), EPOCH_SECONDS, EPOCH_MILLIS, EPOCH_NANOSECONDS.
* `expire_after_days` - Number of days after the value in the criteria.dateField when MongoDB Cloud archives data in the specified cluster.

The only field required for criteria type `CUSTOM`

* `query` - JSON query to use to select documents for archiving. Atlas uses the specified query with the db.collection.find(query) command. The empty document {} to return all documents is not supported.

### Schedule

* `type`          - Type of schedule (`DEFAULT`, `DAILY`, `MONTHLY`, `WEEKLY`).
* `start_hour`    - Hour of the day when the when the scheduled window to run one online archive starts.  
* `end_hour`      - Hour of the day when the scheduled window to run one online archive ends.
* `start_minute`   - Minute of the hour when the scheduled window to run one online archive starts.
* `end_minute`     - Minute of the hour when the scheduled window to run one online archive ends.
* `day_of_month`   - Day of the month when the scheduled archive starts. This field should be provided only when schedule `type` is `MONTHLY`.
* `day_of_week`     - Day of the week when the scheduled archive starts. The week starts with Monday (1) and ends with Sunday (7). This field should be provided only when schedule `type` is `WEEKLY`.

### Partition
* `field_name` - Human-readable label that identifies the parameter that MongoDB Cloud uses to partition data. To specify a nested parameter, use the dot notation.
* `order` - Sequence in which MongoDB Cloud slices the collection data to create partitions. The resource expresses this sequence starting with zero. The value of the `criteria.dateField` parameter defaults as the first item in the partition sequence.
* `field_type` - Data type of the parameter that that MongoDB Cloud uses to partition data. Partition parameters of type UUID must be of binary subtype 4. MongoDB Cloud skips partition parameters of type UUID with subtype 3. Valid values: `date`, `int`, `long`, `objectId`, `string`, `uuid`.

## Attributes Reference
* `archive_id` - ID of the online archive.
* `state`    - Status of the online archive. Valid values are: Pending, Archiving, Idle, Pausing, Paused, Orphaned and Deleted

## Import 

```bash
terraform import mongodbatlas_online_archive.users_archive <project_id>-<cluster_name>-<archive_id>
```

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/online-archive-create-one/) Documentation for more information.
