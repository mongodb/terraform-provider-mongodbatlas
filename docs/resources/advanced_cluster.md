---
subcategory: "Clusters"
---

# Resource: mongodbatlas_advanced_cluster

`mongodbatlas_advanced_cluster` provides an Advanced Cluster resource. The resource lets you create, edit and delete advanced clusters.

~> **IMPORTANT:** If upgrading from our provider versions 1.x.x to 2.0.0 or later, you will be required to update your `mongodbatlas_advanced_cluster` resource configuration. Please refer [this guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/migrate-to-advanced-cluster-2.0) for details. This new implementation uses the recommended Terraform Plugin Framework, which, in addition to providing a better user experience and other features, adds support for the `moved` block between different resource types.

~> **IMPORTANT:** We recommend all new MongoDB Atlas Terraform users start with the [`mongodbatlas_advanced_cluster`](advanced_cluster) resource.  Key differences between [`mongodbatlas_cluster`](cluster) and [`mongodbatlas_advanced_cluster`](advanced_cluster) include support for [Multi-Cloud Clusters](https://www.mongodb.com/blog/post/introducing-multicloud-clusters-on-mongodb-atlas), [Asymmetric Sharding](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema), and [Independent Scaling of Analytics Node Tiers](https://www.mongodb.com/blog/post/introducing-ability-independently-scale-atlas-analytics-node-tiers). For existing [`mongodbatlas_cluster`](cluster) resource users see our [Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

~> **IMPORTANT:** When modifying cluster configurations, you may see `(known after apply)` markers for many attributes, even those you haven't changed. This is expected behavior. See the ["known after apply" verbosity](#known-after-apply-verbosity) section below for details.

-> **NOTE:** If Backup Compliance Policy is enabled for the project for which this backup schedule is defined, you cannot modify the backup schedule for an individual cluster below the minimum requirements set in the Backup Compliance Policy.  See [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

-> **NOTE:** A network container is created for each provider/region combination on the advanced cluster. This can be referenced via a computed attribute for use with other resources. Refer to the `replication_specs[#].container_id` attribute in the [Attributes Reference](#attributes_reference) for more information.

-> **NOTE:** To enable Cluster Extended Storage Sizes use the `is_extended_storage_sizes_enabled` parameter in the [mongodbatlas_project resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project).

-> **NOTE:** The Low-CPU instance clusters are prefixed with `R`, for example `R40`. For complete list of Low-CPU instance clusters see Cluster Configuration Options under each [Cloud Provider](https://www.mongodb.com/docs/atlas/reference/cloud-providers).

-> **NOTE:** Groups and projects are synonymous terms. You might find group_id in the official documentation.

-> **NOTE:** This resource supports Flex clusters. Additionally, you can upgrade [M0 clusters to Flex](#example-tenant-cluster-upgrade-to-flex) and [Flex clusters to Dedicated](#Example-Flex-Cluster-Upgrade). When creating a Flex cluster, make sure to set the priority value to 7.

## Example Usage

### Example single provider and single region

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"
  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"
            node_count    = 3
          }
          analytics_specs = {
            instance_size = "M10"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]
}
```

### Example Tenant Cluster

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
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
```

-> **NOTE** Upgrading the tenant cluster to a Flex cluster or a dedicated cluster is supported. When upgrading to a Flex cluster, change the `provider_name` from "TENANT" to "FLEX". See [Example Tenant Cluster Upgrade to Flex](#example-tenant-cluster-upgrade-to-flex) below. When upgrading to a dedicated cluster, change the `provider_name` to your preferred provider (AWS, GCP or Azure) and remove the variable `backing_provider_name`. See the [Example Tenant Cluster Upgrade](#Example-Tenant-Cluster-Upgrade) below. You can upgrade a tenant cluster only to a single provider on an M10-tier cluster or greater.

When upgrading from the tenant, *only* the upgrade changes will be applied. This helps avoid a corrupt state file in the event that the upgrade succeeds but subsequent updates fail within the same `terraform apply`. To apply additional cluster changes, run a secondary `terraform apply` after the upgrade succeeds.


### Example Tenant Cluster Upgrade

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"
          }
          provider_name         = "AWS"
          region_name           = "US_EAST_1"
          priority              = 7
        }
      ]
    }
  ]
}
```

### Example Tenant Cluster Upgrade to Flex

```terraform
resource "mongodbatlas_advanced_cluster" "example-flex" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
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
```

### Example Flex Cluster

```terraform
resource "mongodbatlas_advanced_cluster" "example-flex" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
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
```

**NOTE**: Upgrading the Flex cluster is supported. When upgrading from a Flex cluster, change the `provider_name` from "TENANT" to your preferred provider (AWS, GCP or Azure) and remove the variable `backing_provider_name`.  See the [Example Flex Cluster Upgrade](#Example-Flex-Cluster-Upgrade) below. You can upgrade a Flex cluster only to a single provider on an M10-tier cluster or greater. 

When upgrading from a flex cluster, *only* the upgrade changes will be applied. This helps avoid a corrupt state file in the event that the upgrade succeeds but subsequent updates fail within the same `terraform apply`. To apply additional cluster changes, run a secondary `terraform apply` after the upgrade succeeds.


### Example Flex Cluster Upgrade

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"
          }
          provider_name         = "AWS"
          region_name           = "US_EAST_1"
          priority              = 7
        }
      ]
    }
  ]
}
```

### Example Multi-Cloud Cluster
```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = "PROJECT ID"
  name         = "NAME OF CLUSTER"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"
            node_count    = 3
          }
          analytics_specs = {
            instance_size = "M10"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }, 
        {
          electable_specs = {
            instance_size = "M10"
            node_count    = 2
          }
          provider_name = "GCP"
          priority      = 6
          region_name   = "NORTH_AMERICA_NORTHEAST_1"
        }
      ]
    }
  ]
}
```
### Example of a Multi Cloud Sharded Cluster with 2 shards

```terraform
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = var.cluster_name
  cluster_type = "SHARDED"
  backup_enabled = true

  replication_specs = [
    {   # shard 1
      region_configs = [
        { 
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }, 
        { 
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
        }
      ]
    }, 
    {   # shard 2
      region_configs = [
        { 
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }, 
        { 
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
        }
      ]
    }
  ]

  advanced_configuration = {
    javascript_enabled                   = true
    oplog_size_mb                        = 991
    sample_refresh_interval_bi_connector = 300
  }
}
```

### Example of a Global Cluster with 2 zones
```terraform
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "GEOSHARDED"
  backup_enabled = true

  replication_specs = [
    { # shard 1 - zone n1
      zone_name  = "zone n1"

      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }, 
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
        }
      ]
    }, 
    {  # shard 2 - zone n1
      zone_name  = "zone n1"

      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }, 
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "US_EAST_2"
        }
      ]
    }, 
    {  # shard 1 - zone n2
      zone_name  = "zone n2"

      region_configs = [
        { 
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }, 
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "EUROPE_NORTH"
        }
      ]
    }, 
    {  # shard 2 - zone n2
      zone_name  = "zone n2"

      region_configs = [
        { 
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }, {
          electable_specs ={
            instance_size = "M30"
            node_count    = 2
          }
          provider_name = "AZURE"
          priority      = 6
          region_name   = "EUROPE_NORTH"
        }
      ]
    }
  ]

  advanced_configuration = {
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
    value = mongodbatlas_advanced_cluster.cluster.connection_strings.standard
}
# Example return string: standard = "mongodb://cluster-atlas-shard-00-00.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-01.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-02.ygo1m.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-12diht-shard-0"
```
Standard srv
```terraform
output "standard_srv" {
    value = mongodbatlas_advanced_cluster.cluster.connection_strings.standard_srv
}
# Example return string: standard_srv = "mongodb+srv://cluster-atlas.ygo1m.mongodb.net"
```
Private with Network peering and Custom DNS AWS enabled
```terraform
output "private" {
    value = mongodbatlas_advanced_cluster.cluster.connection_strings.private
}
# Example return string: private = "mongodb://cluster-atlas-shard-00-00-pri.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-01-pri.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-02-pri.ygo1m.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-12diht-shard-0"
```
Private srv with Network peering and Custom DNS AWS enabled
```terraform
output "private_srv" {
    value = mongodbatlas_advanced_cluster.cluster.connection_strings.private_srv
}
# Example return string: private_srv = "mongodb+srv://cluster-atlas-pri.ygo1m.mongodb.net"
```

By endpoint_service_id
```terraform
locals {
  endpoint_service_id = google_compute_network.default.name
  private_endpoints   = coalesce(mongodbatlas_advanced_cluster.cluster.connection_strings.private_endpoint, [])
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
* [GCP Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_privatelink_endpoint/gcp)
* [Azure Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_privatelink_endpoint/azure)
* [AWS, Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster)
* [AWS, Regionalized Private Endpoints](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster-geosharded)


### Further Examples
- [Asymmetric Sharded Cluster](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_advanced_cluster/asymmetric-sharded-cluster)
- [Auto-Scaling Per Shard](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_advanced_cluster/auto-scaling-per-shard)
- [Global Cluster](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_advanced_cluster/global-cluster)
- [Multi-Cloud](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_advanced_cluster/multi-cloud)
- [Tenant Upgrade](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_advanced_cluster/tenant-upgrade)
- [Version Upgrade with Pinned FCV](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/mongodbatlas_advanced_cluster/version-upgrade-with-pinned-fcv)
- [Migrate Cluster to Advanced Cluster](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_cluster_to_advanced_cluster/basic)

## Argument Reference

* `project_id` - (Required) Unique ID for the project to create the cluster.
* `name` - (Required) Name of the cluster as it appears in Atlas. Once the cluster is created, its name cannot be changed. **WARNING** Changing the name will result in destruction of the existing cluster and the creation of a new cluster.

* `backup_enabled` - (Optional) Flag that indicates whether the cluster can perform backups.
  If `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters.

  Backup uses:
  [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/#std-label-backup-cloud-provider) for dedicated clusters.
  [Flex Cluster Backups](https://www.mongodb.com/docs/atlas/backup/cloud-backup/flex-cluster-backup/) for flex clusters.
  If "`backup_enabled`"  is `false` (default), the cluster doesn't use Atlas backups.

- `retain_backups_enabled` - (Optional) Set to true to retain backup snapshots for the deleted cluster. This parameter applies to the Delete operation and only affects M10 and above clusters. If you encounter the `CANNOT_DELETE_SNAPSHOT_WITH_BACKUP_COMPLIANCE_POLICY` error code, see [how to delete a cluster with Backup Compliance Policy](../guides/delete-cluster-with-backup-compliance-policy.md).

-> **NOTE** Prior version of provider had parameter as `bi_connector` state will migrate it to new value you only need to update parameter in your terraform file

* `bi_connector_config` - (Optional) Configuration settings applied to BI Connector for Atlas on this cluster. The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger clusters. The BI Connector is a powerful tool which provides users SQL-based access to their MongoDB databases. As a result, the BI Connector performs operations which may be CPU and memory intensive. Given the limited hardware resources on M10 and M20 cluster tiers, you may experience performance degradation of the cluster when enabling the BI Connector. If this occurs, upgrade to an M30 or larger cluster or disable the BI Connector. See [below](#bi_connector_config).
* `cluster_type` - (Required)Type of the cluster that you want to create.
    Accepted values include:
      - `REPLICASET` Replica set
      - `SHARDED`	Sharded cluster
      - `GEOSHARDED` Global Cluster

* `encryption_at_rest_provider` - (Optional) Possible values are AWS, GCP, AZURE or NONE.  Only needed if you desire to manage the keys, see [Encryption at Rest using Customer Key Management](https://docs.atlas.mongodb.com/security-kms-encryption/) for complete documentation.  You must configure encryption at rest for the Atlas project before enabling it on any cluster in the project. For Documentation, see [AWS](https://docs.atlas.mongodb.com/security-aws-kms/), [GCP](https://docs.atlas.mongodb.com/security-kms-encryption/) and [Azure](https://docs.atlas.mongodb.com/security-azure-kms/#std-label-security-azure-kms). Requirements are if `replication_specs[#].region_configs[#].<type>Specs.instance_size` is M10 or greater and `backup_enabled` is false or omitted.   
* `tags` - (Optional) Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#tags).
* `labels` - (Optional) Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#labels). **DEPRECATED** Use `tags` instead.
* `mongo_db_major_version` - (Optional) Version of the cluster to deploy. Atlas supports all the MongoDB versions that have **not** reached [End of Live](https://www.mongodb.com/legal/support-policy/lifecycles) for M10+ clusters. If omitted, Atlas deploys the cluster with the default version. For more details, see [documentation](https://www.mongodb.com/docs/atlas/reference/faq/database/#which-versions-of-mongodb-do-service-clusters-use-). Atlas always deploys the cluster with the latest stable release of the specified version.  If you set a value to this parameter and set `version_release_system` `CONTINUOUS`, the resource returns an error. Either clear this parameter or set `version_release_system`: `LTS`.
* `pinned_fcv` - (Optional) Pins the Feature Compatibility Version (FCV) to the current MongoDB version with a provided expiration date. To unpin the FCV the `pinned_fcv` attribute must be removed. This operation can take several minutes as the request processes through the MongoDB data plane. Once FCV is unpinned it will not be possible to downgrade the `mongo_db_major_version`. It is advised that updates to `pinned_fcv` are done isolated from other cluster changes. If a plan contains multiple changes, the FCV change will be applied first. If FCV is unpinned past the expiration date the `pinned_fcv` attribute must be removed. The following [knowledge hub article](https://kb.corp.mongodb.com/article/000021785/) and [FCV documentation](https://www.mongodb.com/docs/atlas/tutorial/major-version-change/#manage-feature-compatibility--fcv--during-upgrades) can be referenced for more details. See [below](#pinned_fcv).
* `pit_enabled` - (Optional) Flag that indicates if the cluster uses Continuous Cloud Backup.
* `replication_specs` - List of settings that configure your cluster regions. This attribute has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations. The `replication_specs` configuration for all shards within the same zone must be the same, with the exception of `instance_size` and `disk_iops` that can scale independently. Note that independent `disk_iops` values are only supported for AWS provisioned IOPS, or Azure regions that support Extended IOPS. See [below](#replication_specs).
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
* `timeouts`- (Optional) The duration of time to wait for Cluster to be created, updated, or deleted. The timeout value is defined by a signed sequence of decimal numbers with a time unit suffix such as: `1h45m`, `300s`, `10m`, etc. The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Advanced Cluster create & delete is `3h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).
* `accept_data_risks_and_force_replica_set_reconfig` - (Optional) If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set `accept_data_risks_and_force_replica_set_reconfig` to the current date. Learn more about Reconfiguring a Replica Set during a regional outage [here](https://dochub.mongodb.org/core/regional-outage-reconfigure-replica-set).
* `global_cluster_self_managed_sharding` - (Optional) Flag that indicates if cluster uses Atlas-Managed Sharding (false, default) or Self-Managed Sharding (true). It can only be enabled for Global Clusters (`GEOSHARDED`). It cannot be changed once the cluster is created. Use this mode if you're an advanced user and the default configuration is too restrictive for your workload. If you select this option, you must manually configure the sharding strategy, more information [here](https://www.mongodb.com/docs/atlas/tutorial/create-global-cluster/#select-your-sharding-configuration).
* `replica_set_scaling_strategy` - (Optional) Replica set scaling mode for your cluster. Valid values are `WORKLOAD_TYPE`, `SEQUENTIAL` and `NODE_TYPE`. By default, Atlas scales under `WORKLOAD_TYPE`. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes. When configured as `SEQUENTIAL`, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads. When configured as `NODE_TYPE`, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads. [Modify the Replica Set Scaling Mode](https://dochub.mongodb.org/core/scale-nodes)
* `redact_client_log_data` - (Optional) Flag that enables or disables log redaction, see the [manual](https://www.mongodb.com/docs/manual/administration/monitoring/#log-redaction) for more information. Use this in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements. **Note**: Changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.
* `config_server_management_mode` - (Optional) Config Server Management Mode for creating or updating a sharded cluster. Valid values are `ATLAS_MANAGED` (default) and `FIXED_TO_DEDICATED`. When configured as `ATLAS_MANAGED`, Atlas may automatically switch the cluster's config server type for optimal performance and savings. When configured as `FIXED_TO_DEDICATED`, the cluster will always use a dedicated config server. To learn more, see the [Sharded Cluster Config Servers documentation](https://dochub.mongodb.org/docs/manual/core/sharded-cluster-config-servers/).
- `delete_on_create_timeout`- (Optional) Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.

### bi_connector_config

Specifies BI Connector for Atlas configuration.

```terraform
bi_connector_config = {
  enabled         = true
  read_preference = "secondary"
}
```

* `enabled` - (Optional) Specifies whether or not BI Connector for Atlas is enabled on the cluster.
*
  - Set to `true` to enable BI Connector for Atlas.
  - Set to `false` to disable BI Connector for Atlas.

* `read_preference` - (Optional) Specifies the read preference to be used by BI Connector for Atlas on the cluster. Each BI Connector for Atlas read preference contains a distinct combination of [readPreference](https://docs.mongodb.com/manual/core/read-preference/) and [readPreferenceTags](https://docs.mongodb.com/manual/core/read-preference/#tag-sets) options. For details on BI Connector for Atlas read preferences, refer to the [BI Connector Read Preferences Table](https://docs.atlas.mongodb.com/tutorial/create-global-writes-cluster/#bic-read-preferences).

  - Set to "primary" to have BI Connector for Atlas read from the primary.

  - Set to "secondary" to have BI Connector for Atlas read from a secondary member. Default if there are no analytics nodes in the cluster.

  - Set to "analytics" to have BI Connector for Atlas read from an analytics node. Default if the cluster contains analytics nodes.

### Advanced Configuration Options

-> **NOTE:** Prior to setting these options please ensure you read https://docs.atlas.mongodb.com/cluster-config/additional-options/.

-> **NOTE:** Once you set some `advanced_configuration` attributes, we recommended to explicitly set those attributes to their intended value instead of removing them from the configuration. For example, if you set `javascript_enabled` to `false`, and later you want to go back to the default value (true), you must set it back to `true` instead of removing it.

Include **desired options** within advanced_configuration:

```terraform
// Nest options within advanced_configuration
 advanced_configuration = {
   javascript_enabled                   = false
   minimum_enabled_tls_protocol         = "TLS1_2"
 }
```

* `default_write_concern` - (Optional) [Default level of acknowledgment requested from MongoDB for write operations](https://docs.mongodb.com/manual/reference/write-concern/) set for this cluster. MongoDB 6.0 clusters default to [majority](https://docs.mongodb.com/manual/reference/write-concern/).
* `javascript_enabled` - (Optional) When true (default), the cluster allows execution of operations that perform server-side executions of JavaScript. When false, the cluster disables execution of those operations.
* `minimum_enabled_tls_protocol` - (Optional) Sets the minimum Transport Layer Security (TLS) version the cluster accepts for incoming connections. Valid values are:
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
* `change_stream_options_pre_and_post_images_expire_after_seconds` - (Optional) The minimum pre- and post-image retention time in seconds. This option corresponds to the `changeStreamOptions.preAndPostImages.expireAfterSeconds` cluster parameter. Defaults to `-1`(off). This setting controls the retention policy of change stream pre- and post-images. Pre- and post-images are the versions of a document before and after document modification, respectively. `expireAfterSeconds` controls how long MongoDB retains pre- and post-images. When set to -1 (off), MongoDB uses the default retention policy: pre- and post-images are retained until the corresponding change stream events are removed from the oplog. To set the minimum pre- and post-image retention time, specify an integer value greater than zero. Setting this too low could increase the risk of interrupting Realm sync or triggers processing. This parameter is only supported for MongoDB version 6.0 and above.
* `default_max_time_ms` - (Optional) Default time limit in milliseconds for individual read operations to complete. This option corresponds to the [defaultMaxTimeMS](https://www.mongodb.com/docs/upcoming/reference/cluster-parameters/defaultMaxTimeMS/) cluster parameter. This parameter is supported only for MongoDB version 8.0 and above.
* `tls_cipher_config_mode` - (Optional) The TLS cipher suite configuration mode. Valid values include `CUSTOM` or `DEFAULT`. The `DEFAULT` mode uses the default cipher suites. The `CUSTOM` mode allows you to specify custom cipher suites for both TLS 1.2 and TLS 1.3. To unset, this should be set back to `DEFAULT`.
* `custom_openssl_cipher_config_tls12` - (Optional) The custom OpenSSL cipher suite list for TLS 1.2. This field is only valid when `tls_cipher_config_mode` is set to `CUSTOM`.


### tags

 ```terraform
 tags = {
    "Key 1" = "Value 1"
    "Key 2" = "Value 2"
    Key3    = "Value 3"
  }
```

Key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.

* `key` - (Required) Constant that defines the set of the tag.
* `value` - (Required) Variable that belongs to the set of the tag.

To learn more, see [Resource Tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas).

### labels

 ```terraform
 labels = {
    "Key 1" = "Value 1"
    "Key 2" = "Value 2"
    Key3    = "Value 3"
  }
```

Key-value pairs that categorize the cluster. Each key and value has a maximum length of 255 characters.  You cannot set the key `Infrastructure Tool`, it is used for internal purposes to track aggregate usage.

* `key` - The key that you want to write.
* `value` - The value that you want to write.

-> **NOTE:** MongoDB Atlas doesn't display your labels.


### replication_specs

~> **NOTE:**  We recommend reviewing our [Best Practices](#remove-or-disable-functionality) before disabling or removing any elements of replication_specs.

```terraform
//Example Multicloud
replication_specs = [
  {
    region_configs = [
      {
        electable_specs = {
          instance_size = "M10"
          node_count    = 3
        }
        analytics_specs = {
          instance_size = "M10"
          node_count    = 1
        }
        provider_name = "AWS"
        priority      = 7
        region_name   = "US_EAST_1"
      }, 
      {
        electable_specs = {
          instance_size = "M10"
          node_count    = 2
        }
        provider_name = "GCP"
        priority      = 6
        region_name   = "NORTH_AMERICA_NORTHEAST_1"
      }
    ]
  }
]
```

* `external_id` - Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. This value corresponds to Shard ID displayed in the UI.
* `region_configs` - (Optional) Configuration for the hardware specifications for nodes set for a given region. Each `region_configs` object describes the region's priority in elections and the number and type of MongoDB nodes that Atlas deploys to the region. Each `region_configs` object must have either an `analytics_specs` object, `electable_specs` object, or `read_only_specs` object. See [below](#region_configs).
* `zone_name` - (Optional) Name for the zone in a Global Cluster.
* `zone_id` - Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. If clusterType is GEOSHARDED, this value indicates the zone that the given shard belongs to and can be used to configure Global Cluster backup policies.


### region_configs

~> **NOTE:**  We recommend reviewing our [Best Practices](#remove-or-disable-functionality) before disabling or removing any elements of region_configs.

* `analytics_specs` - (Optional) Hardware specifications for [analytics nodes](https://docs.atlas.mongodb.com/reference/faq/deployment/#std-label-analytics-nodes-overview) needed in the region. Analytics nodes handle analytic data such as reporting queries from BI Connector for Atlas. Analytics nodes are read-only and can never become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary). If you don't specify this parameter, no analytics nodes deploy to this region. See [below](#specs).
* `auto_scaling` - (Optional) Configuration for the collection of settings that configures auto-scaling information for the cluster. The values for the `auto_scaling` attribute must be the same for all `region_configs` of a cluster. See [below](#auto_scaling).
* `analytics_auto_scaling` - (Optional) Configuration for the Collection of settings that configures analytics-auto-scaling information for the cluster. The values for the `analytics_auto_scaling` attribute must be the same for all `region_configs` of a cluster. See [below](#analytics_auto_scaling).
* `backing_provider_name` - (Optional) Cloud service provider on which you provision the host for a multi-tenant cluster. Use this only when a `provider_name` is `TENANT` and `instance_size` is `M0`.
* `electable_specs` - (Optional) Hardware specifications for electable nodes in the region. All `electable_specs` in the `region_configs` of a `replication_specs` must have the same `instance_size`. Electable nodes can become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary) and can enable local reads. If you do not specify this option, no electable nodes are deployed to the region. See [below](#specs).
* `priority` - (Optional)  Election priority of the region. For regions with only read-only nodes, set this value to 0.
  * If you have multiple `region_configs` objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is 7.
  * If your region has set `region_configs[#].electable_specs.node_count` to 1 or higher, it must have a priority of exactly one (1) less than another region in the `replication_specs[#].region_configs[#]` array. The highest-priority region must have a priority of 7. The lowest possible priority is 1.
* `provider_name` - (Optional) Cloud service provider on which the servers are provisioned.
  The possible values are:
  - `AWS` - Amazon AWS
  - `GCP` - Google Cloud Platform
  - `AZURE` - Microsoft Azure
  - `TENANT` - M0 multi-tenant cluster. Use `replication_specs.[#].region_configs[#].backing_provider_name` to set the cloud service provider.
* `read_only_specs` - (Optional) Hardware specifications for read-only nodes in the region. All `read_only_specs` in the `region_configs` of a `replication_specs` must have the same `instance_size` as `electable_specs`. Read-only nodes can become the [primary](https://docs.atlas.mongodb.com/reference/glossary/#std-term-primary) and can enable local reads. If you don't specify this parameter, no read-only nodes are deployed to the region. See [below](#specs).
* `region_name` - (Optional) Physical location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases.  Requires the **Atlas region name**, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).

### electable_specs

* `instance_size` - (Required) Hardware specification for the instance sizes in this region. Each instance size has a default storage and memory capacity. The instance size you select applies to all the data-bearing hosts in your instance size. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
* `disk_iops` - (Optional) Target IOPS (Input/Output Operations Per Second) desired for storage attached to this hardware. Define this attribute only if you selected AWS as your cloud service provider, `instance_size` is set to "M30" or greater (not including "Mxx_NVME" tiers), and `ebs_volume_type` is "PROVISIONED". You can't set this attribute for a multi-cloud cluster.
* `ebs_volume_type` - (Optional) Type of storage you want to attach to your AWS-provisioned cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster. Valid values are:
    * `STANDARD` volume types can't exceed the default IOPS rate for the selected volume size.
    * `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size.
* `node_count` - (Optional) Number of nodes of the given type for MongoDB Atlas to deploy to the region.
* `disk_size_gb` - (Optional) Storage capacity that the host's root volume possesses expressed in gigabytes. This value must be equal for all shards and node types. If disk size specified is below the minimum (10 GB), this parameter defaults to the minimum disk size value. Storage charge calculations depend on whether you choose the default value or a custom value.  The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier. **Note:** Using `disk_size_gb` with Standard IOPS could lead to errors and configuration issues. Therefore, it should be used only with the [Provisioned IOPS volume type](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#PROVISIONED). When using Provisioned IOPS, the disk_size_gb parameter specifies the storage capacity, but the IOPS are set independently. Ensuring that `disk_size_gb` is used exclusively with Provisioned IOPS will help avoid these issues.


### analytics_specs

* `instance_size` - (Required) Hardware specification for the instance sizes in this region. Each instance size has a default storage and memory capacity. The instance size you select applies to all the data-bearing hosts in your instance size. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
* `disk_iops` - (Optional) Target IOPS (Input/Output Operations Per Second) desired for storage attached to this hardware. Define this attribute only if you selected AWS as your cloud service provider, `instance_size` is set to "M30" or greater (not including "Mxx_NVME" tiers), and `ebs_volume_type` is "PROVISIONED". You can't set this attribute for a multi-cloud cluster.
* `ebs_volume_type` - (Optional) Type of storage you want to attach to your AWS-provisioned cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster. Valid values are:
    * `STANDARD` volume types can't exceed the default IOPS rate for the selected volume size.
    * `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size.
* `node_count` - (Optional) Number of nodes of the given type for MongoDB Atlas to deploy to the region.
* `disk_size_gb` - (Optional) Storage capacity that the host's root volume possesses expressed in gigabytes. This value must be equal for all shards and node types. If disk size specified is below the minimum (10 GB), this parameter defaults to the minimum disk size value. Storage charge calculations depend on whether you choose the default value or a custom value.  The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier. **Note:** Using `disk_size_gb` with Standard IOPS could lead to errors and configuration issues. Therefore, it should be used only with the [Provisioned IOPS volume type](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#PROVISIONED). When using Provisioned IOPS, the disk_size_gb parameter specifies the storage capacity, but the IOPS are set independently. Ensuring that `disk_size_gb` is used exclusively with Provisioned IOPS will help avoid these issues.

### read_only_specs

* `instance_size` - (Required) Hardware specification for the instance sizes in this region. Each instance size has a default storage and memory capacity. The instance size you select applies to all the data-bearing hosts in your instance size. Electable nodes and read-only nodes (known as "base nodes") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards.
* `disk_iops` - (Optional) Target IOPS (Input/Output Operations Per Second) desired for storage attached to this hardware. Define this attribute only if you selected AWS as your cloud service provider, `instance_size` is set to "M30" or greater (not including "Mxx_NVME" tiers), and `ebs_volume_type` is "PROVISIONED". You can't set this attribute for a multi-cloud cluster. This parameter defaults to the cluster tier's standard IOPS value.
* `ebs_volume_type` - (Optional) Type of storage you want to attach to your AWS-provisioned cluster. Set only if you selected AWS as your cloud service provider. You can't set this parameter for a multi-cloud cluster. Valid values are:
    * `STANDARD` volume types can't exceed the default IOPS rate for the selected volume size.
    * `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size.
* `node_count` - (Optional) Number of nodes of the given type for MongoDB Atlas to deploy to the region.
* `disk_size_gb` - (Optional) Storage capacity that the host's root volume possesses expressed in gigabytes. This value must be equal for all shards and node types. If disk size specified is below the minimum (10 GB), this parameter defaults to the minimum disk size value. Storage charge calculations depend on whether you choose the default value or a custom value.  The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier. **Note:** Using `disk_size_gb` with Standard IOPS could lead to errors and configuration issues. Therefore, it should be used only with the [Provisioned IOPS volume type](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#PROVISIONED). When using Provisioned IOPS, the disk_size_gb parameter specifies the storage capacity, but the IOPS are set independently. Ensuring that `disk_size_gb` is used exclusively with Provisioned IOPS will help avoid these issues.

### auto_scaling

* `disk_gb_enabled` - (Optional) Flag that indicates whether this cluster enables disk auto-scaling. This parameter defaults to false.

* `compute_enabled` - (Optional) Flag that indicates whether instance size auto-scaling is enabled. This parameter defaults to false. If a sharded cluster is making use of the [New Sharding Configuration](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema), auto-scaling of the instance size will be independent for each individual shard. Please reference the [Use Auto-Scaling Per Shard](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema#use-auto-scaling-per-shard) section for more details.

~> **IMPORTANT:** If `disk_gb_enabled` or `compute_enabled` is true, Atlas automatically scales the cluster up or down.
This will cause the value of `replication_specs[#].region_config[#].(electable_specs|read_only_specs).disk_size_gb` or `replication_specs[#].region_config[#].(electable_specs|read_only_specs).instance_size` returned to potentially be different than what is specified in the Terraform config. If you then apply a plan, not noting this, Terraform will scale the cluster back to the original values in the config.
To prevent unintended changes when enabling autoscaling, use a lifecycle ignore customization as shown in the example below. To explicitly change `disk_size_gb` or `instance_size` values, comment out the `lifecycle` block and run `terraform apply`. Please be sure to uncomment the `lifecycle` block once done to prevent any accidental changes.

```terraform
// Example: ignore disk_size_gb and instance_size changes in a replica set
lifecycle {
  ignore_changes = [
    replication_specs[0].region_configs[0].electable_specs.disk_size_gb,
    replication_specs[0].region_configs[0].electable_specs.instance_size,
    replication_specs[0].region_configs[0].electable_specs.disk_iops // instance_size change can affect disk_iops in case that you are using it
  ]
}
```

* `compute_scale_down_enabled` - (Optional) Flag that indicates whether the instance size may scale down. Atlas requires this parameter if `replication_specs[#].region_configs[#].auto_scaling.compute_enabled` : true. If you enable this option, specify a value for `replication_specs[#].region_configs[#].auto_scaling.compute_min_instance_size`.
* `compute_min_instance_size` - (Optional) Minimum instance size to which your cluster can automatically scale (such as M10). Atlas requires this parameter if `replication_specs[#].region_configs[#].auto_scaling.compute_scale_down_enabled` is true.
* `compute_max_instance_size` - (Optional) Maximum instance size to which your cluster can automatically scale (such as M40). Atlas requires this parameter if `replication_specs[#].region_configs[#].auto_scaling.compute_enabled` is true.

### analytics_auto_scaling

* `disk_gb_enabled` - (Optional) Flag that indicates whether this cluster enables disk auto-scaling. This parameter defaults to false.
* `compute_enabled` - (Optional) Flag that indicates whether analytics instance size auto-scaling is enabled. This parameter defaults to false. If a sharded cluster is making use of the [New Sharding Configuration](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema), auto-scaling of analytics instance size will be independent for each individual shard. Please reference the [Use Auto-Scaling Per Shard](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema#use-auto-scaling-per-shard) section for more details.

~> **IMPORTANT:** If `disk_gb_enabled` or `compute_enabled` is true, Atlas automatically scales the cluster up or down.
This will cause the value of `replication_specs[#].region_config[#].analytics_specs.disk_size_gb` or `replication_specs[#].region_config[#].analytics_specs.instance_size` returned to potentially be different than what is specified in the Terraform config. If you then apply a plan, not noting this, Terraform will scale the cluster back to the original values in the config.
To prevent unintended changes when enabling autoscaling, use a lifecycle ignore customization as shown in the example below. To explicitly change `disk_size_gb` or `instance_size` values, comment out the `lifecycle` block and run `terraform apply`. Please be sure to uncomment the `lifecycle` block once done to prevent any accidental changes.

```terraform
// Example: ignore disk_size_gb and instance_size changes in a replica set
lifecycle {
  ignore_changes = [
    replication_specs[0].region_configs[0].analytics_specs.disk_size_gb,
    replication_specs[0].region_configs[0].analytics_specs.instance_size,
    replication_specs[0].region_configs[0].analytics_specs.disk_iops // instance_size change can affect disk_iops in case that you are using it
  ]
}
```

* `compute_scale_down_enabled` - (Optional) Flag that indicates whether the instance size may scale down. Atlas requires this parameter if `replication_specs[#].region_configs[#].analytics_auto_scaling.compute_enabled` : true. If you enable this option, specify a value for `replication_specs[#].region_configs[#].analytics_auto_scaling.compute_min_instance_size`.
* `compute_min_instance_size` - (Optional) Minimum instance size to which your cluster can automatically scale (such as M10). Atlas requires this parameter if `replication_specs[#].region_configs[#].analytics_auto_scaling.compute_scale_down_enabled` is true.
* `compute_max_instance_size` - (Optional) Maximum instance size to which your cluster can automatically scale (such as M40). Atlas requires this parameter if `replication_specs[#].region_configs[#].analytics_auto_scaling.compute_enabled` is true.

### pinned_fcv

* `expiration_date` - (Required) Expiration date of the fixed FCV. This value is in the ISO 8601 timestamp format (e.g. "2024-12-04T16:25:00Z"). Note that this field cannot exceed 4 weeks from the pinned date.
* `version` - Feature compatibility version of the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - The cluster ID.
* `mongo_db_version` - Version of MongoDB the cluster runs, in `major-version`.`minor-version` format.
* `connection_strings` - Set of connection strings that your applications use to connect to this cluster. More information in [Connection-strings](https://docs.mongodb.com/manual/reference/connection-string/). Use the parameters in this object to connect your applications to this cluster. To learn more about the formats of connection strings, see [Connection String Options](https://docs.atlas.mongodb.com/reference/faq/connection-changes/). NOTE: Atlas returns the contents of this object after the cluster is operational, not while it builds the cluster.

   **NOTE** Connection strings must be returned as a list, therefore to refer to a specific attribute value add index notation. Example: mongodbatlas_advanced_cluster.cluster-test.connection_strings.0.standard_srv

   Private connection strings are not available until the respective `mongodbatlas_privatelink_endpoint_service` resources are fully applied. Add a `depends_on = [mongodbatlas_privatelink_endpoint_service.example]` to ensure `connection_strings` are available following `terraform apply`. If the expected connection string(s) do not contain a value, a `terraform refresh` may need to be performed to obtain the value. One can also view the status of the peered connection in the [Atlas UI](https://docs.atlas.mongodb.com/security-vpc-peering/).

    - `connection_strings.standard` -   Public mongodb:// connection string for this cluster.
    - `connection_strings.standard_srv` - Public mongodb+srv:// connection string for this cluster. The mongodb+srv protocol tells the driver to look up the seed list of hosts in DNS. Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don’t need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn’t  , use connectionStrings.standard.
    - `connection_strings.private` -   [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_srv` -  [Network-peering-endpoint-aware](https://docs.atlas.mongodb.com/security-vpc-peering/#vpc-peering) mongodb+srv://connection strings for each interface VPC endpoint you configured to connect to this cluster. Returned only if you created a network peering connection to this cluster.
    - `connection_strings.private_endpoint` - Private endpoint connection strings. Each object describes the connection strings you can use to connect to this cluster through a private endpoint. Atlas returns this parameter only if you deployed a private endpoint to all regions to which you deployed this cluster's nodes.
    - `connection_strings.private_endpoint[#].connection_string` - Private-endpoint-aware `mongodb://`connection string for this private endpoint.
    - `connection_strings.private_endpoint[#].srv_connection_string` - Private-endpoint-aware `mongodb+srv://` connection string for this private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in DNS . Atlas synchronizes this list with the nodes in a cluster. If the connection string uses this URI format, you don't need to: Append the seed list or Change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use `connection_strings.private_endpoint[#].connection_string`
    - `connection_strings.private_endpoint[#].srv_shard_optimized_connection_string` - Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster support it. If it doesn't, use and consult the documentation for `connection_strings.private_endpoint[#].srv_connection_string`.
    - `connection_strings.private_endpoint[#].type` - Type of MongoDB process that you connect to with the connection strings. Atlas returns `MONGOD` for replica sets, or `MONGOS` for sharded clusters.
    - `connection_strings.private_endpoint[#].endpoints` - Private endpoint through which you connect to Atlas when you use `connection_strings.private_endpoint[#].connection_string` or `connection_strings.private_endpoint[#].srv_connection_string`
    - `connection_strings.private_endpoint[#].endpoints[#].endpoint_id` - Unique identifier of the private endpoint.
    - `connection_strings.private_endpoint[#].endpoints[#].provider_name` - Cloud provider to which you deployed the private endpoint. Atlas returns `AWS` or `AZURE`.
    - `connection_strings.private_endpoint[#].endpoints[#].region` - Region to which you deployed the private endpoint.
* `state_name` - Current state of the cluster. The possible states are:
    - IDLE
    - CREATING
    - UPDATING
    - DELETING
    - DELETED
    - REPAIRING
* `replication_specs[#].container_id` - A key-value map of the Network Peering Container ID(s) for the configuration specified in `region_configs`. The Container ID is the id of the container created when the first cluster in the region (AWS/Azure) or project (GCP) was created.  The syntax is `"providerName:regionName" = "containerId"`. Example `AWS:US_EAST_1" = "61e0797dde08fb498ca11a71`.
* `config_server_type` Describes a sharded cluster's config server type. Valid values are `DEDICATED` and `EMBEDDED`. To learn more, see the [Sharded Cluster Config Servers documentation](https://dochub.mongodb.org/docs/manual/core/sharded-cluster-config-servers/).
* `pinned_fcv.version` - Feature compatibility version of the cluster.


## Import

Clusters can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_advanced_cluster.my_cluster 1112222b3bf99403840e8934-Cluster0
```

See detailed information for arguments and attributes: [MongoDB API Advanced Clusters](https://docs.atlas.mongodb.com/reference/api/cluster-advanced/create-one-cluster-advanced/)

~> **IMPORTANT:**
<br> &#8226; When a cluster is imported, the resulting schema structure will always return the new schema including `replication_specs` per independent shards of the cluster.

## Move

`mongodbatlas__cluster` resources can be moved to `mongodbatlas_advanced_cluster` in Terraform v1.8 and later, e.g.: 

```terraform
moved {
  from = mongodbatlas_cluster.cluster
  to   = mongodbatlas_advanced_cluster.cluster
}
```

More information about moving resources can be found in our [Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide) and in the Terraform documentation [here](https://developer.hashicorp.com/terraform/language/moved) and [here](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring).

## Considerations and Best Practices

### "known after apply" verbosity

When modifying cluster configurations, you may see `(known after apply)` markers in your Terraform plan output, even for attributes you haven't modified. This is expected behavior.

▸Understanding the Behavior

The provider v2.x uses the Terraform [Plugin Framework (TPF)](https://developer.hashicorp.com/terraform/plugin/framework), which is more strict and verbose with computed values than the legacy [SDKv2 framework](https://developer.hashicorp.com/terraform/plugin/sdkv2) used in v1.x. Key points:

- "(known after apply)" doesn't mean the value will change - It indicates a computed value that can't be known in advance, even if the value remains the same.
- Optional/Computed attributes show as "known after apply" when not explicitly set, but won't actually change.
- Actual changes are marked with an arrow (`->`) in the plan - These values will truly change.
- Dependent attributes may change - Some changes can affect related attributes (e.g., change to `zone_name` may update `zone_id`, `region_name` may update `container_id`, `instance_size` may update `disk_iops`, or `provider_name` may update `ebs_volume_type`).

▸Mitigating Plan Verbosity

To reduce the number of `(known after apply)` entries in your plan output:

1. Explicitly declare known values in your configuration where possible:
   ```terraform
   replication_specs = [
     {
       region_configs = [
         {
           electable_specs = {
             instance_size   = "M30"
             node_count      = 3
             disk_size_gb    = 100  # Explicitly set even if it's the default
             disk_iops       = 3000 # Explicitly set if known
             ebs_volume_type = "STANDARD" # Explicitly set the volume type
           }
           # ... other configuration
         }
       ]
     }
   ]
   ```

2. Use lifecycle ignore_changes for attributes that frequently show as unknown but don't affect your infrastructure, for example:
   ```terraform
   lifecycle {
     ignore_changes = [
       state_name,
       replication_specs[0].container_id,
       replication_specs[0].external_id,
       replication_specs[0].zone_id
     ]
   }
   ```

3. Review the plan carefully to distinguish between:
   - Actual changes: Attributes you're intentionally modifying.
   - Computed updates: Attributes marked as `(known after apply)` that will be recalculated but won't cause operational changes.

▸Important Notes

- `(known after apply)` markers don't represent actual changes—only values with an arrow (`->`) will change, along with any dependent attributes affected by those changes.
- Plans with `(known after apply)` entries are safe to apply.
- The MongoDB team is working to reduce plan verbosity, though no timeline is available yet.

### Remove or disable functionality

To disable or remove functionalities, we recommended to explicitly set those attributes to their intended value instead of removing them from the configuration. This will ensure no ambiguity in what the final terraform resource state will be. For example, if you have a `read_only_specs` block in your cluster definition like this one:
```terraform
...
region_configs = [
  {
    read_only_specs =  {
      instance_size = "M10"
      node_count    = 1
    }
    electable_specs = {
      instance_size = "M10"
      node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "US_WEST_1"
  }
]
...
```
and your intention is to delete the read-only nodes, you should set the `node_count` attribute to `0` instead of removing the block:
```terraform
...
region_configs = [
  {
    read_only_specs =  {
      instance_size = "M10"
      node_count    = 0
    }
    electable_specs = {
      instance_size = "M10"
      node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "US_WEST_1"
  }
]
...
```
Similarly, if you have compute and disk auto-scaling enabled:
```terraform
...
auto_scaling = {
  disk_gb_enabled = true
  compute_enabled = true
  compute_scale_down_enabled = true
  compute_min_instance_size = "M30"
  compute_max_instance_size = "M50"
}
...
``` 
and you want to disable them, you should set the `disk_gb_enabled` and `compute_enabled` attributes to `false` instead of removing the block:
```terraform
...
auto_scaling = {
  disk_gb_enabled = false
  compute_enabled = false
  compute_scale_down_enabled = false
}
...
```
