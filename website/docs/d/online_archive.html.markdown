---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: online_archive"
sidebar_current: "docs-mongodbatlas-datasource-online-archive"
description: |-
    Describes an Online Archive
---

# Data Source: mongodbatlas_online_archive

`mongodbatlas_online_archive` describes an Online Archive

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.


## Example Usage

```terraform 

data "mongodbatlas_online_archive" "test" {
    project_id   = var.project_id
    cluster_name = var.cluster_name
    archive_id     = "5ebad3c1fe9c0ab8d37d61e1"
}
```

## Argument Reference

* `project_id`    - (Required) The unique ID for the project.
* `cluster_name`  - (Required) Name of the cluster that contains the collection.
* `archive_id`      - (Required) ID of the online archive.

## Attributes reference
* `db_name`          -  Name of the database that contains the collection.
* `coll_name`        -  Name of the collection.
* `collection_type`  -  Classification of MongoDB database collection that you want to return, "TIMESERIES" or "STANDARD". Default is "STANDARD". 
* `criteria`         -  Criteria to use for archiving data.
* `criteria.type`          - Type of criteria (DATE, CUSTOM)
* `criteria.date_field`    - Name of an already indexed date field from the documents. Data is archived when the current date is greater than the value of the date field specified here plus the number of days specified via the `expire_after_days` parameter.
* `criteria.date_format`   - the date format. Valid values:  ISODATE (default), EPOCH_SECONDS, EPOCH_MILLIS, EPOCH_NANOSECONDS
* `criteria.expire_after_days` - Number of days that specifies the age limit for the data in the live Atlas cluster.
* `criteria.query` - JSON query to use to select documents for archiving. Only for `CUSTOM` type
* `partition_fields` -  Fields to use to partition data.
* `partition_fields.field_name` - Name of the field. To specify a nested field, use the dot notation.
* `partition_fields.order` - Position of the field in the partition. Value can be: 0,1,2
By default, the date field specified in the criteria.dateField parameter is in the first position of the partition.
* `partitio_fields.field_type` - Type of the partition field
* `state`    - Status of the online archive. Valid values are: Pending, Archiving, Idle, Pausing, Paused, Orphaned and Deleted

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/online-archive-get-one/) Documentation for more information.


