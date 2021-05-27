---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_online_archive"
sidebar_current: "docs-mongodbatlas-resource-online-archive"
description: |-
    Provides a Online Archive resource for creation, update, and delete
---

# mongodbatlas_online_archive

`mongodbatlas_online_archive` resource provides access to create, edit, pause and resume an online archive for a collection. 

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** The collection must exists before performing an online archive.

~> **IMPORTANT:** There are fields that are immutable after creation, i.e if `date_field` value does not exist in the collection, the online archive state will be pending forever, and this field cannot be updated, that means a destroy is required, known error `ONLINE_ARCHIVE_CANNOT_MODIFY_FIELD`

## Example Usages
```hcl
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

}
```

For custom criteria example

```hcl
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
* `criteria`         -  (Required) Criteria to use for archiving data.
* `partition_fields` -  (Recommended) Fields to use to partition data. You can specify up to two frequently queried fields to use for partitioning data. Note that queries that don’t contain the specified fields will require a full collection scan of all archived documents, which will take longer and increase your costs. To learn more about how partition improves query performance, see [Data Structure in S3](https://docs.mongodb.com/datalake/admin/optimize-query-performance/#data-structure-in-s3). The value of a partition field can be up to a maximum of 700 characters. Documents with values exceeding 700 characters are not archived.
* `paused`           - (Optional) State of the online archive. This is required for pausing an active or resume a paused online archive. The resume request will fail if the collection has another active online archive.

### Criteria details

There are two types of criteria, `DATE` to select documents for archiving based on a date and
`CUSTOM` to select documents for archiving based on a custom JSON query.

* `criteria.type`          - Type of criteria (DATE, CUSTOM)

The following fields are required for criteria type `DATE`

* `criteria.date_field`    - Name of an already indexed date field from the documents. Data is archived when the current date is greater than the value of the date field specified here plus the number of days specified via the `expire_after_days` parameter.
* `criteria.date_format`   - the date format. Valid values:  ISODATE (default), EPOCH_SECONDS, EPOCH_MILLIS, EPOCH_NANOSECONDS
* `criteria.expire_after_days` - Number of days that specifies the age limit for the data in the live Atlas cluster. Data is archived when the current date is greater than the value of the date field specified via the `date_field` parameter plus the number of days specified here.

The only field required for criteria type `CUSTOM`

* `criteria.query` - JSON query to use to select documents for archiving. Atlas uses the specified query with the db.collection.find(query) command. The empty document {} to return all documents is not supported.

### Partition fields details
* `partition_fields.field_name` - (Required) Name of the field. To specify a nested field, use the dot notation.
* `partition_fields.order` - (Required) Position of the field in the partition. Value can be: 0,1,2
By default, the date field specified in the criteria.dateField parameter is in the first position of the partition.
* `partitio_fields.field_type` - (Optional) type of the partition field

## Attributes Reference
* `archive_id` - ID of the online archive.
* `state`    - Status of the online archive. Valid values are: Pending, Archiving, Idle, Pausing, Paused, Orphaned and Deleted

## Import 

```bash
terraform import mongodbatlas_online_archive.users_archive <project_id>-<cluster_name>-<archive_id>
```

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/online-archive-create-one/) Documentation for more information.
