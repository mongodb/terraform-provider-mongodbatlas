---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search indexes"
sidebar_current: "docs-mongodbatlas-datasource-search-indexes"
description: |-
Describes a Search Indexes.
---

# Data Source: mongodbatlas_search_indexes

`mongodbatlas_search_indexes` describe all search indexes. This represents search indexes that have been created.

> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.


## Example Usage

```terraform
data "mongodbatlas_search_indexes" "test" {
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  database_name ="<DATABASE_NAME>"
  collection_name = "<COLLECTION_NAME>"
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the [project](https://docs.atlas.mongodb.com/organizations-projects/#std-label-projects) that contains the specified cluster.
* `cluster_name` - (Required) Name of the cluster containing the collection with one or more Atlas Search indexes.
* `database_name` - (Required) Name of the database containing the collection with one or more Atlas Search indexes.
* `collection_name` - (Required) Name of the collection with one or more Atlas Search indexes.

## Attributes Reference
* `total_count` - Represents the total of the search indexes
* `results` - A list where each represents a search index.

### Results

* `name` - Name of the index.
* `status` - Current status of the index.
* `analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when creating the index.
* `analyzers` - [Custom analyzers](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-custom-analyzers) to use in this index (this is an array of objects).
* `collection_name` - (Required) Name of the collection the index is on.
* `database` - (Required) Name of the database the collection is in.
* `mappings_dynamic` - Flag indicating whether the index uses dynamic or static mappings.
* `mappings_fields` - Object containing one or more field specifications.
* `search_analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when searching the index.
* `synonyms` - 	Synonyms mapping definition to use in this index.
* `synonyms.#.name` - Name of the [synonym mapping definition](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#std-label-synonyms-ref).
* `synonyms.#.source_collection` - Name of the source MongoDB collection for the synonyms.
* `synonyms.#.analyzer` - Name of the [analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use with this synonym mapping.




For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
