---
page_title: "How-To Guide"
---

# MongoDB Atlas Provider: How-To Guide

The Terraform MongoDB Atlas Provider guide to perform common tasks with the provider.

##How to Get A Pre-existing Container ID

The following is an end to end example of how to get an existing container id. 

1) Start with an empty project

2) Empty state file

3) Apply a curl command to build cluster

4) Run `terraform apply` to retrieve the container id

The following illustrates step 3 and 4 above, assuming 1 & 2 were true:

1) Create a cluster using a curl command to simulate non-Terraform created cluster.  This will also create a container.  

```
curl --user "pub:priv" --digest \
--header "Content-Type: application/json" \
--include \
--request POST "https://cloud.mongodb.com/api/atlas/v1.0/groups/grpid/clusters?pretty=true" \
--data '
{
  "name": "SingleRegionCluster",
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
  "autoScaling": {
    "diskGBEnabled": true
  }
}'
```

 

2) Then apply this Terraform config to then read the information from the appropriate Data Sources and output the container id.  

 
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
 

This example was tested using versions:
- darwin_amd64
- provider registry.terraform.io/hashicorp/aws v4.26.0
- provider registry.terraform.io/mongodb/mongodbatlas v1.4.3


### Helpful Links

* [Report bugs](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)

* [Request Features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)

* [Contact Support](https://docs.atlas.mongodb.com/support/) covered by MongoDB Atlas support plans, Developer and above.
  