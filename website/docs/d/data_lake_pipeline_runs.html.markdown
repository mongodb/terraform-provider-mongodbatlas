---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: data_lake_pipeline_runs"
sidebar_current: "docs-mongodbatlas-datasource-data-lake-pipeline-runs"
description: |-
    Describes Data Lake Pipeline Runs.
---

# Data Source: mongodbatlas_data_lake_pipeline_runs

`mongodbatlas_data_lake_pipeline_run` describe Data Lake Pipeline Runs.


-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_data_lake_pipeline" "pipeline" {
      //assuming we've already setup project and cluster in another block
      project_id       = mongodbatlas_project.projectTest.project_id
      name             = "DataLakePipelineName"
      sink {
        type = "DLS"
        partition_fields {
            name = "access"
            order = 0
        }
      }
      source {
        type = "ON_DEMAND_CPS"
        cluster_name = mongodbatlas_cluster.clusterTest.name
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

data "mongodbatlas_data_lake_pipeline_runs" "test" {
  project_id       = mongodbatlas_project.projectTest.project_id
  pipeline_name    = mongodbatlas_data_lake_pipeline.pipeline.name
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `pipeline_name` - (Required) Human-readable label that identifies the Data Lake Pipeline.

## Attributes Reference
* `results` - A list where each represents a Data Lake Pipeline Run.

### Data Lake Pipeline Run

* `id` - Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
* `pipeline_run_id` - Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
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

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/listPipelineRuns) Documentation for more information.