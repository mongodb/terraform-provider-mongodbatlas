---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cluster"
sidebar_current: "docs-mongodbatlas-datasource-clusters"
description: |-
    Describe all Clusters in Project.
---

# mongodb_atlas_cluster

`mongodb_atlas_cluster` describes all Clusters by the provided project_id. The data source requires your Project ID.

~> **IMPORTANT:** Changes to cluster configurations can affect costs. Before making changes, please see [Billing](https://docs.atlas.mongodb.com/billing/).

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** If your Atlas project contains a custom role that uses actions introduced in a specific MongoDB version, you cannot create a cluster with a MongoDB version less than that version unless you delete the custom role.

## Example Usage

```hcl
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "cluster-test"
  disk_size_gb = 100
  num_shards   = 1

  replication_factor           = 3
  backup_enabled               = true
  auto_scaling_disk_gb_enabled = true

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_disk_iops          = 300
  provider_encrypt_ebs_volume = false
  provider_instance_size_name = "M40"
  provider_region_name        = "US_EAST_1"
}

data "mongodbatlas_clusters" "test" {
	project_id = mongodbatlas_cluster.test.project_id // To get dependency.
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get the clusters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The cluster ID.
* `results` - A list where each represents a Cluster. See [Cluster](#cluster) below for more details.

### Cluster

* `name` - (Required) Name of the cluster as it appears in Atlas. Once the cluster is created, its name cannot be changed.
*  `mongo_db_version` - Version of MongoDB the cluster runs, in `major-version`.`minor-version` format.
* `mongo_uri` - Base connection string for the cluster. Atlas only displays this field after the cluster is operational, not while it builds the cluster.
* `mongo_uri_updated` - Lists when the connection string was last updated. The connection string changes, for example, if you change a replica set to a sharded cluster.
* `mongo_uri_with_options` - connection string for connecting to the Atlas cluster. Includes the replicaSet, ssl, and authSource query parameters in the connection string with values appropriate for the cluster.

    To review the connection string format, see the connection string format documentation. To add MongoDB users to a Atlas project, see Configure MongoDB Users.

    Atlas only displays this field after the cluster is operational, not while it builds the cluster.
* `paused` - Flag that indicates whether the cluster is paused or not.
* `srv_address` - Connection string for connecting to the Atlas cluster. The +srv modifier forces the connection to use TLS/SSL. See the mongoURI for additional options.
* `state_name` - Current state of the cluster. The possible states are:
    - IDLE
    - CREATING
    - UPDATING
    - DELETING
    - DELETED
    - REPAIRING
* `auto_scaling_disk_gb_enabled` - Specifies whether disk auto-scaling is enabled. The default is true.
    - Set to `true` to enable disk auto-scaling.
    - Set to `false` to disable disk auto-scaling.

* `backup_enabled` - Set to true to enable Atlas continuous backups for the cluster.

    Set to false to disable continuous backups for the cluster. Atlas deletes any stored snapshots. See the continuous backup Snapshot Schedule for more information.

    You cannot enable continuous backups if you have an existing cluster in the project with Cloud Provider Snapshots enabled.

    The default value is false.
* `bi_connector` - Specifies BI Connector for Atlas configuration on this cluster. BI Connector for Atlas is only available for M10+ clusters. See [BI Connector](#bi-connector) below for more details.
* `cluster_type` - Specifies the type of the cluster that you want to modify. You cannot convert a sharded cluster deployment to a replica set deployment.

    -> **WHEN SHOULD YOU USE CLUSTERTYPE?** 

        You set replicationSpecs.(Required)
        You are deploying Global Clusters. (Required)
        You are deploying non-Global replica sets and sharded clusters.

    Accepted values include:

        - `REPLICASET` Replica set
        - `SHARDED`	Sharded cluster
        - `GEOSHARDED` Global Cluster

* `disk_size_gb` - The size in gigabytes of the server’s root volume. You can add capacity by increasing this number, up to a maximum possible value of 4096 (i.e., 4 TB). This value must be a positive integer.

    The minimum disk size for dedicated clusters is 10GB for AWS and GCP, and 32GB for Azure. If you specify diskSizeGB with a lower disk size, Atlas defaults to the minimum disk size value.

* `encryption_at_rest_provider` - Set the Encryption at Rest parameter.

* `name` - (Required) Name of the cluster as it appears in Atlas. Once the cluster is created, its name cannot be changed.
* `mongo_db_major_version` - Version of the cluster to deploy. Atlas supports the following MongoDB versions for M10+ clusters:

    - 3.4
    - 3.6
    - 4.0

    You must set this value to 4.0 if `provider_instance_size_name` is either M2 or M5.
* `num_shards` - Selects whether the cluster is a replica set or a sharded cluster. If you use the replicationSpecs parameter, you must set num_shards.
* `provider_backup_enabled` - Flag indicating if the cluster uses Cloud Provider Snapshots for backups.

    If true, the cluster uses Cloud Provider Snapshots for backups. If providerBackupEnabled and backupEnabled are false, the cluster does not use Atlas backups.

    You cannot enable cloud provider snapshots if you have an existing cluster in the project with Continuous Backups enabled.
* `provider_instance_size_name` - (Required) Atlas provides different instance sizes, each with a default storage capacity and RAM size. The instance size you select is used for all the data-bearing servers in your cluster. See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/) `providerSettings.instanceSizeName` for valid values and default resources.
* `provider_name` - (Required) Cloud service provider on which the servers are provisioned.

    The possible values are:

    - AWS - Amazon AWS
    - GCP - Google Cloud Platform
    - AZURE - Microsoft Azure
    - TENANT - A multi-tenant deployment on one of the supported cloud service providers. Only valid when providerSettings.instanceSizeName is either M2 or M5.
* `backing_provider_name` - Cloud service provider on which the server for a multi-tenant cluster is provisioned.

    This setting is only valid when providerSetting.providerName is TENANT and providerSetting.instanceSizeName is M2 or M5.

    The possible values are:

    - AWS - Amazon AWS
    - GCP - Google Cloud Platform
    - AZURE - Microsoft Azure 

* `provider_disk_iops` - The maximum input/output operations per second (IOPS) the system can perform. The possible values depend on the selected providerSettings.instanceSizeName and diskSizeGB.
* `provider_disk_type_name` - Azure disk type of the server’s root volume. If omitted, Atlas uses the default disk type for the selected providerSettings.instanceSizeName.
* `provider_encrypt_ebs_volume` - If enabled, the Amazon EBS encryption feature encrypts the server’s root volume for both data at rest within the volume and for data moving between the volume and the instance.
* `provider_region_name` - Physical location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases.

    Do not specify this field when creating a multi-region cluster using the replicationSpec document or a Global Cluster with the replicationSpecs array.
* `provider_volume_type` - The type of the volume. The possible values are: `STANDARD` and `PROVISIONED`.
* `replication_factor` - Number of replica set members. Each member keeps a copy of your databases, providing high availability and data redundancy. The possible values are 3, 5, or 7. The default value is 3.

* `replication_specs` - Configuration for cluster regions.  See [Replication Spec](#replication-spec) below for more details.



### BI Connector

Specifies BI Connector for Atlas configuration.

* `enabled` - Specifies whether or not BI Connector for Atlas is enabled on the cluster.
    - Set to `true` to enable BI Connector for Atlas.
    - Set to `false` to disable BI Connector for Atlas.

* `read_preference` - Specifies the read preference to be used by BI Connector for Atlas on the cluster. Each BI Connector for Atlas read preference contains a distinct combination of [readPreference](https://docs.mongodb.com/manual/core/read-preference/) and [readPreferenceTags](https://docs.mongodb.com/manual/core/read-preference/#tag-sets) options. For details on BI Connector for Atlas read preferences, refer to the [BI Connector Read Preferences Table](https://docs.atlas.mongodb.com/tutorial/create-global-writes-cluster/#bic-read-preferences).

    - Set to "primary" to have BI Connector for Atlas read from the primary.

    - Set to "secondary" to have BI Connector for Atlas read from a secondary member. Default if there are no analytics nodes in the cluster.

    - Set to "analytics" to have BI Connector for Atlas read from an analytics node. Default if the cluster contains analytics nodes.

### Replication Spec

Configuration for cluster regions. 

* `id` - Unique identifer of the replication document for a zone in a Global Cluster.
* `num_shards` - (Required) Number of shards to deploy in the specified zone.
* `regions_config` - Physical location of the region. Each regionsConfig document describes the region’s priority in elections and the number and type of MongoDB nodes Atlas deploys to the region. You must order each regionsConfigs document by regionsConfig.priority, descending. See [Region Config](#region-config) below for more details.
* `zone_name` - Name for the zone in a Global Cluster.


### Region Config

Physical location of the region. 

* `region_name` - Name for the region specified.
* `electable_nodes` - Number of electable nodes for Atlas to deploy to the region. Electable nodes can become the primary and can facilitate local reads.
* `priority` -  Election priority of the region. For regions with only read-only nodes, set this value to 0.
* `read_only_nodes` - Number of read-only nodes for Atlas to deploy to the region. Read-only nodes can never become the primary, but can facilitate local-reads. Specify 0 if you do not want any read-only nodes in the region.
* `analytics_nodes` - The number of analytics nodes for Atlas to deploy to the region. Analytics nodes are useful for handling analytic data such as reporting queries from BI Connector for Atlas. Analytics nodes are read-only, and can never become the primary.

    If you do not specify this option, no analytics nodes are deployed to the region.



See detailed information for arguments and attributes: [MongoDB API Clusters](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/)