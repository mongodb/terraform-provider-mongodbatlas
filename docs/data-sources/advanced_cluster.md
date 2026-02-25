---
subcategory: "Clusters"
---

# Data Source: mongodbatlas_advanced_cluster

`mongodbatlas_advanced_cluster` describes an Advanced Cluster. The data source requires your Project ID.


-> **NOTE:** Groups and projects are synonymous terms. You might find group_id in the official documentation.

~> **IMPORTANT:**
<br> &#8226; Changes to cluster configurations can affect costs. Before making changes, please see [Billing](https://docs.atlas.mongodb.com/billing/).
<br> &#8226; If your Atlas project contains a custom role that uses actions introduced in a specific MongoDB version, you cannot create a cluster with a MongoDB version less than that version unless you delete the custom role.

-> **NOTE:** To delete an Atlas cluster that has an associated `mongodbatlas_cloud_backup_schedule` resource and an enabled Backup Compliance Policy, first instruct Terraform to remove the `mongodbatlas_cloud_backup_schedule` resource from the state and then use Terraform to delete the cluster. To learn more, see [Delete a Cluster with a Backup Compliance Policy](../guides/delete-cluster-with-backup-compliance-policy.md).

-> **NOTE:** This data source also includes Flex clusters.

## Example Usage

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "cluster-test"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M0"
          }
          provider_name         = "TENANT"
          backing_provider_name = "AWS"
          region_name           = "US_EAST_1"
          priority              = 7
        }
      ]
    }
  ]
}

data "mongodbatlas_advanced_cluster" "this" {
	project_id = mongodbatlas_advanced_cluster.this.project_id
	name 	   = mongodbatlas_advanced_cluster.this.name
}
```

## Example using effective fields with auto-scaling

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id            = "<YOUR-PROJECT-ID>"
  name                  = "auto-scale-cluster"
  cluster_type          = "REPLICASET"
  use_effective_fields  = true

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10" # Initial size value that won't change in Terraform state, actual size in Atlas may differ due to auto-scaling
            node_count    = 3
          }
          auto_scaling = {
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M10"
            compute_max_instance_size  = "M30"
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]
}

# Read effective values after Atlas auto-scales the cluster
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}

output "configured_instance_size" {
  value = data.mongodbatlas_advanced_cluster.this.replication_specs[0].region_configs[0].electable_specs.instance_size
}

output "actual_instance_size" {
  value = data.mongodbatlas_advanced_cluster.this.replication_specs[0].region_configs[0].effective_electable_specs.instance_size
}
```

## Example using latest sharding configurations with independent shard scaling in the cluster

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = "<YOUR-PROJECT-ID>"
  name           = "cluster-test"
  backup_enabled = false
  cluster_type   = "SHARDED"

  replication_specs = [
    {    # Sharded cluster with 2 asymmetric shards (M30 and M40)
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            disk_iops     = 3000
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    },
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M40"
            disk_iops     = 3000
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    }
  ]
}

data "mongodbatlas_advanced_cluster" "this" {
  project_id                     = mongodbatlas_advanced_cluster.this.project_id
  name                           = mongodbatlas_advanced_cluster.this.name
}
```

## Example using Flex cluster

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "flex-cluster"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          provider_name = "FLEX"
          backing_provider_name = "AWS"
          region_name = "US_EAST_1"
          priority = 7
        }
      ]
    }
  ]
}

data "mongodbatlas_advanced_cluster" "this" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Human-readable label that identifies this cluster.
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.

**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.

### Optional

- `use_effective_fields` (Boolean) Controls how hardware specification fields are returned in the response. When set to true, the non-effective specs (`electable_specs`, `read_only_specs`, `analytics_specs`) fields return the hardware specifications that the client provided. When set to false (default), the non-effective specs fields show the **current** hardware specifications. Cluster auto-scaling is the primary cause for differences between initial and current hardware specifications.

### Read-Only

- `advanced_configuration` (Attributes) Additional settings for an Atlas cluster. (see [below for nested schema](#nestedatt--advanced_configuration))
- `backup_enabled` (Boolean) Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.
- `bi_connector_config` (Attributes) Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster. (see [below for nested schema](#nestedatt--bi_connector_config))
- `cluster_id` (String) Unique 24-hexadecimal digit string that identifies the cluster.
- `cluster_type` (String) Configuration of nodes that comprise the cluster.
- `config_server_management_mode` (String) Config Server Management Mode for creating or updating a sharded cluster.

When configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.

When configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.
- `config_server_type` (String) Describes a sharded cluster's config server type.
- `connection_strings` (Attributes) Collection of Uniform Resource Locators that point to the MongoDB database. (see [below for nested schema](#nestedatt--connection_strings))
- `create_date` (String) Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.
- `encryption_at_rest_provider` (String) Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster **replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize** setting must be `M10` or higher and `"backupEnabled" : false` or omitted entirely.
- `global_cluster_self_managed_sharding` (Boolean) Set this field to configure the Sharding Management Mode when creating a new Global Cluster.

When set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.

When set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.

This setting cannot be changed once the cluster is deployed.
- `labels` (Map of String) Map of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.

Cluster labels are deprecated and will be removed in a future release. We strongly recommend that you use [resource tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas) instead.
- `mongo_db_major_version` (String) MongoDB major version of the cluster.

On creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLTSVersions).

 On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version.
- `mongo_db_version` (String) Version of MongoDB that the cluster runs.
- `paused` (Boolean) Flag that indicates whether the cluster is paused.
- `pinned_fcv` (Attributes) Pins the Feature Compatibility Version (FCV) to the current MongoDB version with a provided expiration date. To unpin the FCV the `pinned_fcv` attribute must be removed. This operation can take several minutes as the request processes through the MongoDB data plane. Once FCV is unpinned it will not be possible to downgrade the `mongo_db_major_version`. It is advised that updates to `pinned_fcv` are done isolated from other cluster changes. If a plan contains multiple changes, the FCV change will be applied first. If FCV is unpinned past the expiration date the `pinned_fcv` attribute must be removed. The following [knowledge hub article](https://kb.corp.mongodb.com/article/000021785/) and [FCV documentation](https://www.mongodb.com/docs/atlas/tutorial/major-version-change/#manage-feature-compatibility--fcv--during-upgrades) can be referenced for more details. (see [below for nested schema](#nestedatt--pinned_fcv))
- `pit_enabled` (Boolean) Flag that indicates whether the cluster uses continuous cloud backups.
- `redact_client_log_data` (Boolean) Enable or disable log redaction.

This setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging. This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.

Use ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.

*Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.
- `replica_set_scaling_strategy` (String) Set this field to configure the replica set scaling mode for your cluster.

By default, Atlas scales under WORKLOAD_TYPE. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.

When configured as SEQUENTIAL, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads.

When configured as NODE_TYPE, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.
- `replication_specs` (Attributes List) List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations. (see [below for nested schema](#nestedatt--replication_specs))
- `root_cert_type` (String) Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group.
- `state_name` (String) Human-readable label that indicates the current operating condition of this cluster.
- `tags` (Map of String) Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.
- `termination_protection_enabled` (Boolean) Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.
- `use_aws_time_based_snapshot_copy_for_fast_initial_sync` (Boolean) Flag that indicates whether time-based snapshot copies will be used instead of slower standard snapshot copies during fast Atlas cross-region initial syncs. This flag is only relevant for clusters containing AWS nodes.
- `version_release_system` (String) Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify **mongoDBMajorVersion**.

<a id="nestedatt--advanced_configuration"></a>
### Nested Schema for `advanced_configuration`

Read-Only:

- `change_stream_options_pre_and_post_images_expire_after_seconds` (Number) The minimum pre- and post-image retention time in seconds.
- `custom_openssl_cipher_config_tls12` (Set of String) The custom OpenSSL cipher suite list for TLS 1.2. This field is only valid when `tls_cipher_config_mode` is set to `CUSTOM`.
- `custom_openssl_cipher_config_tls13` (Set of String) The custom OpenSSL cipher suite list for TLS 1.3. This field is only valid when `tls_cipher_config_mode` is set to `CUSTOM`.
- `default_max_time_ms` (Number) Default time limit in milliseconds for individual read operations to complete. This parameter is supported only for MongoDB version 8.0 and above.
- `default_write_concern` (String) Default level of acknowledgment requested from MongoDB for write operations when none is specified by the driver.
- `javascript_enabled` (Boolean) Flag that indicates whether the cluster allows execution of operations that perform server-side executions of JavaScript. When using 8.0+, we recommend disabling server-side JavaScript and using operators of aggregation pipeline as more performant alternative.
- `minimum_enabled_tls_protocol` (String) Minimum Transport Layer Security (TLS) version that the cluster accepts for incoming connections. Clusters using TLS 1.0 or 1.1 should consider setting TLS 1.2 as the minimum TLS protocol version.
- `no_table_scan` (Boolean) Flag that indicates whether the cluster disables executing any query that requires a collection scan to return results.
- `oplog_min_retention_hours` (Number) Minimum retention window for cluster's oplog expressed in hours. A value of null indicates that the cluster uses the default minimum oplog window that MongoDB Cloud calculates.
- `oplog_size_mb` (Number) Storage limit of cluster's oplog expressed in megabytes. A value of null indicates that the cluster uses the default oplog size that MongoDB Cloud calculates.
- `sample_refresh_interval_bi_connector` (Number) Interval in seconds at which the mongosqld process re-samples data to create its relational schema.
- `sample_size_bi_connector` (Number) Number of documents per database to sample when gathering schema information.
- `tls_cipher_config_mode` (String) The TLS cipher suite configuration mode. Valid values include `CUSTOM` or `DEFAULT`. The `DEFAULT` mode uses the default cipher suites. The `CUSTOM` mode allows you to specify custom cipher suites for both TLS 1.2 and TLS 1.3. To unset, this should be set back to `DEFAULT`.
- `transaction_lifetime_limit_seconds` (Number) Lifetime, in seconds, of multi-document transactions. Atlas considers the transactions that exceed this limit as expired and so aborts them through a periodic cleanup process.


<a id="nestedatt--bi_connector_config"></a>
### Nested Schema for `bi_connector_config`

Read-Only:

- `enabled` (Boolean) Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.
- `read_preference` (String) Data source node designated for the MongoDB Connector for Business Intelligence on MongoDB Cloud. The MongoDB Connector for Business Intelligence on MongoDB Cloud reads data from the primary, secondary, or analytics node based on your read preferences. Defaults to `ANALYTICS` node, or `SECONDARY` if there are no `ANALYTICS` nodes.


<a id="nestedatt--connection_strings"></a>
### Nested Schema for `connection_strings`

Read-Only:

- `private` (String) Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter once someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private. For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.
- `private_endpoint` (Attributes List) List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint. This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes. (see [below for nested schema](#nestedatt--connection_strings--private_endpoint))
- `private_srv` (String) Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter when someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this parameter returns only if you [enable custom DNS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/).
- `standard` (String) Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol.
- `standard_srv` (String) Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb+srv://` protocol.

<a id="nestedatt--connection_strings--private_endpoint"></a>
### Nested Schema for `connection_strings.private_endpoint`

Read-Only:

- `connection_string` (String) Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to MongoDB Cloud through a private endpoint.
- `endpoints` (Attributes List) List that contains the private endpoints through which you connect to MongoDB Cloud when you use **connectionStrings.privateEndpoint[n].connectionString** or **connectionStrings.privateEndpoint[n].srvConnectionString**. (see [below for nested schema](#nestedatt--connection_strings--private_endpoint--endpoints))
- `srv_connection_string` (String) Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application supports it. If it doesn't, use connectionStrings.privateEndpoint[n].connectionString.
- `srv_shard_optimized_connection_string` (String) Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for connectionStrings.privateEndpoint[n].srvConnectionString.
- `type` (String) MongoDB process type to which your application connects. Use `MONGOD` for replica sets and `MONGOS` for sharded clusters.

<a id="nestedatt--connection_strings--private_endpoint--endpoints"></a>
### Nested Schema for `connection_strings.private_endpoint.endpoints`

Read-Only:

- `endpoint_id` (String) Unique string that the cloud provider uses to identify the private endpoint.
- `provider_name` (String) Cloud provider in which MongoDB Cloud deploys the private endpoint.
- `region` (String) Region where the private endpoint is deployed.




<a id="nestedatt--pinned_fcv"></a>
### Nested Schema for `pinned_fcv`

Read-Only:

- `expiration_date` (String) Expiration date of the fixed FCV. This value is in the ISO 8601 timestamp format (e.g. 2024-12-04T16:25:00Z). Note that this field cannot exceed 4 weeks from the pinned date.
- `version` (String) Feature compatibility version of the cluster.


<a id="nestedatt--replication_specs"></a>
### Nested Schema for `replication_specs`

Read-Only:

- `container_id` (Map of String) A key-value map of the Network Peering Container ID(s) for the configuration specified in region_configs. The Container ID is the id of the container created when the first cluster in the region (AWS/Azure) or project (GCP) was created.
- `external_id` (String) Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. This value corresponds to Shard ID displayed in the UI.
- `region_configs` (Attributes List) Hardware specifications for nodes set for a given region. Each **regionConfigs** object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each **regionConfigs** object must have either an **analyticsSpecs** object, **electableSpecs** object, or **readOnlySpecs** object. Tenant clusters only require **electableSpecs. Dedicated** clusters can specify any of these specifications, but must have at least one **electableSpecs** object within a **replicationSpec**.

**Example:**

If you set `"replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize" : "M30"`, set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : `"M30"` if you have electable nodes and `"replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize" : `"M30"` if you have read-only nodes. (see [below for nested schema](#nestedatt--replication_specs--region_configs))
- `zone_id` (String) Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.
- `zone_name` (String) Human-readable label that describes the zone this shard belongs to in a Global Cluster. Provide this value only if "clusterType" : "GEOSHARDED" but not "selfManagedSharding" : true.

<a id="nestedatt--replication_specs--region_configs"></a>
### Nested Schema for `replication_specs.region_configs`

Read-Only:

- `analytics_auto_scaling` (Attributes) Options that determine how this cluster handles resource scaling. (see [below for nested schema](#nestedatt--replication_specs--region_configs--analytics_auto_scaling))
- `analytics_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--analytics_specs))
- `auto_scaling` (Attributes) Options that determine how this cluster handles resource scaling. (see [below for nested schema](#nestedatt--replication_specs--region_configs--auto_scaling))
- `backing_provider_name` (String) Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when **providerName** is `TENANT` and **electableSpecs.instanceSize** is `M0`.
- `effective_analytics_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--effective_analytics_specs))
- `effective_electable_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--effective_electable_specs))
- `effective_read_only_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--effective_read_only_specs))
- `electable_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--electable_specs))
- `priority` (Number) Precedence is given to this region when a primary election occurs. If your **regionConfigs** has only **readOnlySpecs**, **analyticsSpecs**, or both, set this value to `0`. If you have multiple **regionConfigs** objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.

**Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.
- `provider_name` (String) Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.
- `read_only_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--read_only_specs))
- `region_name` (String) Physical location of your MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. The region name is only returned in the response for single-region clusters. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. It assigns the VPC a Classless Inter-Domain Routing (CIDR) block. To limit a new VPC peering connection to one Classless Inter-Domain Routing (CIDR) block and region, create the connection first. Deploy the cluster after the connection starts. GCP Clusters and Multi-region clusters require one VPC peering connection for each region. MongoDB nodes can use only the peering connection that resides in the same region as the nodes to communicate with the peered VPC.

<a id="nestedatt--replication_specs--region_configs--analytics_auto_scaling"></a>
### Nested Schema for `replication_specs.region_configs.analytics_auto_scaling`

Read-Only:

- `compute_enabled` (Boolean) Flag that indicates whether someone enabled instance size auto-scaling.

- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.
- Set to `false` to disable instance size automatic scaling.
- `compute_max_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_min_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_scale_down_enabled` (Boolean) Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.
- `disk_gb_enabled` (Boolean) Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.


<a id="nestedatt--replication_specs--region_configs--analytics_specs"></a>
### Nested Schema for `replication_specs.region_configs.analytics_specs`

Read-Only:

- `disk_iops` (Number) Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:

- set `"replicationSpecs[n].regionConfigs[m].providerName" : "Azure"`.
- set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : "M40"` or greater not including `Mxx_NVME` tiers.

The maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.
This parameter defaults to the cluster tier's standard IOPS value.
Changing this value impacts cluster cost.
- `disk_size_gb` (Number) Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.

 This value must be equal for all shards and node types.

 This value is not configurable on M0/M2/M5 clusters.

 MongoDB Cloud requires this parameter if you set **replicationSpecs**.

 If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. 

 Storage charge calculations depend on whether you choose the default value or a custom value.

 The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
- `ebs_volume_type` (String) Type of storage you want to attach to your AWS-provisioned cluster.

- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. 

- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.
- `instance_size` (String) Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
- `node_count` (Number) Number of nodes of the given type for MongoDB Cloud to deploy to the region.


<a id="nestedatt--replication_specs--region_configs--auto_scaling"></a>
### Nested Schema for `replication_specs.region_configs.auto_scaling`

Read-Only:

- `compute_enabled` (Boolean) Flag that indicates whether someone enabled instance size auto-scaling.

- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.
- Set to `false` to disable instance size automatic scaling.
- `compute_max_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_min_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_scale_down_enabled` (Boolean) Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.
- `disk_gb_enabled` (Boolean) Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.


<a id="nestedatt--replication_specs--region_configs--effective_analytics_specs"></a>
### Nested Schema for `replication_specs.region_configs.effective_analytics_specs`

Read-Only:

- `disk_iops` (Number) Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:

- set `"replicationSpecs[n].regionConfigs[m].providerName" : "Azure"`.
- set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : "M40"` or greater not including `Mxx_NVME` tiers.

The maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.
This parameter defaults to the cluster tier's standard IOPS value.
Changing this value impacts cluster cost.
- `disk_size_gb` (Number) Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.

 This value must be equal for all shards and node types.

 This value is not configurable on M0/M2/M5 clusters.

 MongoDB Cloud requires this parameter if you set **replicationSpecs**.

 If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. 

 Storage charge calculations depend on whether you choose the default value or a custom value.

 The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
- `ebs_volume_type` (String) Type of storage you want to attach to your AWS-provisioned cluster.

- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. 

- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.
- `instance_size` (String) Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
- `node_count` (Number) Number of nodes of the given type for MongoDB Cloud to deploy to the region.


<a id="nestedatt--replication_specs--region_configs--effective_electable_specs"></a>
### Nested Schema for `replication_specs.region_configs.effective_electable_specs`

Read-Only:

- `disk_iops` (Number) Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:

- set `"replicationSpecs[n].regionConfigs[m].providerName" : "Azure"`.
- set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : "M40"` or greater not including `Mxx_NVME` tiers.

The maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.
This parameter defaults to the cluster tier's standard IOPS value.
Changing this value impacts cluster cost.
- `disk_size_gb` (Number) Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.

 This value must be equal for all shards and node types.

 This value is not configurable on M0/M2/M5 clusters.

 MongoDB Cloud requires this parameter if you set **replicationSpecs**.

 If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. 

 Storage charge calculations depend on whether you choose the default value or a custom value.

 The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
- `ebs_volume_type` (String) Type of storage you want to attach to your AWS-provisioned cluster.

- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. 

- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.
- `instance_size` (String) Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
- `node_count` (Number) Number of nodes of the given type for MongoDB Cloud to deploy to the region.


<a id="nestedatt--replication_specs--region_configs--effective_read_only_specs"></a>
### Nested Schema for `replication_specs.region_configs.effective_read_only_specs`

Read-Only:

- `disk_iops` (Number) Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:

- set `"replicationSpecs[n].regionConfigs[m].providerName" : "Azure"`.
- set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : "M40"` or greater not including `Mxx_NVME` tiers.

The maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.
This parameter defaults to the cluster tier's standard IOPS value.
Changing this value impacts cluster cost.
- `disk_size_gb` (Number) Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.

 This value must be equal for all shards and node types.

 This value is not configurable on M0/M2/M5 clusters.

 MongoDB Cloud requires this parameter if you set **replicationSpecs**.

 If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. 

 Storage charge calculations depend on whether you choose the default value or a custom value.

 The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
- `ebs_volume_type` (String) Type of storage you want to attach to your AWS-provisioned cluster.

- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. 

- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.
- `instance_size` (String) Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
- `node_count` (Number) Number of nodes of the given type for MongoDB Cloud to deploy to the region.


<a id="nestedatt--replication_specs--region_configs--electable_specs"></a>
### Nested Schema for `replication_specs.region_configs.electable_specs`

Read-Only:

- `disk_iops` (Number) Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:

- set `"replicationSpecs[n].regionConfigs[m].providerName" : "Azure"`.
- set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : "M40"` or greater not including `Mxx_NVME` tiers.

The maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.
This parameter defaults to the cluster tier's standard IOPS value.
Changing this value impacts cluster cost.
- `disk_size_gb` (Number) Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.

 This value must be equal for all shards and node types.

 This value is not configurable on M0/M2/M5 clusters.

 MongoDB Cloud requires this parameter if you set **replicationSpecs**.

 If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. 

 Storage charge calculations depend on whether you choose the default value or a custom value.

 The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
- `ebs_volume_type` (String) Type of storage you want to attach to your AWS-provisioned cluster.

- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. 

- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.
- `instance_size` (String) Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
- `node_count` (Number) Number of nodes of the given type for MongoDB Cloud to deploy to the region.


<a id="nestedatt--replication_specs--region_configs--read_only_specs"></a>
### Nested Schema for `replication_specs.region_configs.read_only_specs`

Read-Only:

- `disk_iops` (Number) Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:

- set `"replicationSpecs[n].regionConfigs[m].providerName" : "Azure"`.
- set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : "M40"` or greater not including `Mxx_NVME` tiers.

The maximum input/output operations per second (IOPS) depend on the selected **.instanceSize** and **.diskSizeGB**.
This parameter defaults to the cluster tier's standard IOPS value.
Changing this value impacts cluster cost.
- `disk_size_gb` (Number) Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.

 This value must be equal for all shards and node types.

 This value is not configurable on M0/M2/M5 clusters.

 MongoDB Cloud requires this parameter if you set **replicationSpecs**.

 If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. 

 Storage charge calculations depend on whether you choose the default value or a custom value.

 The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
- `ebs_volume_type` (String) Type of storage you want to attach to your AWS-provisioned cluster.

- `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size. 

- `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.
- `instance_size` (String) Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
- `node_count` (Number) Number of nodes of the given type for MongoDB Cloud to deploy to the region.
