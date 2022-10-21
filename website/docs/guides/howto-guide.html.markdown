---
layout: "mongodbatlas"
page_title: "MongoDB Atlas Provider How to Guide"
sidebar_current: "docs-mongodbatlas-guides-how-to-guide"
description: |-
MongoDB Atlas Provider : How to Guide
---

# MongoDB Atlas Provider: How to Guide

The Terraform MongoDB Atlas Provider guide to perform common tasks with the provider.

Document for users how to get a pre-existing container id

1) Empty project

2) empty state file

3) applied curl to build cluster

4) terraform apply to retrieve container id

Validated by testing following scenario container_id populated in state after first apply ...

1) create cluster using curl command to simulate non terraform based cluster

```
curl --user "pub:priv" --digest \
--header "Content-Type: application/json" \
--include \
--request POST "https://cloud.mongodb.com/api/atlas/v1.0/groups/grpid/clusters?pretty=true" \
--data '
{
  "name": "SingleRegionCluster",
  "diskSizeGB": 100,
  "numShards": 1,
  "providerSettings": {
    "providerName": "AWS",
    "instanceSizeName": "M40",
    "regionName": "US_EAST_1"
  },
  "clusterType": "REPLICASET",
  "replicationFactor": 3,
  "replicationSpecs": [
    {
      "numShards": 1,
      "regionsConfig": {
        "US_EAST_1": {
          "analyticsNodes": 0,
          "electableNodes": 3,
          "priority": 7,
          "readOnlyNodes": 0
        }
      },
      "zoneName": "Zone 1"
    }
  ],
  "backupEnabled": false,
  "providerBackupEnabled": true,
  "autoScaling": {
    "diskGBEnabled": true
  }
}'
```

 

2) Apply this terraform to read information from datasources

 
```
data "mongodbatlas_cluster" "admin" { 
  name = "SingleRegionCluster" 
  project_id = local.mongodbatlas_project_id 
}

data "mongodbatlas_network_container" "admin" { 
  project_id = local.mongodbatlas_project_id 
  container_id = data.mongodbatlas_cluster.admin.container_id 
}

output "container" { 
  value = data.mongodbatlas_network_container.admin.container_id 
}

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

container = "62ffe4ecb79e2e007c375935"
```
 

Using versions
- darwin_amd64
- provider registry.terraform.io/hashicorp/aws v4.26.0
- provider registry.terraform.io/mongodb/mongodbatlas v1.4.3


### Helpful Links

* [Report bugs](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)

* [Request Features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)

* [Contact Support](https://docs.atlas.mongodb.com/support/) covered by MongoDB Atlas support plans, Developer and above.
  