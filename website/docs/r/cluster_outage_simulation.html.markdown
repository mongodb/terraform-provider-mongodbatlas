---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cluster_outage_simulation"
sidebar_current: "docs-mongodbatlas-resource-federated-database-instance"
description: |-
    Provides a Cluster Outage Simulation resource.
---

# Resource: mongodbatlas_cluster_outage_simulation

`mongodbatlas_cluster_outage_simulation` provides a Cluster Outage Simulation resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** This resource cannot be updated.
~> **IMPORTANT:** An existing Cluster Outage Simulation cannot be imported as this resource does not support import operation.

## Example Usages


```terraform
resource "mongodbatlas_cluster_outage_simulation" "outage_simulation" {
  project_id = "PROJECT ID"
  cluster_name = "NAME OF THE CLUSTER THAT WILL UNDERGO OUTAGE SIMULATION"
 	outage_filters {
     	cloud_provider = "NAME OF CLOUD PROVIDER OF THE REGION"
     	region_name = "ATLAS REGION NAME"
        type = "REGION"
 	}

    outage_filters {
     	cloud_provider = "NAME OF CLOUD PROVIDER OF THE REGION"
     	region_name = "ATLAS REGION NAME"
        type = "REGION"
 	}
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project that contains the cluster that is/will undergoing outage simulation.
* `cluster_name` - (Required) Name of the Atlas Cluster that is/will undergoing outage simulation.
* `outage_filters` - (Minimum one required) List of settings that specify the type of cluster outage simulation.
  * `outage_filters.0.cloud_provider` - (Required) The cloud provider of the region that undergoes the outage simulation. Following values are supported:
    * `AWS`
    * `GCP`
    * `AZURE`
  * `outage_filters.0.region_name` - (Required) The Atlas name of the region to undergo an outage simulation.
  * `outage_filters.0.type` - (Required) The type of cluster outage to simulate. Following values are supported:
    * `REGION` (Simulates a cluster outage for a region)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `simulation_id` - Unique 24-hexadecimal character string that identifies the outage simulation.
* `start_request_date` - Date and time when MongoDB Cloud started the regional outage simulation.
* `state` - Current phase of the outage simulation:
  * `START_REQUESTED` - User has requested cluster outage simulation.
  * `STARTING` - MongoDB Cloud is starting cluster outage simulation.
  * `SIMULATING` - MongoDB Cloud is simulating cluster outage.
  * `RECOVERY_REQUESTED` - User has requested recovery from the simulated outage.
  * `RECOVERING` - MongoDB Cloud is recovering the cluster from the simulated outage.
  * `COMPLETE` - MongoDB Cloud has completed the cluster outage simulation.

## Import

The `mongodbatlas_cluster_outage_simulation` resource does not support import operation.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation) Documentation for more information.