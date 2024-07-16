# Resource: mongodbatlas_advanced_cluster

`mongodbatlas_advanced_cluster` provides an Advanced Cluster resource. The resource lets you create, edit and delete advanced clusters. The resource requires your Project ID.

More information on considerations for using advanced clusters please see [Considerations](https://docs.atlas.mongodb.com/reference/api/cluster-advanced/create-one-cluster-advanced/#considerations)

~> **IMPORTANT:**
<br> &#8226; The primary difference between [`mongodbatlas_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) and [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) is that `mongodbatlas_advanced_cluster` supports multi-cloud clusters.  We recommend new users start with the `mongodbatlas_advanced_cluster` resource.  

-> **NOTE:** If Backup Compliance Policy is enabled for the project for which this backup schedule is defined, you cannot modify the backup schedule for an individual cluster below the minimum requirements set in the Backup Compliance Policy.  See [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).


<br> &#8226; Upgrading the shared tier is supported. Any change from a shared tier cluster (a tenant) to a different instance size will be considered a tenant upgrade. When upgrading from the shared tier, change the `provider_name` from "TENANT" to your preferred provider (AWS, GCP or Azure) and remove the variable `backing_provider_name`.  See the [Example Tenant Cluster Upgrade](#Example-Tenant-Cluster-Upgrade) below.   Note you can upgrade a shared tier cluster only to a single provider on an M10-tier cluster or greater.  
<br> &#8226; **IMPORTANT NOTE** When upgrading from the shared tier, *only* the upgrade changes will be applied. This helps avoid a corrupt state file in the event that the upgrade succeeds but subsequent updates fail within the same `terraform apply`. To apply additional cluster changes, run a secondary `terraform apply` after the upgrade succeeds.
-> **NOTE:** Groups and projects are synonymous terms. You might find group_id in the official documentation.

-> **NOTE:** A network container is created for each provider/region combination on the advanced cluster. This can be referenced via a computed attribute for use with other resources. Refer to the `replication_specs.#.container_id` attribute in the [Attributes Reference](#attributes_reference) for more information.

-> **NOTE:** To enable Cluster Extended Storage Sizes use the `is_extended_storage_sizes_enabled` parameter in the [mongodbatlas_project resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project).

-> **NOTE:** The Low-CPU instance clusters are prefixed with `R`, i.e. `R40`. For complete list of Low-CPU instance clusters see Cluster Configuration Options under each Cloud Provider (https://www.mongodb.com/docs/atlas/reference/cloud-providers/).


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
      priority      = 7
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
      priority              = 7
    }
  }
}
```

### Example Tenant Cluster Upgrade

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
      }
      provider_name         = "AWS"
      region_name           = "US_EAST_1"
      priority              = 7
    }
  }
}
```

### Example Multi-Cloud Cluster
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
      priority      = 7
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
### Example of a Multi-Cloud Cluster

```terraform
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = var.cluster_name
  cluster_type = "SHARDED"
  backup_enabled = true

  replication_specs {
    num_shards = 3 

    region_configs { # shard n1
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }

    region_configs { # shard n2
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AZURE"
      priority      = 6
      region_name   = "US_EAST_2"
    }

    region_configs { # shard n3
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "GCP"
      priority      = 5
      region_name   = "US_EAST_4"
    }
  }

  advanced_configuration {
    javascript_enabled                   = true
    oplog_size_mb                        = 30
    sample_refresh_interval_bi_connector = 300
  }
}
```

### Example of a Global Cluster
```terraform
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "GEOSHARDED"
  backup_enabled = true

  replication_specs { # zone n1
    zone_name  = "zone n1"
    num_shards = 3 # 3-shard Multi-Cloud Cluster

    region_configs { # shard n1 
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }

    region_configs { # shard n2
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AZURE"
      priority      = 6
      region_name   = "US_EAST_2"
    }

    region_configs { # shard n3
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "GCP"
      priority      = 5
      region_name   = "US_EAST_4"
    }
  }

  replication_specs { # zone n2
    zone_name  = "zone n2"
    num_shards = 2 # 2-shard Multi-Cloud Cluster

    region_configs { # shard n1 
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }

    region_configs { # shard n2
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AZURE"
      priority      = 6
      region_name   = "EUROPE_NORTH"
    }

    region_configs { # shard n3
      electable_specs {
        instance_size = "M10"
        node_count    = 2
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "GCP"
      priority      = 5
      region_name   = "US_EAST_4"
    }
  }

  advanced_configuration {
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }
}
```


### Example - Return a Connection String
Standard
```terraform
output "standard" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].standard
}
# Example return string: standard = "mongodb://cluster-atlas-shard-00-00.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-01.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-02.ygo1m.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-12diht-shard-0"
```
Standard srv
```terraform
output "standard_srv" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].standard_srv
}
# Example return string: standard_srv = "mongodb+srv://cluster-atlas.ygo1m.mongodb.net"
```
Private with Network peering and Custom DNS AWS enabled
```terraform
output "private" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].private
}
# Example return string: private = "mongodb://cluster-atlas-shard-00-00-pri.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-01-pri.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-02-pri.ygo1m.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-12diht-shard-0"
private = "mongodb+srv://cluster-atlas-pri.ygo1m.mongodb.net"
```
Private srv with Network peering and Custom DNS AWS enabled
```terraform
output "private_srv" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].private_srv
}
# Example return string: private_srv = "mongodb+srv://cluster-atlas-pri.ygo1m.mongodb.net"
```

By endpoint_service_id
```terraform
locals {
  endpoint_service_id = google_compute_network.default.name
  private_endpoints   = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings : cs.private_endpoint]), [])
  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id)
  ]
}
output "endpoint_service_connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
# Example return string: connection_string = "mongodb+srv://cluster-atlas-pl-0.ygo1m.mongodb.net"
```
Refer to the following for full privatelink endpoint connection string examples:
* [GCP Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/gcp-atlas-privatelink)
* [Azure Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/azure-atlas-privatelink)
* [AWS, Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/aws-privatelink-endpoint/cluster)
* [AWS, Regionalized Private Endpoints](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/aws-privatelink-endpoint/cluster-geosharded)

## Argument Reference

* `project_id` - (Required) Unique ID for the project to create the database user.
* `name` - (Required) Name of the cluster as it appears in Atlas. Once the cluster is created, its name cannot be changed. **WARNING** Changing the name will result in destruction of the existing cluster and the creation of a new cluster.

* `backup_enabled` - (Optional) Flag that indicates whether the cluster can perform backups.
  If `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters.

  Backup uses:
  [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/#std-label-backup-cloud-provider) for dedicated clusters.
  [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/#std-label-m2-m5-snapshots) for tenant clusters.
  If "`backup_enabled`" : `false`, the cluster doesn't use Atlas backups.

This parameter defaults to false.

* `retain_backups_enabled` - (Optional) Set to true to retain backup snapshots for the deleted cluster. M10 and above only.

**NOTE** Prior version of provider had parameter as `bi_connector` state will migrate it to new value you only need to update parameter in your terraform file

* `bi_connector_config` - (Optional) Configuration settings applied to BI Connector for Atlas on this cluster. The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger clusters. The BI Connector is a powerful tool which provides users SQL-based access to their MongoDB databases. As a result, the BI Connector performs operations which may be CPU and memory intensive. Given the limited hardware resources on M10 and M20 cluster tiers, you may experience performance degradation of the cluster when enabling the BI Connector. If this occurs, upgrade to an M30 or larger cluster or disable the BI Connector. See [below](#bi_connector_config).
* `cluster_type` - (Required)Type of the cluster that you want to create.
    Accepted values include:
      - `REPLICASET` Replica set
      - `SHARDED`	Sharded cluster
      - `GEOSHARDED` Global Cluster

* `disk_size_gb` - (Optional) Capacity, in gigabytes, of the host's root volume. Increase this number to add capacity, up to a maximum possible value of 4096 (i.e., 4 TB). This value must be a positive number. You can't set this value with clusters with local [NVMe SSDs](https://docs.atlas.mongodb.com/cluster-tier/#std-label-nvme-storage). The minimum disk size for dedicated clusters is 10 GB for AWS and GCP. If you specify diskSizeGB with a lower disk size, Atlas defaults to the minimum disk size value. If your cluster includes Azure nodes, this value must correspond to an existing Azure disk type (8, 16, 32, 64, 128, 256, 512, 1024, 2048, or 4095)Atlas calculates storage charges differently depending on whether you choose the default value or a custom value. The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require additional storage space beyond this limitation, consider [upgrading your cluster](https://docs.atlas.mongodb.com/scale-cluster/#std-label-scale-cluster-instance) to a higher tier. If your cluster spans cloud service providers, this value defaults to the minimum default of the providers involved.
* `encryption_at_rest_provider` - (Optional) Possible values are AWS, GCP, AZURE or NONE.  Only needed if you desire to manage the keys, see [Encryption at Rest using Customer Key Management](https://docs.atlas.mongodb.com/security-kms-encryption/) for complete documentation.  You must configure encryption at rest for the Atlas project before enabling it on any cluster in the project. For Documentation, see [AWS](https://docs.atlas.mongodb.com/security-aws-kms/), [GCP](https://docs.atlas.mongodb.com/security-kms-encryption/) and [Azure](https://docs.atlas.mongodb.com/security-azure-kms/#std-label-security-azure-kms). Requirements are if `replication_specs.#.region_configs.#.<type>Specs.instance_size` is M10 or greater and `backup_enabled` is false or omitted.   
* `tags` - (Optional) Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#tags).
* `labels` - (Optional) Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#labels). **DEPRECATED** Use `tags` instead.
* `mongo_db_major_version` - (Optional) Version of the cluster to deploy. Atlas supports the following MongoDB versions for M10+ clusters: `4.4`, `5.0`, `6.0` or `7.0`. If omitted, Atlas deploys a cluster that runs MongoDB 7.0. If `replication_specs#.region_configs#.<type>Specs.instance_size`: `M0`, `M2` or `M5`, Atlas deploys MongoDB 4.4. Atlas always deploys the cluster with the latest stable release of the specified version.  If you set a value to this parameter and set `version_release_system` `CONTINUOUS`, the resource returns an error. Either clear this parameter or set `version_release_system`: `LTS`.
* `pit_enabled` - (Optional) - Flag that indicates if the cluster uses Continuous Cloud Backup.
* `replication_specs` - Configuration for cluster regions and the hardware provisioned in them. See [below](#replication_specs)
* `root_cert_type` - (Optional) - Certificate Authority that MongoDB Atlas clusters use. You can specify ISRGROOTX1 (for ISRG Root X1).
* `termination_protection_enabled` - Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
* `version_release_system` - (Optional) - Release cadence that Atlas uses for this cluster. This parameter defaults to `LTS`. If you set this field to `CONTINUOUS`, you must omit the `mongo_db_major_version` field. Atlas accepts:
  - `CONTINUOUS`:  Atlas creates your cluster using the most recent MongoDB release. Atlas automatically updates your cluster to the latest major and rapid MongoDB releases as they become available.
  - `LTS`: Atlas creates your cluster using the latest patch release of the MongoDB version that you specify in the mongoDBMajorVersion field. Atlas automatically updates your cluster to subsequent patch releases of this MongoDB version. Atlas doesn't update your cluster to newer rapid or major MongoDB releases as they become available.
* `paused` (Optional) - Flag that indicates whether the cluster is paused or not. You can pause M10 or larger clusters.  You cannot initiate pausing for a shared/tenant tier cluster. If you try to update a `paused` cluster you will get a `CANNOT_UPDATE_PAUSED_CLUSTER` error. See [Considerations for Paused Clusters](https://docs.atlas.mongodb.com/pause-terminate-cluster/#considerations-for-paused-clusters).
  **NOTE** Pause lasts for up to 30 days. If you don't resume the cluster within 30 days, Atlas resumes the cluster.  When the cluster resumption happens Terraform will flag the changed state.  If you wish to keep the cluster paused, reapply your Terraform configuration.   If you prefer to allow the automated change of state to unpaused use:
  `lifecycle {
  ignore_changes = [paused]
  }`
* `timeouts`- (Optional) The duration of time to wait for Cluster to be created, updated, or deleted. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Advanced Cluster create & delete is `3h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).
* `accept_data_risks_and_force_replica_set_reconfig` - (Optional) If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set `accept_data_risks_and_force_replica_set_reconfig` to the current date. Learn more about Reconfiguring a Replica Set during a regional outage [here](https://dochub.mongodb.org/core/regional-outage-reconfigure-replica-set).
* `global_cluster_self_managed_sharding` - (Optional) Flag that indicates if cluster uses Atlas-Managed Sharding (false, default) or Self-Managed Sharding (true). It can only be enabled for Global Clusters (`GEOSHARDED`). It cannot be changed once the cluster is created. Use this mode if you're an advanced user and the default configuration is too restrictive for your workload. If you select this option, you must manually configure the sharding strategy, more info [here](https://www.mongodb.com/docs/atlas/tutorial/create-global-cluster/#select-your-sharding-configuration).

### bi_connector_config

Specifies BI Connector for Atlas configuration.

```terraform
bi_connector_config {
  enabled         = true
  read_preference = "secondary"
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
* `oplog_min_retention_hours` - (Optional) Minimum retention window for cluster's oplog expressed in hours. A value of null indicates that the cluster uses the default minimum oplog window that MongoDB Cloud calculates.
* **Note**  A minimum oplog retention is required when seeking to change a cluster's class to Local NVMe SSD. To learn more and for latest guidance see [`oplogMinRetentionHours`](https://www.mongodb.com/docs/manual/core/replica-set-oplog/#std-label-replica-set-minimum-oplog-size) 
* `sample_size_bi_connector` - (Optional) Number of documents per database to sample when gathering schema information. Defaults to 100. Available only for Atlas deployments in which BI Connector for Atlas is enabled.
* `sample_refresh_interval_bi_connector` - (Optional) Interval in seconds at which the mongosqld process re-samples data to create its relational schema. The default value is 300. The specified value must be a positive integer. Available only for Atlas deployments in which BI Connector for Atlas is enabled.
* `transaction_lifetime_limit_seconds` - (Optional) Lifetime, in seconds, of multi-document transactions. Defaults to 60 seconds.


### Tags

 ```terraform
 tags {
        key   = "Key 1"
        value = "Value 1"
  }
 tags {
        key   = "Key 2"
        value = "Value 2"
  }
```

Key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.

* `key` - (Required) Constant that defines the set of the tag.
* `value` - (Required) Variable that belongs to the set of the tag.

To learn more, see [Resource Tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas).

### labels

**WARNING:** This property is deprecated and will be removed in the future, use the `tags` attribute instead.

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

Key-value pairs that categorize the cluster. Each key and value has a maximum length of 255 characters.  You cannot set the key `Infrastructure Tool`, it is used for internal purposes to track aggregate usage.

* `key` - The key that you want to write.
* `value` - The value that you want to write.

-> **NOTE:** MongoDB Atlas doesn't display your labels.


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
    priority      = 7
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

* `num_shards` - (Optional) Provide this value if you set a `cluster_type` of SHARDED or GEOSHARDED. Omit this value if you selected a `cluster_type` of REPLICASET. This API resource accepts 1 through 50, inclusive. This parameter defaults to 1. If you specify a `num_shards` value of 1 and a `cluster_type` of SHARDED, Atlas deploys a single-shard [sharded cluster](https://docs.atlas.mongodb.com/reference/glossary/#std-term-sharded-cluster). Don't create a sharded cluster with a single shard for production environments. Single-shard sharded clusters don't provide the same benefits as multi-shard configurations.
If you are upgrading a replica set to a sharded cluster, you cannot increase the number of shards in the same update request. You should wait until after the cluster has completed upgrading to sharded and you have reconnected all application clients to the MongoDB router before adding additional shards. Otherwise, your data might become inconsistent once MongoDB Cloud begins distributing data across shards. To learn more, see [Convert a replica set to a sharded cluster documentation](https://www.mongodb.com/docs/atlas/scale-cluster/#convert-a-replica-set-to-a-sharded-cluster) and [Convert a replica set to a sharded cluster tutorial](https://www.mongodb.com/docs/upcoming/tutorial/convert-replica-set-to-replicated-shard-cluster).
* `region_configs` - (Optional) Configuration for the hardware specifications for nodes set for a given regionEach `region_configs` object describes the region's priority in elections and the number and type of MongoDB nodes that Atlas deploys to the region. Each `region_configs` object must have either an `analytics_specs` object, `electable_specs` object, or `read_only_specs` object. See [below](#region_configs)
* `zone_name` - (Optional) Name for the zone in a Global Cluster.


### region_configs

* `analytics_specs` - (Optional) Hardware specifications for [analytics nodes](https://docs.atlas.mongodb.com/reference/faq/deployment/#std-label-analytics-nodes-overview) needed in the region. Analytics nodes handle analytic data such as reporting queries from BI Connector for Atlas. Analytics nodes are read-only and can never become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary). If you don't specify this parameter, no analytics nodes deploy to this region. See [below](#specs)
* `auto_scaling` - (Optional) Configuration for the Collection of settings that configures auto-scaling information for the cluster. The values for the `auto_scaling` parameter must be the same for every item in the `replication_specs` array. See [below](#auto_scaling)
* `analytics_auto_scaling` - (Optional) Configuration for the Collection of settings that configures analytics-auto-scaling information for the cluster. The values for the `analytics_auto_scaling` parameter must be the same for every item in the `replication_specs` array. See [below](#analytics_auto_scaling)
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

### electable_specs

* `instance_size` - (Required) Hardware specification for the instance sizes in this region. Each instance size has a default storage and memory capacity. The instance size you select applies to all the data-bearing hosts in your instance size.
* `disk_iops` - (Optional) Target throughput (IOPS) desired for AWS storage attached to your cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster.
* `ebs_volume_type` - (Optional) Type of storage you want to attach to your AWS-provisioned cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster. Valid values are:
    * `STANDARD` volume types can't exceed the default IOPS rate for the selected volume size.
    * `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size.
* `node_count` - (Optional) Number of nodes of the given type for MongoDB Atlas to deploy to the region.

### analytics_specs

* `disk_iops` - (Optional) Target throughput (IOPS) desired for AWS storage attached to your cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster.
* `ebs_volume_type` - (Optional) Type of storage you want to attach to your AWS-provisioned cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster. Valid values are:
    * `STANDARD` volume types can't exceed the default IOPS rate for the selected volume size.
    * `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size.
* `instance_size` - (Optional) Hardware specification for the instance sizes in this region. Each instance size has a default storage and memory capacity. The instance size you select applies to all the data-bearing hosts in your instance size.
* `node_count` - (Optional) Number of nodes of the given type for MongoDB Atlas to deploy to the region.

### read_only_specs

* `disk_iops` - (Optional) Target throughput (IOPS) desired for AWS storage attached to your cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster.
* `ebs_volume_type` - (Optional) Type of storage you want to attach to your AWS-provisioned cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster. Valid values are:
    * `STANDARD` volume types can't exceed the default IOPS rate for the selected volume size.
    * `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size.
* `instance_size` - (Optional) Hardware specification for the instance sizes in this region. Each instance size has a default storage and memory capacity. The instance size you select applies to all the data-bearing hosts in your instance size.
* `node_count` - (Optional) Number of nodes of the given type for MongoDB Atlas to deploy to the region.

### auto_scaling

* `disk_gb_enabled` - (Optional) Flag that indicates whether this cluster enables disk auto-scaling. This parameter defaults to true.
    - Set to `true` to enable disk auto-scaling.
    - Set to `false` to disable disk auto-scaling.

~> **IMPORTANT:** If `disk_gb_enabled` is true, then Atlas will automatically scale disk size up and down.
This will cause the value of `disk_size_gb` returned to potentially be different than what is specified in the Terraform config and if one then applies a plan, not noting this, Terraform will scale the cluster disk size back to the original `disk_size_gb` value.
To prevent this a lifecycle customization should be used, i.e.:  
`lifecycle {
  ignore_changes = [disk_size_gb]
}`
After adding the `lifecycle` block to explicitly change `disk_size_gb` comment out the `lifecycle` block and run `terraform apply`. Please be sure to uncomment the `lifecycle` block once done to prevent any accidental changes.

```terraform
// Example: ignore disk_size_gb and instance_size changes in a replica set
lifecycle {
  ignore_changes = [
    disk_size_gb,
    replication_specs[0].region_configs[0].electable_specs[0].instance_size,
    replication_specs[0].region_configs[1].electable_specs[0].instance_size,
    replication_specs[0].region_configs[2].electable_specs[0].instance_size,
  ]
}
```

* `compute_enabled` - (Optional) Flag that indicates whether instance size auto-scaling is enabled. This parameter defaults to false.

~> **IMPORTANT:** If `compute_enabled` is true, then Atlas will automatically scale up to the maximum provided and down to the minimum, if provided.
This will cause the value of `instance_size` returned to potentially be different than what is specified in the Terraform config and if one then applies a plan, not noting this, Terraform will scale the cluster back to the original `instance_size` value.
To prevent this a lifecycle customization should be used, i.e.:  
`lifecycle {
  ignore_changes = [instance_size]
}`
After adding the `lifecycle` block to explicitly change `instance_size` comment out the `lifecycle` block and run `terraform apply`. Please be sure to uncomment the `lifecycle` block once done to prevent any accidental changes.

* `compute_scale_down_enabled` - (Optional) Flag that indicates whether the instance size may scale down. Atlas requires this parameter if `replication_specs.#.region_configs.#.auto_scaling.0.compute_enabled` : true. If you enable this option, specify a value for `replication_specs.#.region_configs.#.auto_scaling.0.compute_min_instance_size`.
* `compute_min_instance_size` - (Optional) Minimum instance size to which your cluster can automatically scale (such as M10). Atlas requires this parameter if `replication_specs.#.region_configs.#.auto_scaling.0.compute_scale_down_enabled` is true.
* `compute_max_instance_size` - (Optional) Maximum instance size to which your cluster can automatically scale (such as M40). Atlas requires this parameter if `replication_specs.#.region_configs.#.auto_scaling.0.compute_enabled` is true.

### analytics_auto_scaling

* `disk_gb_enabled` - (Optional) Flag that indicates whether this cluster enables disk auto-scaling. This parameter defaults to true.
* `compute_enabled` - (Optional) Flag that indicates whether instance size auto-scaling is enabled. This parameter defaults to false.

~> **IMPORTANT:** If `compute_enabled` is true, then Atlas will automatically scale up to the maximum provided and down to the minimum, if provided.
This will cause the value of `instance_size` returned to potential be different than what is specified in the Terraform config and if one then applies a plan, not noting this, Terraform will scale the cluster back down to the original `instance_size` value.
To prevent this a lifecycle customization should be used, i.e.:  
`lifecycle {
  ignore_changes = [instance_size]
}`
After adding the `lifecycle` block to explicitly change `instance_size` comment out the `lifecycle` block and run `terraform apply`. Please be sure to uncomment the `lifecycle` block once done to prevent any accidental changes.

* `compute_scale_down_enabled` - (Optional) Flag that indicates whether the instance size may scale down. Atlas requires this parameter if `replication_specs.#.region_configs.#.analytics_auto_scaling.0.compute_enabled` : true. If you enable this option, specify a value for `replication_specs.#.region_configs.#.analytics_auto_scaling.0.compute_min_instance_size`.
* `compute_min_instance_size` - (Optional) Minimum instance size to which your cluster can automatically scale (such as M10). Atlas requires this parameter if `replication_specs.#.region_configs.#.analytics_auto_scaling.0.compute_scale_down_enabled` is true.
* `compute_max_instance_size` - (Optional) Maximum instance size to which your cluster can automatically scale (such as M40). Atlas requires this parameter if `replication_specs.#.region_configs.#.analytics_auto_scaling.0.compute_enabled` is true.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - The cluster ID.
*  `mongo_db_version` - Version of MongoDB the cluster runs, in `major-version`.`minor-version` format.
* `id` -	The Terraform's unique identifier used internally for state management.
* `connection_strings` - Set of connection strings that your applications use to connect to this cluster. More info in [Connection-strings](https://docs.mongodb.com/manual/reference/connection-string/). Use the parameters in this object to connect your applications to this cluster. To learn more about the formats of connection strings, see [Connection String Options](https://docs.atlas.mongodb.com/reference/faq/connection-changes/). NOTE: Atlas returns the contents of this object after the cluster is operational, not while it builds the cluster.

   **NOTE** Connection strings must be returned as a list, therefore to refer to a specific attribute value add index notation. Example: mongodbatlas_advanced_cluster.cluster-test.connection_strings.0.standard_srv

   Private connection strings are not available until the respective `mongodbatlas_privatelink_endpoint_service` resources are fully applied. Add a `depends_on = [mongodbatlas_privatelink_endpoint_service.example]` to ensure `connection_strings` are available following `terraform apply`.\
   If the expected connection string(s) do not contain a value, a `terraform refresh` may need to be performed to obtain the value. One can also view the status of the peered connection in the [Atlas UI](https://docs.atlas.mongodb.com/security-vpc-peering/).

    - `connection_strings.standard` -   Public mongodb:// connection string for this cluster.
    - `connection_strings.standard_srv` - Public mongodb+srv:// connection string for this cluster. The mongodb+srv protocol tells the driver to look up the seed list of hosts in DNS. Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don’t need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn’t  , use connectionStrings.standard.
    - `connection_strings.private` -   [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_srv` -  [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb+srv://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_endpoint` - Private endpoint connection strings. Each object describes the connection strings you can use to connect to this cluster through a private endpoint. Atlas returns this parameter only if you deployed a private endpoint to all regions to which you deployed this cluster's nodes.
    - `connection_strings.private_endpoint.#.connection_string` - Private-endpoint-aware `mongodb://`connection string for this private endpoint.
    - `connection_strings.private_endpoint.#.srv_connection_string` - Private-endpoint-aware `mongodb+srv://` connection string for this private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in DNS . Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don't need to: Append the seed list or Change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use `connection_strings.private_endpoint[n].connection_string`
    - `connection_strings.private_endpoint.#.srv_shard_optimized_connection_string` - Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for connectionStrings.privateEndpoint[n].srvConnectionString.
    - `connection_strings.private_endpoint.#.type` - Type of MongoDB process that you connect to with the connection strings. Atlas returns `MONGOD` for replica sets, or `MONGOS` for sharded clusters.
    - `connection_strings.private_endpoint.#.endpoints` - Private endpoint through which you connect to Atlas when you use `connection_strings.private_endpoint[n].connection_string` or `connection_strings.private_endpoint[n].srv_connection_string`
    - `connection_strings.private_endpoint.#.endpoints.#.endpoint_id` - Unique identifier of the private endpoint.
    - `connection_strings.private_endpoint.#.endpoints.#.provider_name` - Cloud provider to which you deployed the private endpoint. Atlas returns `AWS` or `AZURE`.
    - `connection_strings.private_endpoint.#.endpoints.#.region` - Region to which you deployed the private endpoint.
* `state_name` - Current state of the cluster. The possible states are:
    - IDLE
    - CREATING
    - UPDATING
    - DELETING
    - DELETED
    - REPAIRING
* `replication_specs` - Set of replication specifications for the cluster. Primary usage is covered under the [replication_specs argument reference](#replication_specs), though there are some computed attributes:
  - `replication_specs.#.container_id` - A key-value map of the Network Peering Container ID(s) for the configuration specified in `region_configs`. The Container ID is the id of the container created when the first cluster in the region (AWS/Azure) or project (GCP) was created.  The syntax is `"providerName:regionName" = "containerId"`. Example `AWS:US_EAST_1" = "61e0797dde08fb498ca11a71`.

## Import

Clusters can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_advanced_cluster.my_cluster 1112222b3bf99403840e8934-Cluster0
```

See detailed information for arguments and attributes: [MongoDB API Advanced Clusters](https://docs.atlas.mongodb.com/reference/api/cluster-advanced/create-one-cluster-advanced/)
