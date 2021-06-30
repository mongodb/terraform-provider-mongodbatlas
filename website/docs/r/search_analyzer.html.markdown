---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search analyzer"
sidebar_current: "docs-mongodbatlas-resource-search-analyzer"
description: |-
Provides a Search Analyzers resource.
---

# mongodbatlas_search_analyzer

`mongodbatlas_search_analyzer` creates new user-defined Atlas Search [analyzers](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) and replaces existing ones. If you have any user-defined analyzers for your Atlas Search [index](https://docs.atlas.mongodb.com/atlas-search/#std-label-fts-top-ref), they will be replaced by the new analyzer or analyzers specified in this request.

## Example Usage

```hcl
resource "mongodbatlas_search_analyzer" "test" {
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  
  search_analyzers {
    base_analyzer = "lucene.standard"
    ignore_case = yes
    max_token_length = 10
    name = "search_analyzer_1"
  }
}
```

## Argument Reference
* `project_id` - (Required) The ID of the organization or project you want to create the search analyzer within.
* `cluster_name` - (Required) The name of the cluster where you want to create the search analyzer within.

* `base_analyzer` - (Required) [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) on which the user-defined analyzer is based.
* `ignore_case` - Specify whether the index is case-sensitive.
* `max_token_length` - Longest text unit to analyze. Atlas Search excludes anything longer from the index.
* `name` - (Required) Name of the user-defined analyzer.
* `stem_exclusion_set` - Words to exclude from [stemming](https://en.wikipedia.org/wiki/Stemming) by the language analyzer.
* `stopwords` - Strings to ignore when creating the index.



For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
