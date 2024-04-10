---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cluster"
sidebar_current: "docs-mongodbatlas-datasource-clusters"
description: |-
    Describe all Clusters in Project.
---

# Data Source: mongodbatlas_clusters

`mongodbatlas_cluster` describes all Clusters by the provided project_id. The data source requires your Project ID.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:**
<br> &#8226; Multi Region Cluster: The `mongodbatlas_cluster` data source doesn't return the `container_id` for each region utilized by the cluster. For retrieving the `container_id`, we recommend the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/advanced_clusters) data source instead.
<br> &#8226; Changes to cluster configurations can affect costs. Before making changes, please see [Billing](https://docs.atlas.mongodb.com/billing/).
<br> &#8226; If your Atlas project contains a custom role that uses actions introduced in a specific MongoDB version, you cannot create a cluster with a MongoDB version less than that version unless you delete the custom role.

## Example Usage

```terraform
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "cluster-test"
  cluster_type = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup                = true
  auto_scaling_disk_gb_enabled = true

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_instance_size_name = "M40"
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

* `name` - Name of the cluster as it appears in Atlas.
* `mongo_db_version` - Version of MongoDB the cluster runs, in `major-version`.`minor-version` format.
* `mongo_uri` - Base connection string for the cluster. Atlas only displays this field after the cluster is operational, not while it builds the cluster.
* `mongo_uri_updated` - Lists when the connection string was last updated. The connection string changes, for example, if you change a replica set to a sharded cluster.
* `mongo_uri_with_options` - Describes connection string for connecting to the Atlas cluster. Includes the replicaSet, ssl, and authSource query parameters in the connection string with values appropriate for the cluster.

    To review the connection string format, see the connection string format documentation. To add MongoDB users to a Atlas project, see Configure MongoDB Users.

    Atlas only displays this field after the cluster is operational, not while it builds the cluster.
* `paused` - Flag that indicates whether the cluster is paused or not.
* `pit_enabled` - Flag that indicates if the cluster uses Continuous Cloud Backup.
* `srv_address` - Connection string for connecting to the Atlas cluster. The +srv modifier forces the connection to use TLS/SSL. See the mongoURI for additional options.
* `state_name` - Indicates the current state of the cluster. The possible states are:
    - IDLE
    - CREATING
    - UPDATING
    - DELETING
    - DELETED
    - REPAIRING
* `auto_scaling_disk_gb_enabled` - Indicates whether disk auto-scaling is enabled.
* `auto_scaling_compute_enabled` - Specifies whether cluster tier auto-scaling is enabled. The default is false.
* `auto_scaling_compute_scale_down_enabled` -  * `auto_scaling_compute_scale_down_enabled` - Specifies whether cluster tier auto-down-scaling is enabled.
* `backup_enabled` - Legacy Option, Indicates whether Atlas continuous backups are enabled for the cluster.
* `bi_connector_config` - Indicates BI Connector for Atlas configuration on this cluster. BI Connector for Atlas is only available for M10+ clusters. See [BI Connector](#bi-connector) below for more details.
* `cluster_type` - Indicates the type of the cluster that you want to modify. You cannot convert a sharded cluster deployment to a replica set deployment.
* `connection_strings` - Set of connection strings that your applications use to connect to this cluster. More info in [Connection-strings](https://docs.mongodb.com/manual/reference/connection-string/). Use the parameters in this object to connect your applications to this cluster. To learn more about the formats of connection strings, see [Connection String Options](https://docs.atlas.mongodb.com/reference/faq/connection-changes/). NOTE: Atlas returns the contents of this object after the cluster is operational, not while it builds the cluster.
    - `connection_strings.standard` -   Public mongodb:// connection string for this cluster.
    - `connection_strings.standard_srv` - Public mongodb+srv:// connection string for this cluster. The mongodb+srv protocol tells the driver to look up the seed list of hosts in DNS. Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don’t need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn’t, use connectionStrings.standard.
    - `connection_strings.private` -   [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_srv` -  [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb+srv://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_endpoint.#.connection_string` - Private-endpoint-aware `mongodb://`connection string for this private endpoint.
    - `connection_strings.private_endpoint.#.srv_connection_string` - Private-endpoint-aware `mongodb+srv://` connection string for this private endpoint.
    - `connection_strings.private_endpoint.#.srv_shard_optimized_connection_string` - Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. 
    - `connection_strings.private_endpoint.#.type` - Type of MongoDB process that you connect to with the connection strings. Atlas returns `MONGOD` for replica sets, or `MONGOS` for sharded clusters.
    - `connection_strings.private_endpoint.#.endpoints` - Private endpoint through which you connect to Atlas when you use `connection_strings.private_endpoint[n].connection_string` or `connection_strings.private_endpoint[n].srv_connection_string`
    - `connection_strings.private_endpoint.#.endpoints.#.endpoint_id` - Unique identifier of the private endpoint.
    - `connection_strings.private_endpoint.#.endpoints.#.provider_name` - Cloud provider to which you deployed the private endpoint. Atlas returns `AWS` or `AZURE`.
    - `connection_strings.private_endpoint.#.endpoints.#.region` - Region to which you deployed the private endpoint.
* `disk_size_gb` - Indicates the size in gigabytes of the server’s root volume (AWS/GCP Only).
* `encryption_at_rest_provider` - Indicates whether Encryption at Rest is enabled or disabled.
* `tags` - Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#tags).
* `labels` - Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#labels). **DEPRECATED** Use `tags` instead.
* `mongo_db_major_version` - Indicates the version of the cluster to deploy.
* `num_shards` - Indicates whether the cluster is a replica set or a sharded cluster.
* `provider_backup_enabled` - Flag indicating if the cluster uses Cloud Backup Snapshots for backups. **DEPRECATED** Use `cloud_backup` instead.
* `cloud_backup` - Flag indicating if the cluster uses Cloud Backup Snapshots for backups.
* `termination_protection_enabled` - Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
* `provider_instance_size_name` - Atlas provides different instance sizes, each with a default storage capacity and RAM size.
* `provider_name` - Indicates the cloud service provider on which the servers are provisioned.
* `backing_provider_name` - Indicates Cloud service provider on which the server for a multi-tenant cluster is provisioned.
* `provider_disk_iops` - Indicates the maximum input/output operations per second (IOPS) the system can perform. The possible values depend on the selected providerSettings.instanceSizeName and diskSizeGB.
* `provider_disk_type_name` - Describes Azure disk type of the server’s root volume (Azure Only).
* `provider_encrypt_ebs_volume` - **(DEPRECATED)** Indicates whether the Amazon EBS encryption is enabled. This feature encrypts the server’s root volume for both data at rest within the volume and data moving between the volume and the instance. By default this attribute is always enabled, per deprecation process showing the real value at `provider_encrypt_ebs_volume_flag` computed attribute. 
* `provider_region_name` - Indicates Physical location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases. Requires the Atlas Region name, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
* `provider_volume_type` - Indicates the type of the volume. The possible values are: `STANDARD` and `PROVISIONED`.
-> **NOTE:** `STANDARD` is not available for NVME clusters.
* `provider_auto_scaling_compute_min_instance_size` - Minimum instance size to which your cluster can automatically scale.
* `provider_auto_scaling_compute_max_instance_size` - Maximum instance size to which your cluster can automatically scale.
* `replication_factor` - (Deprecated) Number of replica set members. Each member keeps a copy of your databases, providing high availability and data redundancy. The possible values are 3, 5, or 7. The default value is 3.

* `replication_specs` - Configuration for cluster regions.  See [Replication Spec](#replication-spec) below for more details.

* `container_id` - The Network Peering Container ID. 
  
-> **NOTE:** If you need to get an existing container ID see the [How-To Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/howto-guide.html).

* `version_release_system` - Release cadence that Atlas uses for this cluster.

* `advanced_configuration` - Get the advanced configuration options. See [Advanced Configuration](#advanced-configuration) below for more details.

### BI Connector

Indicates BI Connector for Atlas configuration.

* `enabled` - Indicates whether or not BI Connector for Atlas is enabled on the cluster.
* `read_preference` - Indicates the read preference to be used by BI Connector for Atlas on the cluster. Each BI Connector for Atlas read preference contains a distinct combination of [readPreference](https://docs.mongodb.com/manual/core/read-preference/) and [readPreferenceTags](https://docs.mongodb.com/manual/core/read-preference/#tag-sets) options. For details on BI Connector for Atlas read preferences, refer to the [BI Connector Read Preferences Table](https://docs.atlas.mongodb.com/tutorial/create-global-writes-cluster/#bic-read-preferences).

### Replication Spec

Configuration for cluster regions.

* `id` - Unique identifer of the replication document for a zone in a Global Cluster.
* `num_shards` - Number of shards to deploy in the specified zone.
* `regions_config` - Describes the physical location of the region. Each regionsConfig document describes the region’s priority in elections and the number and type of MongoDB nodes Atlas deploys to the region. You must order each regionsConfigs document by regionsConfig.priority, descending. See [Region Config](#region-config) below for more details.
* `zone_name` - Indicates the n ame for the zone in a Global Cluster.


### Region Config

Physical location of the region.

* `region_name` - Name for the region specified.
* `electable_nodes` - Number of electable nodes for Atlas to deploy to the region.
* `priority` -  Election priority of the region. For regions with only read-only nodes, set this value to 0.
* `read_only_nodes` - Number of read-only nodes for Atlas to deploy to the region. Read-only nodes can never become the primary, but can facilitate local-reads. Specify 0 if you do not want any read-only nodes in the region.
* `analytics_nodes` - Indicates the number of analytics nodes for Atlas to deploy to the region. Analytics nodes are useful for handling analytic data such as reporting queries from BI Connector for Atlas. Analytics nodes are read-only, and can never become the primary.

### Tags

Key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.

* `key` - Constant that defines the set of the tag.
* `value` - Variable that belongs to the set of the tag.

To learn more, see [Resource Tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas).

### Labels

**WARNING:** This property is deprecated and will be removed in the future, use the `tags` attribute instead.

Key-value pairs that categorize the cluster. Each key and value has a maximum length of 255 characters.  You cannot set the key `Infrastructure Tool`, it is used for internal purposes to track aggregate usage.

* `key` - The key that you want to write.
* `value` - The value that you want to write.

-> **NOTE:** MongoDB Atlas doesn't display your labels.

### Plugin
Contains a key-value pair that tags that the cluster was created by a Terraform Provider and notes the version.

* `name` - The name of the current plugin
* `version` - The current version of the plugin.

### Cloud Provider Snapshot Backup Policy
**WARNING:** This property is deprecated, use `mongodbatlas_cloud_backup_schedule` datasource instead.

* `snapshot_backup_policy` - current snapshot schedule and retention settings for the cluster.

* `snapshot_backup_policy.#.cluster_id` - Unique identifier of the Atlas cluster.
* `snapshot_backup_policy.#.cluster_name` - Name of the Atlas cluster that contains the snapshot backup policy.
* `snapshot_backup_policy.#.next_snapshot` - UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
* `snapshot_backup_policy.#.reference_hour_of_day` - UTC Hour of day between 0 and 23 representing which hour of the day that Atlas takes a snapshot.
* `snapshot_backup_policy.#.reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the referenceHourOfDay that Atlas takes the snapshot.
* `snapshot_backup_policy.#.restore_window_days` - Specifies a restore window in days for the cloud provider backup to maintain.
* `snapshot_backup_policy.#.update_snapshots` - Specifies it's true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.

### Policies
**WARNING:** This property is deprecated, use `mongodbatlas_cloud_backup_schedule` datasource instead.

* `snapshot_backup_policy.#.policies` - A list of policy definitions for the cluster.
* `snapshot_backup_policy.#.policies.#.id` - Unique identifier of the backup policy.

#### Policy Item
**WARNING:** This property is deprecated, use `mongodbatlas_cloud_backup_schedule` datasource instead.

* `snapshot_backup_policy.#.policies.#.policy_item` - A list of specifications for a policy.
* `snapshot_backup_policy.#.policies.#.policy_item.#.id` - Unique identifier for this policy item.
* `snapshot_backup_policy.#.policies.#.policy_item.#.frequency_interval` - The frequency interval for a set of snapshots.
* `snapshot_backup_policy.#.policies.#.policy_item.#.frequency_type` - A type of frequency (hourly, daily, weekly, monthly).
* `snapshot_backup_policy.#.policies.#.policy_item.#.retention_unit` - The unit of time in which snapshot retention is measured (days, weeks, months).
* `snapshot_backup_policy.#.policies.#.policy_item.#.retention_value` - The number of days, weeks, or months the snapshot is retained.

#### Advanced Configuration

* `default_read_concern` - [Default level of acknowledgment requested from MongoDB for read operations](https://docs.mongodb.com/manual/reference/read-concern/) set for this cluster. MongoDB 4.4 clusters default to [available](https://docs.mongodb.com/manual/reference/read-concern-available/).
* `default_write_concern` -  [Default level of acknowledgment requested from MongoDB for write operations](https://docs.mongodb.com/manual/reference/write-concern/) set for this cluster. MongoDB 4.4 clusters default to [1](https://docs.mongodb.com/manual/reference/write-concern/).
* `fail_index_key_too_long` - When true, documents can only be updated or inserted if, for all indexed fields on the target collection, the corresponding index entries do not exceed 1024 bytes. When false, mongod writes documents that exceed the limit but does not index them.
* `javascript_enabled` - When true, the cluster allows execution of operations that perform server-side executions of JavaScript. When false, the cluster disables execution of those operations.
* `minimum_enabled_tls_protocol` - Sets the minimum Transport Layer Security (TLS) version the cluster accepts for incoming connections.Valid values are:

  - TLS1_0
  - TLS1_1
  - TLS1_2

* `no_table_scan` - When true, the cluster disables the execution of any query that requires a collection scan to return results. When false, the cluster allows the execution of those operations.
* `oplog_size_mb` - The custom oplog size of the cluster. Without a value that indicates that the cluster uses the default oplog size calculated by Atlas.
* `oplog_min_retention_hours` - Minimum retention window for cluster's oplog expressed in hours. A value of null indicates that the cluster uses the default minimum oplog window that MongoDB Cloud calculates.
* `sample_size_bi_connector` - Number of documents per database to sample when gathering schema information. Defaults to 100. Available only for Atlas deployments in which BI Connector for Atlas is enabled.
* `sample_refresh_interval_bi_connector` - Interval in seconds at which the mongosqld process re-samples data to create its relational schema. The default value is 300. The specified value must be a positive integer. Available only for Atlas deployments in which BI Connector for Atlas is enabled.


See detailed information for arguments and attributes: [MongoDB API Clusters](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/)