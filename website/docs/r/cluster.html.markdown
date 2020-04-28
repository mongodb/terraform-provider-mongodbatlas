---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cluster"
sidebar_current: "docs-mongodbatlas-resource-cluster"
description: |-
    Provides a Cluster resource.
---

# mongodbatlas_cluster

`mongodbatlas_cluster` provides a Cluster resource. The resource lets you create, edit and delete clusters. The resource requires your Project ID.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:**
<br> &#8226; Free tier cluster creation (M0) is not supported via API or by this Provider.
<br> &#8226; Shared tier clusters (M2, M5) cannot be upgraded to higher tiers via API or by this Provider.
<br> &#8226; Changes to cluster configurations can affect costs. Before making changes, please see [Billing](https://docs.atlas.mongodb.com/billing/).
<br> &#8226; If your Atlas project contains a custom role that uses actions introduced in a specific MongoDB version, you cannot create a cluster with a MongoDB version less than that version unless you delete the custom role.

## Example Usage

### Example AWS cluster

```hcl
resource "mongodbatlas_cluster" "cluster-test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "cluster-test"
  num_shards   = 1

  replication_factor           = 3
  provider_backup_enabled      = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "4.0"

  //Provider Settings "block"
  provider_name               = "AWS"
  disk_size_gb                = 100
  provider_disk_iops          = 300
  provider_volume_type        = "STANDARD"
  provider_encrypt_ebs_volume = true
  provider_instance_size_name = "M40"
  provider_region_name        = "US_EAST_1"
}
```

### Example Azure cluster.

```hcl
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "test"
  num_shards   = 1

  replication_factor           = 3
  provider_backup_enabled      = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "4.0"

  //Provider Settings "block"
  provider_name               = "AZURE"
  provider_disk_type_name     = "P6"
  provider_instance_size_name = "M30"
  provider_region_name        = "US_EAST_2"
}
```

### Example GCP cluster

```hcl
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "test"
  num_shards   = 1

  replication_factor           = 3
  provider_backup_enabled      = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "4.0"

  //Provider Settings "block"
  provider_name               = "GCP"
  disk_size_gb                = 40
  provider_instance_size_name = "M30"
  provider_region_name        = "US_EAST_4"
}
```

### Example Multi Region cluster

```hcl
resource "mongodbatlas_cluster" "cluster-test" {
  project_id               = "<YOUR-PROJECT-ID>"
  name                     = "cluster-test-multi-region"
  disk_size_gb             = 100
  num_shards               = 1
  provider_backup_enabled  = true
  cluster_type             = "REPLICASET"

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_disk_iops          = 300
  provider_volume_type        = "STANDARD"
  provider_instance_size_name = "M10"

  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
    regions_config {
      region_name     = "US_EAST_2"
      electable_nodes = 2
      priority        = 6
      read_only_nodes = 0
    }
    regions_config {
      region_name     = "US_WEST_1"
      electable_nodes = 2
      priority        = 5
      read_only_nodes = 2
    }
  }
}
```

### Example Global cluster

```hcl
resource "mongodbatlas_cluster" "cluster-test" {
  project_id              = "<YOUR-PROJECT-ID>"
  name                    = "cluster-test-global"
  disk_size_gb            = 80
  num_shards              = 1
  provider_backup_enabled = true
  cluster_type            = "GEOSHARDED"

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_disk_iops          = 240
  provider_volume_type        = "STANDARD"
  provider_instance_size_name = "M30"

  replication_specs {
    zone_name  = "Zone 1"
    num_shards = 2
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }

  replication_specs {
    zone_name  = "Zone 2"
    num_shards = 2
    regions_config {
      region_name     = "EU_CENTRAL_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
}
```
### Example AWS Shared Tier cluster
```hcl
resource "mongodbatlas_cluster" "cluster-test" {
  project_id              = "<YOUR-PROJECT-ID>"
  name                    = "cluster-test-global"
  //M2 must be 2, M5 must be 5
  disk_size_gb            = "2"

  //Provider Settings "block"
  provider_name = "TENANT"
  backing_provider_name = "AWS"
  provider_region_name = "US_EAST_1"
  provider_instance_size_name = "M2"

  //These must be the following values
  mongo_db_major_version = "4.2"
  auto_scaling_disk_gb_enabled = "false"
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `provider_name` - (Required) Cloud service provider on which the servers are provisioned.

    The possible values are:

    - `AWS` - Amazon AWS
    - `GCP` - Google Cloud Platform
    - `AZURE` - Microsoft Azure
    - `TENANT` - A multi-tenant deployment on one of the supported cloud service providers. Only valid when providerSettings.instanceSizeName is either M2 or M5.
* `name` - (Required) Name of the cluster as it appears in Atlas. Once the cluster is created, its name cannot be changed.
* `provider_instance_size_name` - (Required) Atlas provides different instance sizes, each with a default storage capacity and RAM size. The instance size you select is used for all the data-bearing servers in your cluster. See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/) `providerSettings.instanceSizeName` for valid values and default resources.

Note free tier (M0) creation is not supported by the Atlas API and hence not supported by this provider.

* `auto_scaling_disk_gb_enabled` - (Optional) Specifies whether disk auto-scaling is enabled. The default is true.
    - Set to `true` to enable disk auto-scaling.
    - Set to `false` to disable disk auto-scaling.

* `backup_enabled` - (Optional) Legacy Option - Set to true to enable Atlas continuous backups for the cluster.

    Important - Clusters running MongoDB 4.2 and any net new Atlas clusters of any type do not support this parameter. These clusters must use `provider_backup_enabled` to enable Cloud Provider Snapshots.  If you create a new Atlas cluster and set `backup_enabled` to true, the Provider will respond with an error.  This change doesn’t affect existing clusters that use continuous backups.

    Set to false to disable continuous backups for the cluster. Atlas deletes any stored snapshots.

    The default value is false.  M10 and above only.

* `bi_connector` - (Optional) Specifies BI Connector for Atlas configuration on this cluster. BI Connector for Atlas is only available for M10+ clusters. See [BI Connector](#bi-connector) below for more details.
* `cluster_type` - (Optional) Specifies the type of the cluster that you want to modify. You cannot convert a sharded cluster deployment to a replica set deployment.

    -> **WHEN SHOULD YOU USE CLUSTERTYPE?**
      When you set replication_specs, when you are deploying Global Clusters or when you are deploying non-Global replica sets and sharded clusters.

    Accepted values include:
      - `REPLICASET` Replica set
      - `SHARDED`	Sharded cluster
      - `GEOSHARDED` Global Cluster

* `disk_size_gb` - (Optional - GCP/AWS Only) The size in gigabytes of the server’s root volume. You can add capacity by increasing this number, up to a maximum possible value of 4096 (i.e., 4 TB). This value must be a positive integer.

    The minimum disk size for dedicated clusters is 10GB for AWS and GCP. If you specify diskSizeGB with a lower disk size, Atlas defaults to the minimum disk size value.

* `encryption_at_rest_provider` - (Optional) Set the Encryption at Rest parameter.  Possible values are AWS, GCP, AZURE or NONE.  Requires M10 or greater and for backup_enabled to be false or omitted.
* `labels` - (Optional) Array containing key-value pairs that tag and categorize the cluster. Each key and value has a maximum length of 255 characters. You cannot set the key `Infrastructure Tool`, it is used for internal purposes to track aggregate usage.
* `mongo_db_major_version` - (Optional) Version of the cluster to deploy. Atlas supports the following MongoDB versions for M10+ clusters: `3.6`, `4.0`, or `4.2`. You must set this value to `4.2` if `provider_instance_size_name` is either M2 or M5.
* `num_shards` - (Optional) Selects whether the cluster is a replica set or a sharded cluster. If you use the replicationSpecs parameter, you must set num_shards.
* `pit_enabled` - (Optional) - Flag that indicates if the cluster uses Point-in-Time backups. If set to true, provider_backup_enabled must also be set to true.
* `provider_backup_enabled` - (Optional) Flag indicating if the cluster uses Cloud Provider Snapshots for backups.

    If true, the cluster uses Cloud Provider Snapshots for backups. If provider_backup_enabled and backup_enabled are false, the cluster does not use Atlas backups.

    You cannot enable cloud provider snapshots if you have an existing cluster in the project with Continuous Backups enabled.

    ~> **IMPORTANT:** If setting to true for an existing cluster or imported cluster be sure to run terraform refresh after applying to enable modification of the Cloud Provider Snapshot Backup Policy going forward.

* `backing_provider_name` - (Optional) Cloud service provider on which the server for a multi-tenant cluster is provisioned.

    This setting is only valid when providerSetting.providerName is TENANT and providerSetting.instanceSizeName is M2 or M5.

    The possible values are:

    - AWS - Amazon AWS
    - GCP - Google Cloud Platform
    - AZURE - Microsoft Azure

* `provider_disk_iops` - (Optional) The maximum input/output operations per second (IOPS) the system can perform. The possible values depend on the selected providerSettings.instanceSizeName and diskSizeGB.
* `provider_disk_type_name` - (Optional - Azure Only) Azure disk type of the server’s root volume. If omitted, Atlas uses the default disk type for the selected providerSettings.instanceSizeName.  Example disk types and associated storage sizes: P4 - 32GB, P6 - 64GB, P10 - 128GB, P20 - 512GB, P30 - 1024GB, P40 - 2048GB, P50 - 4095GB.  More information and the most update to date disk types/storage sizes can be located at https://docs.atlas.mongodb.com/reference/api/clusters-create-one/.
* `provider_encrypt_ebs_volume` - (Optional) If enabled, the Amazon EBS encryption feature encrypts the server’s root volume for both data at rest within the volume and for data moving between the volume and the instance.
* `provider_region_name` - (Optional) Physical location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases.  Requires the Atlas Region name, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
    Do not specify this field when creating a multi-region cluster using the replicationSpec document or a Global Cluster with the replicationSpecs array.
* `provider_volume_type` - (AWS - Optional) The type of the volume. The possible values are: `STANDARD` and `PROVISIONED`.  `PROVISIONED` required if setting IOPS higher than the default instance IOPS.
* `replication_factor` - (Optional) Number of replica set members. Each member keeps a copy of your databases, providing high availability and data redundancy. The possible values are 3, 5, or 7. The default value is 3.

* `replication_specs` - (Optional) Configuration for cluster regions.  See [Replication Spec](#replication-spec) below for more details.



### BI Connector

Specifies BI Connector for Atlas configuration.

* `enabled` - (Optional) Specifies whether or not BI Connector for Atlas is enabled on the cluster.
    - Set to `true` to enable BI Connector for Atlas.
    - Set to `false` to disable BI Connector for Atlas.

* `read_preference` - (Optional) Specifies the read preference to be used by BI Connector for Atlas on the cluster. Each BI Connector for Atlas read preference contains a distinct combination of [readPreference](https://docs.mongodb.com/manual/core/read-preference/) and [readPreferenceTags](https://docs.mongodb.com/manual/core/read-preference/#tag-sets) options. For details on BI Connector for Atlas read preferences, refer to the [BI Connector Read Preferences Table](https://docs.atlas.mongodb.com/tutorial/create-global-writes-cluster/#bic-read-preferences).

    - Set to "primary" to have BI Connector for Atlas read from the primary.

    - Set to "secondary" to have BI Connector for Atlas read from a secondary member. Default if there are no analytics nodes in the cluster.

    - Set to "analytics" to have BI Connector for Atlas read from an analytics node. Default if the cluster contains analytics nodes.

### Replication Spec

Configuration for cluster regions.

* `num_shards` - (Required) Number of shards to deploy in the specified zone.
* `id` - (Optional) Unique identifer of the replication document for a zone in a Global Cluster.
* `regions_config` - (Optional) Physical location of the region. Each regionsConfig document describes the region’s priority in elections and the number and type of MongoDB nodes Atlas deploys to the region. You must order each regionsConfigs document by regionsConfig.priority, descending. See [Region Config](#region-config) below for more details.
* `zone_name` - (Optional) Name for the zone in a Global Cluster.


### Region Config

Physical location of the region.

* `region_name` - (Optional) Name for the region specified.
* `electable_nodes` - (Optional) Number of electable nodes for Atlas to deploy to the region. Electable nodes can become the primary and can facilitate local reads.
* `priority` - (Optional)  Election priority of the region. For regions with only read-only nodes, set this value to 0.
* `read_only_nodes` - (Optional) Number of read-only nodes for Atlas to deploy to the region. Read-only nodes can never become the primary, but can facilitate local-reads. Specify 0 if you do not want any read-only nodes in the region.
* `analytics_nodes` - (Optional) The number of analytics nodes for Atlas to deploy to the region. Analytics nodes are useful for handling analytic data such as reporting queries from BI Connector for Atlas. Analytics nodes are read-only, and can never become the primary.

    If you do not specify this option, no analytics nodes are deployed to the region.

### Advanced Configuration Options

-> **NOTE:** Prior to setting these options please ensure you read https://docs.atlas.mongodb.com/cluster-config/additional-options/.

Include **desired options** within advanced_configuration:

```hcl
// Nest options within advanced_configuration
 advanced_configuration = {
   javascript_enabled                   = false
   minimum_enabled_tls_protocol         = "TLS1_2"
 }
```

* `fail_index_key_too_long` - (Optional) When true, documents can only be updated or inserted if, for all indexed fields on the target collection, the corresponding index entries do not exceed 1024 bytes. When false, mongod writes documents that exceed the limit but does not index them.
* `javascript_enabled` - (Optional) When true, the cluster allows execution of operations that perform server-side executions of JavaScript. When false, the cluster disables execution of those operations.
* `minimum_enabled_tls_protocol` - (Optional) Sets the minimum Transport Layer Security (TLS) version the cluster accepts for incoming connections.Valid values are:

  - TLS1_0
  - TLS1_1
  - TLS1_2

* `no_table_scan` - (Optional) When true, the cluster disables the execution of any query that requires a collection scan to return results. When false, the cluster allows the execution of those operations.
* `oplog_size_mb` - (Optional) The custom oplog size of the cluster. Without a value that indicates that the cluster uses the default oplog size calculated by Atlas.
* `sample_size_bi_connector` - (Optional) Number of documents per database to sample when gathering schema information. Defaults to 100. Available only for Atlas deployments in which BI Connector for Atlas is enabled.
* `sample_refresh_interval_bi_connector` - (Optional) Interval in seconds at which the mongosqld process re-samples data to create its relational schema. The default value is 300. The specified value must be a positive integer. Available only for Atlas deployments in which BI Connector for Atlas is enabled.

### Labels
Contains key-value pairs that tag and categorize the cluster. Each key and value has a maximum length of 255 characters.

* `key` - The key that you want to write.
* `value` - The value that you want to write.



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - The cluster ID.
*  `mongo_db_version` - Version of MongoDB the cluster runs, in `major-version`.`minor-version` format.
* `id` -	The Terraform's unique identifier used internally for state management.
* `mongo_uri` - Base connection string for the cluster. Atlas only displays this field after the cluster is operational, not while it builds the cluster.
* `mongo_uri_updated` - Lists when the connection string was last updated. The connection string changes, for example, if you change a replica set to a sharded cluster.
* `mongo_uri_with_options` - connection string for connecting to the Atlas cluster. Includes the replicaSet, ssl, and authSource query parameters in the connection string with values appropriate for the cluster.
* `connection_strings` - Set of connection strings that your applications use to connect to this cluster. More info in [Connection-strings](https://docs.mongodb.com/manual/reference/connection-string/). Use the parameters in this object to connect your applications to this cluster. To learn more about the formats of connection strings, see [Connection String Options](https://docs.atlas.mongodb.com/reference/faq/connection-changes/). NOTE: Atlas returns the contents of this object after the cluster is operational, not while it builds the cluster.
    - `connection_strings.standard` -   Public mongodb:// connection string for this cluster.
    - `connection_strings.standard_srv` - Public mongodb+srv:// connection string for this cluster. The mongodb+srv protocol tells the driver to look up the seed list of hosts in DNS. Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don’t need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn’t, use connectionStrings.standard.
    - `connection_strings.aws_private_link` -  [Private-endpoint-aware](https://docs.atlas.mongodb.com/security-private-endpoint/#private-endpoint-connection-strings) mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a AWS PrivateLink connection to this cluster.
    - `connection_strings.aws_private_link_srv` - [Private-endpoint-aware](https://docs.atlas.mongodb.com/security-private-endpoint/#private-endpoint-connection-strings) mongodb+srv://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a AWS PrivateLink connection to this cluster. Use this URI format if your driver supports it. If it doesn’t, use connectionStrings.awsPrivateLink.
    - `connection_strings.private` -   [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_srv` -  [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb+srv://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
* `container_id` - The Network Peering Container ID.

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

### Cloud Provider Snapshot Backup Policy
Cloud Provider Snapshot Backup Policy will be added if provider_backup_enabled is enabled, if not, returned values will be empty.

* `snapshot_backup_policy` - current snapshot schedule and retention settings for the cluster.

* `snapshot_backup_policy.#.cluster_id` - Unique identifier of the Atlas cluster.
* `snapshot_backup_policy.#.cluster_name` - Name of the Atlas cluster that contains the snapshot backup policy.
* `snapshot_backup_policy.#.next_snapshot` - UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
* `snapshot_backup_policy.#.reference_hour_of_day` - UTC Hour of day between 0 and 23 representing which hour of the day that Atlas takes a snapshot.
* `snapshot_backup_policy.#.reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the referenceHourOfDay that Atlas takes the snapshot.
* `snapshot_backup_policy.#.restore_window_days` - Specifies a restore window in days for the cloud provider backup to maintain.
* `snapshot_backup_policy.#.update_snapshots` - Specifies it's true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.

### Policies
* `snapshot_backup_policy.#.policies` - A list of policy definitions for the cluster.
* `snapshot_backup_policy.#.policies.#.id` - Unique identifier of the backup policy.

#### Policy Item
* `snapshot_backup_policy.#.policies.#.policy_item` - A list of specifications for a policy.
* `snapshot_backup_policy.#.policies.#.policy_item.#.id` - Unique identifier for this policy item.
* `snapshot_backup_policy.#.policies.#.policy_item.#.frequency_interval` - The frequency interval for a set of snapshots.
* `snapshot_backup_policy.#.policies.#.policy_item.#.frequency_type` - A type of frequency (hourly, daily, weekly, monthly).
* `snapshot_backup_policy.#.policies.#.policy_item.#.retention_unit` - The unit of time in which snapshot retention is measured (days, weeks, months).
* `snapshot_backup_policy.#.policies.#.policy_item.#.retention_value` - The number of days, weeks, or months the snapshot is retained.

## Import

Clusters can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cluster.my_cluster 1112222b3bf99403840e8934-Cluster0
```

See detailed information for arguments and attributes: [MongoDB API Clusters](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/)
