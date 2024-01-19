---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search index"
sidebar_current: "docs-mongodbatlas-datasource-search-index"
description: |-
Describes a Search Index.
---

# Data Source: mongodbatlas_search_index

`mongodbatlas_search_index` describe a single search indexes. This represents a single search index that have been created.

> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.


## Example Usage

```terraform
data "mongodbatlas_search_index" "test" {
  index_id   = "<INDEX_ID"
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  
}
```

## Argument Reference

* `index_id` - (Required) The unique identifier of the Atlas Search index. Use the `mongodbatlas_search_indexes`datasource to find the IDs of all Atlas Search indexes.
* `project_id` - (Required) The unique identifier for the [project](https://docs.atlas.mongodb.com/organizations-projects/#std-label-projects) that contains the specified cluster.
* `cluster_name` - (Required) The name of the cluster containing the collection with one or more Atlas Search indexes.

## Attributes Reference

* `name` - Name of the index.
* `status` - Current status of the index.
* `analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when creating the index.
* `analyzers` - [Custom analyzers](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-custom-analyzers) to use in this index (this is an array of objects).
* `collection_name` - Name of the collection the index is on.
* `database` - Name of the database the collection is in.
* `mappings_dynamic` - Flag indicating whether the index uses dynamic or static mappings.
* `mappings_fields` - Object containing one or more field specifications.
* `search_analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when searching the index.
* `synonyms` - 	Synonyms mapping definition to use in this index.
* `synonyms.#.name` - Name of the [synonym mapping definition](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#std-label-synonyms-ref).
* `synonyms.#.source_collection` - Name of the source MongoDB collection for the synonyms.
* `synonyms.#.analyzer` - Name of the [analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use with this synonym mapping. 




For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
