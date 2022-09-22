---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: advanced_cluster"
sidebar_current: "docs-mongodbatlas-resource-advanced-cluster"
description: |-
    Provides an Advanced Cluster resource.
---

# Resource: mongodbatlas_advanced_cluster

`mongodbatlas_advanced_cluster` provides an Advanced Cluster resource. The resource lets you create, edit and delete advanced clusters. The resource requires your Project ID.

More information on considerations for using advanced clusters please see [Considerations](https://docs.atlas.mongodb.com/reference/api/cluster-advanced/create-one-cluster-advanced/#considerations)

~> **IMPORTANT:**
<br> &#8226; The primary difference between [`mongodbatlas_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) and [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) is that `mongodbatlas_advanced_cluster` supports multi-cloud clusters.  We recommend new users start with the `mongodbatlas_advanced_cluster` resource.  

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** A network container is created for a advanced cluster to reside in if one does not yet exist in the project.  To  use this automatically created container with another resource, such as peering, the `container_id` is exported after creation.

## Example Usage


### Example single provider and single region

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"
  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 1
      region_name   = "US_EAST_1"
    }
  }
}
```

### Example Tenant Cluster

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M5"
      }
      provider_name         = "TENANT"
      backing_provider_name = "AWS"
      region_name           = "US_EAST_1"
      priority              = 1
    }
  }
}
```

### Example Multicloud.

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 1
      region_name   = "US_EAST_1"
    }
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      provider_name = "GCP"
      priority      = 6
      region_name   = "NORTH_AMERICA_NORTHEAST_1"
    }
  }
}
```

## Argument Reference

* `project_id` - (Required) Unique ID for the project to create the database user.
* `name` - (Required) Name of the cluster as it appears in Atlas. Once the cluster is created, its name cannot be changed.
* `cluster_type` - (Required) Atlas provides different instance sizes, each with a default storage capacity and RAM size. The instance size you select is used for all the data-bearing servers in your cluster. See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/) `providerSettings.instanceSizeName` for valid values and default resources. 

* `backup_enabled` - (Optional) Flag that indicates whether the cluster can perform backups.
  If `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters.

  Backup uses:
  [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/#std-label-backup-cloud-provider) for dedicated clusters.
  [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/#std-label-m2-m5-snapshots) for tenant clusters.
  If "`backup_enabled`" : `false`, the cluster doesn't use Atlas backups.

This parameter defaults to false.

* `bi_connector` - (Optional) Configuration settings applied to BI Connector for Atlas on this cluster. The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger clusters. The BI Connector is a powerful tool which provides users SQL-based access to their MongoDB databases. As a result, the BI Connector performs operations which may be CPU and memory intensive. Given the limited hardware resources on M10 and M20 cluster tiers, you may experience performance degradation of the cluster when enabling the BI Connector. If this occurs, upgrade to an M30 or larger cluster or disable the BI Connector. See [below](#bi_connector).
* `cluster_type` - (Required)Type of the cluster that you want to create.
    Accepted values include:
      - `REPLICASET` Replica set
      - `SHARDED`	Sharded cluster
      - `GEOSHARDED` Global Cluster

* `disk_size_gb` - (Optional) Capacity, in gigabytes, of the host's root volume. Increase this number to add capacity, up to a maximum possible value of 4096 (i.e., 4 TB). This value must be a positive number. You can't set this value with clusters with local [NVMe SSDs](https://docs.atlas.mongodb.com/cluster-tier/#std-label-nvme-storage). The minimum disk size for dedicated clusters is 10 GB for AWS and GCP. If you specify diskSizeGB with a lower disk size, Atlas defaults to the minimum disk size value. If your cluster includes Azure nodes, this value must correspond to an existing Azure disk type (8, 16, 32, 64, 128, 256, 512, 1024, 2048, or 4095)Atlas calculates storage charges differently depending on whether you choose the default value or a custom value. The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require additional storage space beyond this limitation, consider [upgrading your cluster](https://docs.atlas.mongodb.com/scale-cluster/#std-label-scale-cluster-instance) to a higher tier. If your cluster spans cloud service providers, this value defaults to the minimum default of the providers involved.
* `encryption_at_rest_provider` - (Optional) Possible values are AWS, GCP, AZURE or NONE.  Only needed if you desire to manage the keys, see [Encryption at Rest using Customer Key Management](https://docs.atlas.mongodb.com/security-kms-encryption/) for complete documentation.  You must configure encryption at rest for the Atlas project before enabling it on any cluster in the project. For Documentation, see [AWS](https://docs.atlas.mongodb.com/security-aws-kms/), [GCP](https://docs.atlas.mongodb.com/security-kms-encryption/) and [Azure](https://docs.atlas.mongodb.com/security-azure-kms/#std-label-security-azure-kms). Requirements are if `replication_specs.#.region_configs.#.<type>Specs.instance_size` is M10 or greater and `backup_enabled` is false or omitted.   
* `labels` - (Optional) Configuration for the collection of key-value pairs that tag and categorize the cluster. See [below](#labels).
* `mongo_db_major_version` - (Optional) Version of the cluster to deploy. Atlas supports the following MongoDB versions for M10+ clusters: `4.0`, `4.2`, `4.4`, or `5.0`. If omitted, Atlas deploys a cluster that runs MongoDB 4.4. If `replication_specs#.region_configs#.<type>Specs.instance_size`: `M0`, `M2` or `M5`, Atlas deploys MongoDB 4.4. Atlas always deploys the cluster with the latest stable release of the specified version.  If you set a value to this parameter and set `version_release_system` `CONTINUOUS`, the resource returns an error. Either clear this parameter or set `version_release_system`: `LTS`.
* `pit_enabled` - (Optional) - Flag that indicates if the cluster uses Continuous Cloud Backup.
* `replication_specs` - Configuration for cluster regions and the hardware provisioned in them. See [below](#replication_specs)
* `root_cert_type` - (Optional) - Certificate Authority that MongoDB Atlas clusters use. You can specify ISRGROOTX1 (for ISRG Root X1).
* `version_release_system` - (Optional) - Release cadence that Atlas uses for this cluster. This parameter defaults to `LTS`. If you set this field to `CONTINUOUS`, you must omit the `mongo_db_major_version` field. Atlas accepts:
  - `CONTINUOUS`:  Atlas creates your cluster using the most recent MongoDB release. Atlas automatically updates your cluster to the latest major and rapid MongoDB releases as they become available.
  - `LTS`: Atlas creates your cluster using the latest patch release of the MongoDB version that you specify in the mongoDBMajorVersion field. Atlas automatically updates your cluster to subsequent patch releases of this MongoDB version. Atlas doesn't update your cluster to newer rapid or major MongoDB releases as they become available.
* `paused` (Optional) - Flag that indicates whether the cluster is paused or not. You can pause M10 or larger clusters.  You cannot initiate pausing for a shared/tenant tier cluster.  See [Considerations for Paused Clusters](https://docs.atlas.mongodb.com/pause-terminate-cluster/#considerations-for-paused-clusters)  
  **NOTE** Pause lasts for up to 30 days. If you don't resume the cluster within 30 days, Atlas resumes the cluster.  When the cluster resumption happens Terraform will flag the changed state.  If you wish to keep the cluster paused, reapply your Terraform configuration.   If you prefer to allow the automated change of state to unpaused use:
  `lifecycle {
  ignore_changes = [paused]
  }`
 

### bi_connector

Specifies BI Connector for Atlas configuration.

 ```terraform
 bi_connector {
  enabled         = true
  read_preference = secondary
}
  ```

* `enabled` - (Optional) Specifies whether or not BI Connector for Atlas is enabled on the cluster.l
*
  - Set to `true` to enable BI Connector for Atlas.
  - Set to `false` to disable BI Connector for Atlas.

* `read_preference` - (Optional) Specifies the read preference to be used by BI Connector for Atlas on the cluster. Each BI Connector for Atlas read preference contains a distinct combination of [readPreference](https://docs.mongodb.com/manual/core/read-preference/) and [readPreferenceTags](https://docs.mongodb.com/manual/core/read-preference/#tag-sets) options. For details on BI Connector for Atlas read preferences, refer to the [BI Connector Read Preferences Table](https://docs.atlas.mongodb.com/tutorial/create-global-writes-cluster/#bic-read-preferences).

  - Set to "primary" to have BI Connector for Atlas read from the primary.

  - Set to "secondary" to have BI Connector for Atlas read from a secondary member. Default if there are no analytics nodes in the cluster.

  - Set to "analytics" to have BI Connector for Atlas read from an analytics node. Default if the cluster contains analytics nodes.

### Advanced Configuration Options

-> **NOTE:** Prior to setting these options please ensure you read https://docs.atlas.mongodb.com/cluster-config/additional-options/.

-> **NOTE:** This argument has been changed to type list make sure you have the proper syntax. The list can have only one  item maximum.

Include **desired options** within advanced_configuration:

```terraform
// Nest options within advanced_configuration
 advanced_configuration {
   javascript_enabled                   = false
   minimum_enabled_tls_protocol         = "TLS1_2"
 }
```

* `default_read_concern` - (Optional) [Default level of acknowledgment requested from MongoDB for read operations](https://docs.mongodb.com/manual/reference/read-concern/) set for this cluster. MongoDB 4.4 clusters default to [available](https://docs.mongodb.com/manual/reference/read-concern-available/).
* `default_write_concern` - (Optional) [Default level of acknowledgment requested from MongoDB for write operations](https://docs.mongodb.com/manual/reference/write-concern/) set for this cluster. MongoDB 4.4 clusters default to [1](https://docs.mongodb.com/manual/reference/write-concern/).
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


### labels

 ```terraform
 labels {
        key   = "Key 1"
        value = "Value 1"
  }
 labels {
        key   = "Key 2"
        value = "Value 2"
  }
```

Key-value pairs that tag and categorize the cluster. Each key and value has a maximum length of 255 characters.  You cannot set the key `Infrastructure Tool`, it is used for internal purposes to track aggregate usage.

* `key` - The key that you want to write.
* `value` - The value that you want to write.


### replication_specs 

```terraform
//Example Multicloud
replication_specs {
  region_configs {
    electable_specs {
      instance_size = "M10"
      node_count    = 3
    }
    analytics_specs {
      instance_size = "M10"
      node_count    = 1
    }
    provider_name = "AWS"
    priority      = 1
    region_name   = "US_EAST_1"
  }
  region_configs {
    electable_specs {
      instance_size = "M10"
      node_count    = 2
    }
    provider_name = "GCP"
    priority      = 6
    region_name   = "NORTH_AMERICA_NORTHEAST_1"
  }
}
```

* `num_shards` - (Required) Provide this value if you set a `cluster_type` of SHARDED or GEOSHARDED. Omit this value if you selected a `cluster_type` of REPLICASET. This API resource accepts 1 through 50, inclusive. This parameter defaults to 1. If you specify a `num_shards` value of 1 and a `cluster_type` of SHARDED, Atlas deploys a single-shard [sharded cluster](https://docs.atlas.mongodb.com/reference/glossary/#std-term-sharded-cluster). Don't create a sharded cluster with a single shard for production environments. Single-shard sharded clusters don't provide the same benefits as multi-shard configurations.
* `region_configs` - (Optional) Configuration for the hardware specifications for nodes set for a given regionEach `region_configs` object describes the region's priority in elections and the number and type of MongoDB nodes that Atlas deploys to the region. Each `region_configs` object must have either an `analytics_specs` object, `electable_specs` object, or `read_only_specs` object. See [below](#region_configs)
*  `container_id` - A key-value map of the Network Peering Container ID(s) for the configuration specified in `region_configs`. The Container ID is the id of the container either created programmatically by the user before any clusters existed in a project or when the first cluster in the region (AWS/Azure) or project (GCP) was created.  The syntax is `"providerName:regionName" = "containerId"`. Example `AWS:US_EAST_1" = "61e0797dde08fb498ca11a71`.
* `zone_name` - (Optional) Name for the zone in a Global Cluster.


### region_configs

* `analytics_specs` - (Optional) Hardware specifications for [analytics nodes](https://docs.atlas.mongodb.com/reference/faq/deployment/#std-label-analytics-nodes-overview) needed in the region. Analytics nodes handle analytic data such as reporting queries from BI Connector for Atlas. Analytics nodes are read-only and can never become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary). If you don't specify this parameter, no analytics nodes deploy to this region. See [below](#specs)
* `auto_scaling` - (Optional) Configuration for the Collection of settings that configures auto-scaling information for the cluster. The values for the `auto_scaling` parameter must be the same for every item in the `replication_specs` array. See [below](#auto_scaling)
* `backing_provider_name` - (Optional) Cloud service provider on which you provision the host for a multi-tenant cluster. Use this only when a `provider_name` is `TENANT` and `instance_size` of a specs is `M2` or `M5`.
* `electable_specs` - (Optional) Hardware specifications for electable nodes in the region. Electable nodes can become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary) and can enable local reads. If you do not specify this option, no electable nodes are deployed to the region. See [below](#specs)
* `priority` - (Optional)  Election priority of the region. For regions with only read-only nodes, set this value to 0.
  * If you have multiple `region_configs` objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is 7.
  * If your region has set `region_configs.#.electable_specs.0.node_count` to 1 or higher, it must have a priority of exactly one (1) less than another region in the `replication_specs.#.region_configs.#` array. The highest-priority region must have a priority of 7. The lowest possible priority is 1.
* `provider_name` - (Optional) Cloud service provider on which the servers are provisioned.
  The possible values are:

  - `AWS` - Amazon AWS
  - `GCP` - Google Cloud Platform
  - `AZURE` - Microsoft Azure
  - `TENANT` - M2 or M5 multi-tenant cluster. Use `replication_specs.#.region_configs.#.backing_provider_name` to set the cloud service provider.
* `read_only_specs` - (Optional) Hardware specifications for read-only nodes in the region. Read-only nodes can become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary) and can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region. See [below](#specs)
* `region_name` - (Optional) Physical location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases.  Requires the **Atlas region name**, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).

### specs

* `disk_iops` - (Optional) Target throughput (IOPS) desired for AWS storage attached to your cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster.
* `ebs_volume_type` - (Optional) Type of storage you want to attach to your AWS-provisioned cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster. Valid values are:
    * `STANDARD` volume types can't exceed the default IOPS rate for the selected volume size.
    * `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size.
* `instance_size` - (Optional) Hardware specification for the instance sizes in this region. Each instance size has a default storage and memory capacity. The instance size you select applies to all the data-bearing hosts in your instance size.
* `node_count` - (Optional) Number of read-only nodes for Atlas to deploy to the region. Read-only nodes can never become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary), but can enable local reads.

### auto_scaling

* `disk_gb_enabled` - (Optional) Flag that indicates whether this cluster enables disk auto-scaling. This parameter defaults to true.
* `compute_enabled` - (Optional) Flag that indicates whether instance size auto-scaling is enabled. This parameter defaults to false.

~> **IMPORTANT:** If `compute_enabled` is true,  then Atlas will automatically scale up to the maximum provided and down to the minimum, if provided.
This will cause the value of `instance_size` returned to potential be different than what is specified in the Terraform config and if one then applies a plan, not noting this, Terraform will scale the cluster back down to the original `instance_size` value.
To prevent this a lifecycle customization should be used, i.e.:  
`lifecycle {
  ignore_changes = [instance_size]
}`
But in order to explicitly change `instance_size` comment the `lifecycle` block and run `terraform apply`. Please ensure to uncomment it to prevent any accidental changes.

* `compute_scale_down_enabled` - (Optional) Flag that indicates whether the instance size may scale down. Atlas requires this parameter if `replication_specs.#.region_configs.#.auto_scaling.0.compute_enabled` : true. If you enable this option, specify a value for `replication_specs.#.region_configs.#.auto_scaling.0.compute_min_instance_size`.
* `compute_min_instance_size` - (Optional) Minimum instance size to which your cluster can automatically scale (such as M10). Atlas requires this parameter if `replication_specs.#.region_configs.#.auto_scaling.0.compute_scale_down_enabled` is true.
* `compute_max_instance_size` - (Optional) Maximum instance size to which your cluster can automatically scale (such as M40). Atlas requires this parameter if `replication_specs.#.region_configs.#.auto_scaling.0.compute_enabled` is true.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - The cluster ID.
*  `mongo_db_version` - Version of MongoDB the cluster runs, in `major-version`.`minor-version` format.
* `id` -	The Terraform's unique identifier used internally for state management.
* `connection_strings` - Set of connection strings that your applications use to connect to this cluster. More info in [Connection-strings](https://docs.mongodb.com/manual/reference/connection-string/). Use the parameters in this object to connect your applications to this cluster. To learn more about the formats of connection strings, see [Connection String Options](https://docs.atlas.mongodb.com/reference/faq/connection-changes/). NOTE: Atlas returns the contents of this object after the cluster is operational, not while it builds the cluster.

   **NOTE** Connection strings must be returned as a list, therefore to refer to a specific attribute value add index notation. Example: mongodbatlas_advanced_cluster.cluster-test.connection_strings.0.standard_srv

   Private connection strings may not be available immediately as the reciprocal connections may not have finalized by end of the Terraform run. If the expected connection string(s) do not contain a value a terraform refresh may need to be performed to obtain the value. One can also view the status of the peered connection in the [Atlas UI](https://docs.atlas.mongodb.com/security-vpc-peering/). 

    - `connection_strings.standard` -   Public mongodb:// connection string for this cluster.
    - `connection_strings.standard_srv` - Public mongodb+srv:// connection string for this cluster. The mongodb+srv protocol tells the driver to look up the seed list of hosts in DNS. Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don’t need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn’t  , use connectionStrings.standard.
    - `connection_strings.aws_private_link` -  [Private-endpoint-aware](https://docs.atlas.mongodb.com/security-private-endpoint/#private-endpoint-connection-strings) mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a AWS PrivateLink connection to this cluster. **DEPRECATED** Use `connection_strings.private_endpoint[n].connection_string` instead.
    - `connection_strings.aws_private_link_srv` - [Private-endpoint-aware](https://docs.atlas.mongodb.com/security-private-endpoint/#private-endpoint-connection-strings) mongodb+srv://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a AWS PrivateLink connection to this cluster. Use this URI format if your driver supports it. If it doesn’t, use connectionStrings.awsPrivateLink. **DEPRECATED** Use `connection_strings.private_endpoint[n].srv_connection_string` instead.
    - `connection_strings.private` -   [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_srv` -  [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb+srv://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_endpoint` - Private endpoint connection strings. Each object describes the connection strings you can use to connect to this cluster through a private endpoint. Atlas returns this parameter only if you deployed a private endpoint to all regions to which you deployed this cluster's nodes.
    - `connection_strings.private_endpoint.#.connection_string` - Private-endpoint-aware `mongodb://`connection string for this private endpoint.
    - `connection_strings.private_endpoint.#.srv_connection_string` - Private-endpoint-aware `mongodb+srv://` connection string for this private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in DNS . Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don't need to: Append the seed list or Change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use `connection_strings.private_endpoint[n].connection_string`
    - `connection_strings.private_endpoint.#.type` - Type of MongoDB process that you connect to with the connection strings. Atlas returns `MONGOD` for replica sets, or `MONGOS` for sharded clusters.
    - `connection_strings.private_endpoint.#.endpoints` - Private endpoint through which you connect to Atlas when you use `connection_strings.private_endpoint[n].connection_string` or `connection_strings.private_endpoint[n].srv_connection_string`
    - `connection_strings.private_endoint.#.endpoints.#.endpoint_id` - Unique identifier of the private endpoint.
    - `connection_strings.private_endpoint.#.endpoints.#.provider_name` - Cloud provider to which you deployed the private endpoint. Atlas returns `AWS` or `AZURE`.
    - `connection_strings.private_endpoint.#.endpoints.#.region` - Region to which you deployed the private endpoint.
* `state_name` - Current state of the cluster. The possible states are:
    - IDLE
    - CREATING
    - UPDATING
    - DELETING
    - DELETED
    - REPAIRING

## Import

Clusters can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_advanced_cluster.my_cluster 1112222b3bf99403840e8934-Cluster0
```

See detailed information for arguments and attributes: [MongoDB API Advanced Clusters](https://docs.atlas.mongodb.com/reference/api/cluster-advanced/create-one-cluster-advanced/)
