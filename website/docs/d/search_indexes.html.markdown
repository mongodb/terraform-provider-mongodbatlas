---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search indexes"
sidebar_current: "docs-mongodbatlas-search-indexes"
description: |-
Describes a Search Indexes.
---

# mongodbatlas_search_indexes

`mongodbatlas_search_indexes` describe all Projects. This represents projects that have been created.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.


## Example Usage

```hcl
data "mongodbatlas_search_index" "test" {
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  database_name ="<DATABASE_NAME>"
  collection_name = "<COLLECTION_NAME>"
  
  page_num = 1
  items_per_page = 100
  
}
```

## Argument Reference

* `project_id` - (Required) The ID of the organization or project you want to get the search index within.
* `cluster_name` - (Required) The name of the cluster where you want to get the search index within.
* `database_name` - (Required) Name of the database containing the collection with one or more Atlas Search indexes.
* `collection_name` - (Required) Name of the collection with one or more Atlas Search indexes
* `total_count` - Represents the total of the search indexes

* `page_num` - Page number, starting with one, that Atlas returns of the total number of objects.
  
* `items_per_page` - Number of items that Atlas returns per page, up to a maximum of 500.
* `results` - A list where each represents a search index.



For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
