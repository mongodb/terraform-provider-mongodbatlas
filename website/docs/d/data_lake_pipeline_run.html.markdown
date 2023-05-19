---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: data_lake_pipeline_run"
sidebar_current: "docs-mongodbatlas-datasource-data-lake-pipeline-run"
description: |-
    Describes a Data Lake Pipeline Run.
---

# Data Source: mongodbatlas_data_lake_pipeline_run

`mongodbatlas_data_lake_pipeline_run` describe a Data Lake Pipeline Run.


-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
data "mongodbatlas_data_lake_pipeline_run" "test" {
  project_id = "PROJECT ID"
  name = "DATA LAKE PIPELINE NAME"
  pipeline_run_id = "DATA LAKE PIPELINE RUN ID"
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `name` - (Required) Human-readable label that identifies the Data Lake Pipeline.
* `pipeline_run_id` - (Required) Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
* `created_date` - Timestamp that indicates when the pipeline run was created.
* `last_updated_date` - Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
* `state` - State of the pipeline run.
* `dataset_name` - Human-readable label that identifies the dataset that Atlas generates during this pipeline run. 
* `phase` - Processing phase of the Data Lake Pipeline.
* `pipeline_id` - Unique 24-hexadecimal character string that identifies a Data Lake Pipeline.
* `snapshot_id` - Unique 24-hexadecimal character string that identifies the snapshot of a cluster.
* `backup_frequency_type` - Backup schedule interval of the Data Lake Pipeline.
* `stats` - Runtime statistics for this Data Lake Pipeline run.
  * `bytes_exported` - Total data size in bytes exported for this pipeline run.
  * `num_docs` - Number of docs ingested for a this pipeline run.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/getPipelineRun) Documentation for more information.