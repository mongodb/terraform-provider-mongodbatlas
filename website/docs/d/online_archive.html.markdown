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
* `state`    - Status of the online archive. Valid values are: Pending, Archiving, Idle, Pausing, Paused, Orphaned and Deleted

### Criteria
* `type`          - Type of criteria (DATE, CUSTOM)
* `date_field`   - Indexed database parameter that stores the date that determines when data moves to the online archive. MongoDB Cloud archives the data when the current date exceeds the date in this database parameter plus the number of days specified through the expireAfterDays parameter. Set this parameter when `type` is `DATE`.
* `date_format`   - Syntax used to write the date after which data moves to the online archive. Date can be expressed as ISO 8601 or Epoch timestamps. The Epoch timestamp can be expressed as nanoseconds, milliseconds, or seconds. Set this parameter when `type` is `DATE`. You must set `type` to `DATE` if `collectionType` is `TIMESERIES`. Valid values:  ISODATE (default), EPOCH_SECONDS, EPOCH_MILLIS, EPOCH_NANOSECONDS.
* `expire_after_days` - Number of days after the value in the criteria.dateField when MongoDB Cloud archives data in the specified cluster. Set this parameter when `type` is `DATE`.
* `query` - JSON query to use to select documents for archiving. Atlas uses the specified query with the db.collection.find(query) command. The empty document {} to return all documents is not supported. Set this parameter when `type` is `CUSTOM`.

### Schedule

* `type`          - Type of schedule. Valid values: `DEFAULT`, `DAILY`, `MONTHLY`, `WEEKLY`.
* `start_hour`    - Hour of the day when the when the scheduled window to run one online archive starts.  
* `end_hour`      - Hour of the day when the scheduled window to run one online archive ends.
* `start_minute`   - Minute of the hour when the scheduled window to run one online archive starts.
* `end_minute`     - Minute of the hour when the scheduled window to run one online archive ends.
* `day_of_month`   - Day of the month when the scheduled archive starts. Set this parameter when `type` is `MONTHLY`.
* `day_of_week`     - Day of the week when the scheduled archive starts. The week starts with Monday (1) and ends with Sunday (7).Set this parameter when `type` is `WEEKLY`.

### Partition
* `field_name` - Human-readable label that identifies the parameter that MongoDB Cloud uses to partition data. To specify a nested parameter, use the dot notation.
* `order` - Sequence in which MongoDB Cloud slices the collection data to create partitions. The resource expresses this sequence starting with zero. The value of the `criteria.dateField` parameter defaults as the first item in the partition sequence.
* `field_type` - Data type of the parameter that that MongoDB Cloud uses to partition data. Partition parameters of type UUID must be of binary subtype 4. MongoDB Cloud skips partition parameters of type UUID with subtype 3. Valid values: `date`, `int`, `long`, `objectId`, `string`, `uuid`.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/online-archive-get-one/) Documentation for more information.


