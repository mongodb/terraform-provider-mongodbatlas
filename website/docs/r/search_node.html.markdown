---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search node"
sidebar_current: "docs-mongodbatlas-resource-search-node"
description: |-
Provides a Search Node resource.
---

# Resource: mongodbatlas_search_node

`mongodbatlas_search_node` provides a Search Node resource. The resource lets you create, edit and delete dedicated search nodes in a cluster.

-> **NOTE:** For details on supported cloud providers and existing limitations you can visit the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-nodes-for-workload-isolation).
-> **NOTE:** Only a single search node resource can be defined for each cluster.


## Example Usage

```terraform
resource "mongodbatlas_search_node" "test" {
  project_id = "PROJECT ID"
  cluster_name = "NAME OF CLUSTER"
  specs = [
    {
      instance_size = "S20_HIGHCPU_NVME"
      node_count = 2
    }
  ]
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `cluster_name` - (Required) Label that identifies the cluster to create search nodes for.
* `specs` - (Required) List of settings that configure the search nodes for your cluster. This list is currently limited to defining a single element. See [specs](#specs).

### Specs

TODO: add proper link here
* `instance_size` - (Required) Hardware specification for the search node instance sizes. The [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/) describes the valid values. More details can also be found in the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-tier).
* `node_count` - (Required) Number of search nodes in the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `state_name` - Human-readable label that indicates the current operating condition of this search node deployment.

## Import

Search node resource can be imported using the project ID and cluster name, in the format `PROJECT_ID-CLUSTER_NAME`, e.g.

```
$ terraform import mongodbatlas_search_node.test 650972848269185c55f40ca1-Cluster0
```
TODO: add proper link here
For more information see: [MongoDB Atlas API - Search Node](https://docs.atlas.mongodb.com/reference/api/) Documentation.
