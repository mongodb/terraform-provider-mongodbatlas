---
subcategory: "Clusters"
---

# Resource: mongodbatlas_advanced_cluster

`mongodbatlas_advanced_cluster` provides an Advanced Cluster resource. The resource lets you create, edit and delete advanced clusters.

~> **IMPORTANT:** If upgrading from our provider versions 1.x.x to 2.0.0 or later, you will be required to update your `mongodbatlas_advanced_cluster` resource configuration. Please refer [this guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/migrate-to-advanced-cluster-2.0) for details. This new implementation uses the recommended Terraform Plugin Framework, which, in addition to providing a better user experience and other features, adds support for the `moved` block between different resource types.

~> **IMPORTANT:** We recommend all new MongoDB Atlas Terraform users start with the [`mongodbatlas_advanced_cluster`](advanced_cluster) resource.  Key differences between [`mongodbatlas_cluster`](cluster) and [`mongodbatlas_advanced_cluster`](advanced_cluster) include support for [Multi-Cloud Clusters](https://www.mongodb.com/blog/post/introducing-multicloud-clusters-on-mongodb-atlas), [Asymmetric Sharding](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema), and [Independent Scaling of Analytics Node Tiers](https://www.mongodb.com/blog/post/introducing-ability-independently-scale-atlas-analytics-node-tiers). For existing [`mongodbatlas_cluster`](cluster) resource users see our [Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

~> **IMPORTANT:** When modifying cluster configurations, you may see `(known after apply)` markers for many attributes, even those you haven't changed. This is expected behavior. See the ["known after apply" verbosity](#known-after-apply-verbosity) section below for details.

~> **IMPORTANT:** When configuring auto-scaling, you can now use `use_effective_fields` to simplify your Terraform workflow. See the [Auto-Scaling with Effective Fields](#auto-scaling-with-effective-fields) section below for details.

-> **NOTE:** If Backup Compliance Policy is enabled for the project for which this backup schedule is defined, you cannot modify the backup schedule for an individual cluster below the minimum requirements set in the Backup Compliance Policy.  See [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

-> **NOTE:** A network container is created for each provider/region combination on the advanced cluster. This can be referenced via a computed attribute for use with other resources. Refer to the `replication_specs[#].container_id` attribute in the [Attributes Reference](#attributes_reference) for more information.

-> **NOTE:** To enable Cluster Extended Storage Sizes use the `is_extended_storage_sizes_enabled` parameter in the [mongodbatlas_project resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project).

-> **NOTE:** The Low-CPU instance clusters are prefixed with `R`, for example `R40`. For complete list of Low-CPU instance clusters see Cluster Configuration Options under each [Cloud Provider](https://www.mongodb.com/docs/atlas/reference/cloud-providers).

-> **NOTE:** Groups and projects are synonymous terms. You might find group_id in the official documentation.

-> **NOTE:** This resource supports Flex clusters. Additionally, you can upgrade [M0 clusters to Flex](#example-tenant-cluster-upgrade-to-flex) and [Flex clusters to Dedicated](#Example-Flex-Cluster-Upgrade). When creating a Flex cluster, make sure to set the priority value to 7.

## Example Usage

### Example single provider and single region

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
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

### Example using effective fields with auto-scaling

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id            = var.project_id
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

# Read the effective (actual) values after Atlas scales
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

**For module authors:** See the [Effective Fields Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/effective_fields) for complete examples of using `use_effective_fields` and effective specs in reusable Terraform modules.

### Example Tenant Cluster

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
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
resource "mongodbatlas_advanced_cluster" "this" {
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
resource "mongodbatlas_advanced_cluster" "this" {
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
resource "mongodbatlas_advanced_cluster" "this" {
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
resource "mongodbatlas_advanced_cluster" "this" {
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
resource "mongodbatlas_advanced_cluster" "this" {
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
resource "mongodbatlas_advanced_cluster" "this" {
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
resource "mongodbatlas_advanced_cluster" "this" {
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
* [GCP Private Endpoint (Port-Mapped Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)
* [Azure Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/azure)
* [AWS, Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster)
* [AWS, Regionalized Private Endpoints](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster-geosharded)


### Further Examples

**Cluster Types:**
- [Replicaset](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/replicaset)
- [Symmetric Sharded Cluster](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/symmetric-sharded-cluster)
- [Asymmetric Sharded Cluster](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/asymmetric-sharded-cluster)
- [Global Cluster](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/global-cluster)
- [Multi-Cloud](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/multi-cloud)

**Auto-scaling:**
- [Auto-Scaling Per Shard](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/auto-scaling-per-shard)
- [Effective Fields Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/effective_fields)

**Upgrades & Migrations:**
- [Tenant Upgrade](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/tenant-upgrade)
- [Flex Upgrade](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/flex-upgrade)
- [Version Upgrade with Pinned FCV](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/version-upgrade-with-pinned-fcv)
- [Migrate Cluster to Advanced Cluster](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/migrate_cluster_to_advanced_cluster/basic)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_type` (String) Configuration of nodes that comprise the cluster.
- `name` (String) Human-readable label that identifies this cluster.
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.

**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
- `replication_specs` (Attributes List) List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations. (see [below for nested schema](#nestedatt--replication_specs))

### Optional

- `accept_data_risks_and_force_replica_set_reconfig` (String) If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set **acceptDataRisksAndForceReplicaSetReconfig** to the current date.
- `advanced_configuration` (Attributes) Additional settings for an Atlas cluster. (see [below for nested schema](#nestedatt--advanced_configuration))
- `backup_enabled` (Boolean) Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.
- `bi_connector_config` (Attributes) Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster. (see [below for nested schema](#nestedatt--bi_connector_config))
- `config_server_management_mode` (String) Config Server Management Mode for creating or updating a sharded cluster.

When configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.

When configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.
- `delete_on_create_timeout` (Boolean) Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.
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
- `retain_backups_enabled` (Boolean) Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster.
- `root_cert_type` (String) Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group.
- `tags` (Map of String) Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.
- `termination_protection_enabled` (Boolean) Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `use_aws_time_based_snapshot_copy_for_fast_initial_sync` (Boolean) Flag that indicates whether time-based snapshot copies will be used instead of slower standard snapshot copies during fast Atlas cross-region initial syncs. This flag is only relevant for clusters containing AWS nodes.
- `use_effective_fields` (Boolean) Controls how hardware specification fields are returned in the response. When set to true, the non-effective specs (`electable_specs`, `read_only_specs`, `analytics_specs`) fields return the hardware specifications that the client provided. When set to false (default), the non-effective specs fields show the **current** hardware specifications. Cluster auto-scaling is the primary cause for differences between initial and current hardware specifications.
- `version_release_system` (String) Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify **mongoDBMajorVersion**.

### Read-Only

- `cluster_id` (String) Unique 24-hexadecimal digit string that identifies the cluster.
- `config_server_type` (String) Describes a sharded cluster's config server type.
- `connection_strings` (Attributes) Collection of Uniform Resource Locators that point to the MongoDB database. (see [below for nested schema](#nestedatt--connection_strings))
- `create_date` (String) Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.
- `mongo_db_version` (String) Version of MongoDB that the cluster runs.
- `state_name` (String) Human-readable label that indicates the current operating condition of this cluster.

<a id="nestedatt--replication_specs"></a>
### Nested Schema for `replication_specs`

Required:

- `region_configs` (Attributes List) Hardware specifications for nodes set for a given region. Each **regionConfigs** object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each **regionConfigs** object must have either an **analyticsSpecs** object, **electableSpecs** object, or **readOnlySpecs** object. Tenant clusters only require **electableSpecs. Dedicated** clusters can specify any of these specifications, but must have at least one **electableSpecs** object within a **replicationSpec**.

**Example:**

If you set `"replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize" : "M30"`, set `"replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize" : `"M30"` if you have electable nodes and `"replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize" : `"M30"` if you have read-only nodes. (see [below for nested schema](#nestedatt--replication_specs--region_configs))

Optional:

- `zone_name` (String) Human-readable label that describes the zone this shard belongs to in a Global Cluster. Provide this value only if "clusterType" : "GEOSHARDED" but not "selfManagedSharding" : true.

Read-Only:

- `container_id` (Map of String) A key-value map of the Network Peering Container ID(s) for the configuration specified in region_configs. The Container ID is the id of the container created when the first cluster in the region (AWS/Azure) or project (GCP) was created.
- `external_id` (String) Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. This value corresponds to Shard ID displayed in the UI.
- `zone_id` (String) Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.

<a id="nestedatt--replication_specs--region_configs"></a>
### Nested Schema for `replication_specs.region_configs`

Required:

- `priority` (Number) Precedence is given to this region when a primary election occurs. If your **regionConfigs** has only **readOnlySpecs**, **analyticsSpecs**, or both, set this value to `0`. If you have multiple **regionConfigs** objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.

**Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.
- `provider_name` (String) Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.
- `region_name` (String) Physical location of your MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. The region name is only returned in the response for single-region clusters. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. It assigns the VPC a Classless Inter-Domain Routing (CIDR) block. To limit a new VPC peering connection to one Classless Inter-Domain Routing (CIDR) block and region, create the connection first. Deploy the cluster after the connection starts. GCP Clusters and Multi-region clusters require one VPC peering connection for each region. MongoDB nodes can use only the peering connection that resides in the same region as the nodes to communicate with the peered VPC.

Optional:

- `analytics_auto_scaling` (Attributes) Options that determine how this cluster handles resource scaling. (see [below for nested schema](#nestedatt--replication_specs--region_configs--analytics_auto_scaling))
- `analytics_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--analytics_specs))
- `auto_scaling` (Attributes) Options that determine how this cluster handles resource scaling. (see [below for nested schema](#nestedatt--replication_specs--region_configs--auto_scaling))
- `backing_provider_name` (String) Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when **providerName** is `TENANT` and **electableSpecs.instanceSize** is `M0`.
- `electable_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--electable_specs))
- `read_only_specs` (Attributes) Hardware specifications for nodes deployed in the region. (see [below for nested schema](#nestedatt--replication_specs--region_configs--read_only_specs))

<a id="nestedatt--replication_specs--region_configs--analytics_auto_scaling"></a>
### Nested Schema for `replication_specs.region_configs.analytics_auto_scaling`

Optional:

- `compute_enabled` (Boolean) Flag that indicates whether someone enabled instance size auto-scaling.

- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.
- Set to `false` to disable instance size automatic scaling.
- `compute_max_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_min_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_scale_down_enabled` (Boolean) Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.
- `disk_gb_enabled` (Boolean) Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.


<a id="nestedatt--replication_specs--region_configs--analytics_specs"></a>
### Nested Schema for `replication_specs.region_configs.analytics_specs`

Optional:

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

Optional:

- `compute_enabled` (Boolean) Flag that indicates whether someone enabled instance size auto-scaling.

- Set to `true` to enable instance size auto-scaling. If enabled, you must specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize**.
- Set to `false` to disable instance size automatic scaling.
- `compute_max_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_min_instance_size` (String) Minimum instance size to which your cluster can automatically scale. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.scaleDownEnabled" : true`.
- `compute_scale_down_enabled` (Boolean) Flag that indicates whether the instance size may scale down. MongoDB Cloud requires this parameter if `"replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled" : true`. If you enable this option, specify a value for **replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize**.
- `disk_gb_enabled` (Boolean) Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.


<a id="nestedatt--replication_specs--region_configs--electable_specs"></a>
### Nested Schema for `replication_specs.region_configs.electable_specs`

Optional:

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

Optional:

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




<a id="nestedatt--advanced_configuration"></a>
### Nested Schema for `advanced_configuration`

Optional:

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

Optional:

- `enabled` (Boolean) Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.
- `read_preference` (String) Data source node designated for the MongoDB Connector for Business Intelligence on MongoDB Cloud. The MongoDB Connector for Business Intelligence on MongoDB Cloud reads data from the primary, secondary, or analytics node based on your read preferences. Defaults to `ANALYTICS` node, or `SECONDARY` if there are no `ANALYTICS` nodes.


<a id="nestedatt--pinned_fcv"></a>
### Nested Schema for `pinned_fcv`

Required:

- `expiration_date` (String) Expiration date of the fixed FCV. This value is in the ISO 8601 timestamp format (e.g. 2024-12-04T16:25:00Z). Note that this field cannot exceed 4 weeks from the pinned date.

Read-Only:

- `version` (String) Feature compatibility version of the cluster.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


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

## Auto-Scaling with Effective Fields

The `use_effective_fields` attribute enhances auto-scaling workflows by eliminating the need for `lifecycle.ignore_changes` blocks and providing visibility into Atlas-managed changes. This feature only applies to dedicated clusters (M10+) and is not supported for flex and tenant clusters.

### Why use_effective_fields?

When auto-scaling is enabled on a cluster, Atlas automatically adjusts instance sizes and disk capacity based on workload. Without `use_effective_fields`, `lifecycle.ignore_changes` blocks are required to prevent Terraform from reverting these Atlas-managed changes. This approach has limitations:

- **Configuration drift**: The actual cluster configuration diverges from your Terraform configuration
- **Maintenance overhead**: Careful management of `ignore_changes` blocks is required, including commenting and uncommenting when making intentional changes
- **Limited visibility**: Actual scaled values cannot be easily inspected within Terraform state

### How use_effective_fields works

The `use_effective_fields` attribute changes how the provider handles specification attributes:

**When `use_effective_fields = false` (default - current behavior):**
- Spec attributes (`electable_specs`, `analytics_specs`, `read_only_specs`) behavior:
  - If values are specified in your Terraform configuration (e.g., `instance_size = "M10"`), those values remain in your configuration
  - If values are not specified, Atlas provides default values automatically
- With auto-scaling enabled, Atlas scales your cluster but your configured values do not update to match
- This creates plan drift: Terraform shows differences between your configured values and what Atlas has actually deployed
- `lifecycle.ignore_changes` must be used to prevent Terraform from reverting Atlas auto-scaling changes back to your original configuration

**When `use_effective_fields = true` (new behavior):**
- **Clear separation of concerns**:
  - Spec attributes remain exactly as defined in your Terraform configuration
  - Atlas-computed values (defaults and auto-scaled values) are available separately in effective specs
- No plan drift occurs when Atlas auto-scales your cluster
- Use data sources to read `effective_electable_specs`, `effective_analytics_specs`, and `effective_read_only_specs` for actual values

**Key difference:** With `use_effective_fields = true`, your configuration stays clean and represents your intent, while effective specs show the reality of what Atlas has provisioned.

See the [Example using effective fields with auto-scaling](#example-using-effective-fields-with-auto-scaling) in the Example Usage section.

### Manually Updating Specs with use_effective_fields

When `use_effective_fields = true` and auto-scaling is enabled, you can update `instance_size`, `disk_size_gb`, or `disk_iops` in your configuration at any time without validation errors. However, Atlas echoes these values back in state while continuing to use auto-scaled values for actual cluster operations. To have your configured values take effect, temporarily disable auto-scaling:

1. Set `compute_enabled = false` and `disk_gb_enabled = false` in the [`auto_scaling`](#auto_scaling) block and apply.
2. Update `instance_size`, `disk_size_gb`, or `disk_iops` to your desired values and apply.
3. Re-enable auto-scaling by setting `compute_enabled` and/or `disk_gb_enabled` back to `true` and apply.

This workflow allows you to set specific baseline values from which auto-scaling will resume dynamic adjustments based on workload.

### Terraform Modules

`use_effective_fields` is particularly valuable for reusable Terraform modules. Without it, separate module implementations are required (one with lifecycle blocks for auto-scaling, one without). With `use_effective_fields`, a single module handles both scenarios without lifecycle blocks. See the [Effective Fields Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_advanced_cluster/effective_fields) for complete implementations.

### Migration path and version 3.x

**Current behavior (provider v2.x):**
- `use_effective_fields` defaults to `false` for full backward compatibility
- Set to `true` to opt into the effective fields behavior
- The attribute will be deprecated later in v2.x releases in preparation for v3.x

**Future behavior (provider v3.x):**
- The effective fields behavior will be enabled by default
- The `use_effective_fields` attribute will be removed, as the new behavior becomes standard
- This change will reduce plan verbosity by making specification fields Optional-only (removing Computed), eliminating unnecessary `(known after apply)` markers for user-configured values

**Potential enhancements (v3.x or later):**
- If customer demand warrants, effective spec fields (`effective_electable_specs`, `effective_analytics_specs`, `effective_read_only_specs`) may be exposed directly in the resource (currently available only via data source)
- This would improve observability by providing direct access to actual operational values from the resource without requiring a separate data source
- Note: Effective fields would still show `(known after apply)` markers, but user-configured spec fields would not, resulting in clearer plan output overall

**Migration recommendation:** Adopt `use_effective_fields = true` in v2.x to prepare for the v3.x transition and benefit from improved auto-scaling workflows immediately. The recommendation is to toggle the flag and remove any existing `lifecycle.ignore_changes` blocks in the same apply, without combining other changes.

## Considerations and Best Practices

### "known after apply" verbosity

When modifying cluster configurations, you may see `(known after apply)` markers in your Terraform plan output, even for attributes you haven't modified. This is expected behavior, for example:
```
# mongodbatlas_advanced_cluster.this will be updated in-place
! resource "mongodbatlas_advanced_cluster" "this" {
!       connection_strings                   = {
+           private          = (known after apply)
!           private_endpoint = [
-               {
-                   connection_string                     = "<REDACTED>" -> null
-                   endpoints                             = [
-                       {
-                           endpoint_id   = "<REDACTED>" -> null
-                           provider_name = "AWS" -> null
-                           region        = "EU_EAST_1" -> null
                        },
                    ] -> null
                    # (1 unchanged attribute hidden)
                },
            ] -> (known after apply)
+          
...
!                       electable_specs        = {
!                           disk_iops       = 3000 -> (known after apply)
!                           disk_size_gb    = 60 -> 80  # CHANGE DONE IN THE CONFIGURATION FILE
!                           ebs_volume_type = "STANDARD" -> (known after apply)
                            # (2 unchanged attributes hidden)
                        }
...
    }
```

The provider v2.x uses the Terraform [Plugin Framework (TPF)](https://developer.hashicorp.com/terraform/plugin/framework), which is more strict and verbose with computed values than the legacy [SDKv2 framework](https://developer.hashicorp.com/terraform/plugin/sdkv2) used in v1.x. For more information, see [this discussion](https://discuss.hashicorp.com/t/best-practices-for-handling-known-after-apply-plan-verbosity-in-tpf-resources/73806). Key points:

- "(known after apply)" doesn't mean the value will change - It indicates a computed value that [can't be known in advance](https://developer.hashicorp.com/terraform/language/expressions/references#values-not-yet-known), even if the value remains the same.
- All attributes which are marked as "known after apply", including their nested attributes, can be safely ignored.
- Dependent attributes may change - Some changes can affect related attributes (e.g., change to `zone_name` may update `zone_id`, `region_name` may update `container_id`, `instance_size` may update `disk_iops`, or `provider_name` may update `ebs_volume_type`).
- Optional/Computed attributes show as "known after apply" when not explicitly set, but only attributes modified in the Terraform configuration files will change along with their dependent attributes.

To reduce the number of `(known after apply)` entries in your plan output, explicitly declare known values in your configuration where possible:
   ```terraform
   replication_specs = [
     {
       region_configs = [
         {
           electable_specs = {
             instance_size   = "M30"
             node_count      = 3
             disk_size_gb    = 100  # Explicitly set if known
             disk_iops       = 3000 # Explicitly set if known
             ebs_volume_type = "STANDARD" # Explicitly set even if it's the default
           }
           # ... other configuration
         }
       ]
     }
   ]
   ```

The MongoDB team is working to reduce plan verbosity, though no timeline is available yet.

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
