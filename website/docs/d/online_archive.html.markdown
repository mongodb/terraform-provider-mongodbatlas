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
* `date_field`    - Name of an already indexed date field from the documents. Data is archived when the current date is greater than the value of the date field specified here plus the number of days specified via the `expire_after_days` parameter.
* `date_format`   - the date format. Valid values:  ISODATE (default), EPOCH_SECONDS, EPOCH_MILLIS, EPOCH_NANOSECONDS
* `expire_after_days` - Number of days that specifies the age limit for the data in the live Atlas cluster. Data is archived when the current date is greater than the value of the date field specified via the `date_field` parameter plus the number of days specified here.
* `query` - JSON query to use to select documents for archiving. Atlas uses the specified query with the db.collection.find(query) command. The empty document {} to return all documents is not supported.

### Schedule

* `type`          - Type of schedule (`DEFAULT`, `DAILY`, `MONTHLY`, `WEEKLY`).
* `start_hour`    - Hour of the day when the when the scheduled window to run one online archive starts.  
* `end_hour`      - Hour of the day when the scheduled window to run one online archive ends.
* `start_minute`   - Minute of the hour when the scheduled window to run one online archive starts.
* `end_minute`     - Minute of the hour when the scheduled window to run one online archive ends.
* `day_of_month`   - Day of the month when the scheduled archive starts.
* `day_of_week`     - Day of the week when the scheduled archive starts. The week starts with Monday (1) and ends with Sunday (7).

### Partition
* `field_name` - Name of the field. To specify a nested field, use the dot notation.
* `order` - Position of the field in the partition. Value can be: 0,1,2
By default, the date field specified in the criteria.dateField parameter is in the first position of the partition.
* `field_type` - (Optional) type of the partition field

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/online-archive-get-one/) Documentation for more information.


