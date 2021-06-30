---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search analyzers"
sidebar_current: "docs-mongodbatlas-resource-search-analyzers"
description: |-
Provides a Search Analyzers resource.
---

# mongodbatlas_search_analyzer

`mongodbatlas_search_analyzer` allows you to retrieve and edit Atlas Search [analyzers](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-custom-analyzers) and [index configurations](https://docs.atlas.mongodb.com/atlas-search/#std-label-fts-top-ref) for the specified cluster.

## Example Usage

```hcl
data "mongodbatlas_search_analyzer" "test" {
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  page_num = 1
  items_per_page = 100
  
}
```

## Argument Reference
* `project_id` - (Required) The ID of the organization or project you want to get the search analyzers within.
* `cluster_name` - (Required) The name of the cluster where you want to get the search analyzers within.
* `total_count` - Represents the total of the search indexes

* `page_num` - Page number, starting with one, that Atlas returns of the total number of objects.

* `items_per_page` - Number of items that Atlas returns per page, up to a maximum of 500.
* `results` - A list where each represents a search index.



For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
