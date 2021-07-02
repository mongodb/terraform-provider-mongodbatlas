---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search index"
sidebar_current: "docs-mongodbatlas-search-index"
description: |-
Describes a Search Index.
---

# mongodbatlas_search_index

`mongodbatlas_search_index` describe all Projects. This represents projects that have been created.

> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.


## Example Usage

```hcl
data "mongodbatlas_search_index" "test" {
  index_id   = "<INDEX_ID"
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  
}
```

## Argument Reference

* `index_id` - (Required) The indexID of the search index you want to create.
* `project_id` - (Required) The ID of the organization or project you want to create the search index within.
* `cluster_name` - (Required) The name of the cluster where you want to create the search index within.



For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
