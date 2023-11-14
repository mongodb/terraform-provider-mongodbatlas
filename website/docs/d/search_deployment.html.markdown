---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search deployment"
sidebar_current: "docs-mongodbatlas-datasource-search-deployment"
description: |-
Describes a Search Deployment.
---

# Data Source: mongodbatlas_search_deployment

`mongodbatlas_search_deployment` describes a search node deployment.

## Example Usage

```terraform
data "mongodbatlas_search_deployment" "test" {
    project_id = "<PROJECT_ID>"
    cluster_name = "<CLUSTER_NAME>"
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier for the [project](https://docs.atlas.mongodb.com/organizations-projects/#std-label-projects) that contains the specified cluster.
* `cluster_name` - (Required) The name of the cluster containing a search node deployment.

## Attributes Reference

* `specs` - List of settings that configure the search nodes for your cluster. See [specs](#specs).
* `state_name` - Human-readable label that indicates the current operating condition of this search node deployment.

### Specs
TODO: add proper link here
* `instance_size` - (Required) Hardware specification for the search node instance sizes. The [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/) describes the valid values. More details can also be found in the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-tier).
* `node_count` - (Required) Number of search nodes in the cluster.


TODO: add proper link here
For more information see: [MongoDB Atlas API - Search Node](https://docs.atlas.mongodb.com/reference/api/) Documentation.
