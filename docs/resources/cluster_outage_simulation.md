# Resource: mongodbatlas_cluster_outage_simulation

`mongodbatlas_cluster_outage_simulation` provides a Cluster Outage Simulation resource. For more details see https://www.mongodb.com/docs/atlas/tutorial/test-resilience/simulate-regional-outage/

Test Outage on Minority of Electable Nodes - Select fewer than half of your electable nodes. 

Test Outage on Majority of Electable Nodes - Select at least one more than half of your electable nodes and keep at least one electable node remaining. 

**IMPORTANT:** Test Outage on Majority of Electable Nodes will leave the Atlas cluster without a majority quorum. There will be no primary so write operations will not succeed, and reads will succeed only when configured with a suitable [readPreference](https://www.mongodb.com/docs/manual/core/read-preference/). To recover the majority quorum, you will have the option to manually reconfigure your cluster by adding new nodes to existing regions or adding new regions at the risk of losing recent writes, or end the simulation.   


-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** This resource cannot be updated.
~> **IMPORTANT:** An existing Cluster Outage Simulation cannot be imported as this resource does not support import operation.

## Example Usages


```terraform
resource "mongodbatlas_cluster_outage_simulation" "outage_simulation" {
  project_id = "64707f06c519c20c3a2b1b03"
  cluster_name = "Cluster0"
 	outage_filters {
     	cloud_provider = "AWS"
     	region_name = "US_EAST_1"
 	}

    outage_filters {
     	cloud_provider = "AWS"
     	region_name = "US_EAST_2"
 	}
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project that contains the cluster that is/will undergoing outage simulation.
* `cluster_name` - (Required) Name of the Atlas Cluster that is/will undergoing outage simulation.
* `outage_filters` - (Minimum one required) List of settings that specify the type of cluster outage simulation.
  * `cloud_provider` - (Required) The cloud provider of the region that undergoes the outage simulation. Following values are supported:
    * `AWS`
    * `GCP`
    * `AZURE`
  * `region_name` - (Required) The Atlas name of the region to undergo an outage simulation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `simulation_id` - Unique 24-hexadecimal character string that identifies the outage simulation.
* `start_request_date` - Date and time when MongoDB Cloud started the regional outage simulation.
* `outage_filters` - List of settings that specify the type of cluster outage simulation.
  * `type` - The type of cluster outage simulation. Following values are supported:
    * `REGION` - Simulates a cluster outage for a region
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
