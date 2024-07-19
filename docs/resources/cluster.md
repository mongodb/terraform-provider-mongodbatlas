---
subcategory: "Deprecated"    
---

**WARNING:** This resource is deprecated, use `mongodbatlas_advanced_cluster`. In order to learn more about how to migrate, please read the [Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

# Resource: mongodbatlas_cluster

`mongodbatlas_cluster` provides a Cluster resource. The resource lets you create, edit and delete clusters. The resource requires your Project ID.

~> **IMPORTANT:**
<br> &#8226; New Users: If you are not already using `mongodbatlas_cluster` for your deployment we recommend starting with the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster).  `mongodbatlas_advanced_cluster` has all the same functionality as `mongodbatlas_cluster` but also supports multi-cloud clusters.  

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** A network container is created for a cluster to reside in. To use this container with another resource, such as peering, reference the computed`container_id` attribute on the cluster.

-> **NOTE:** To enable Cluster Extended Storage Sizes use the `is_extended_storage_sizes_enabled` parameter in the [mongodbatlas_project resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project).

-> **NOTE:** If Backup Compliance Policy is enabled for the project for which this backup schedule is defined, you cannot modify the backup schedule for an individual cluster below the minimum requirements set in the Backup Compliance Policy.  See [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

-> **NOTE:** The Low-CPU instance clusters are prefixed with `R`, i.e. `R40`. For complete list of Low-CPU instance clusters see Cluster Configuration Options under each Cloud Provider (https://www.mongodb.com/docs/atlas/reference/cloud-providers/).

~> **IMPORTANT:**
<br> &#8226; Multi Region Cluster: The `mongodbatlas_cluster` resource doesn't return the `container_id` for each region utilized by the cluster. For retrieving the `container_id`, we recommend to use the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource instead.
<br> &#8226; Free tier cluster creation (M0) is supported.
<br> &#8226; Shared tier clusters (M0, M2, M5) can be upgraded to dedicated tiers (M10+) via this provider. WARNING WHEN UPGRADING TENANT/SHARED CLUSTERS!!! Any change from shared tier to a different instance size will be considered a tenant upgrade. When upgrading from shared tier to dedicated simply change the `provider_name` from "TENANT"  to your preferred provider (AWS, GCP, AZURE) and remove the variable `backing_provider_name`, for example if you have an existing tenant/shared cluster and want to upgrade your Terraform config should be changed from:
```
provider_instance_size_name = "M0"
provider_name               = "TENANT"
backing_provider_name       = "AWS"
```
To:
```
provider_instance_size_name = "M10"
provider_name               = "AWS"
```
<br> &#8226; Changes to cluster configurations can affect costs. Before making changes, please see [Billing](https://docs.atlas.mongodb.com/billing/).   
<br> &#8226; If your Atlas project contains a custom role that uses actions introduced in a specific MongoDB version, you cannot create a cluster with a MongoDB version less than that version unless you delete the custom role.

## Example Usage

### Example AWS cluster

```terraform
resource "mongodbatlas_cluster" "cluster-test" {
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
  cloud_backup = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "7.0"

  # Provider Settings "block"
  provider_name               = "AWS"
  provider_instance_size_name = "M40"
}
```

### Example Azure cluster.

```terraform
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "test"
  cluster_type = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup     = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "7.0"

  # Provider Settings "block"
  provider_name               = "AZURE"
  provider_disk_type_name     = "P6"
  provider_instance_size_name = "M30"
}
```

### Example GCP cluster

```terraform
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "test"
  cluster_type = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "EASTERN_US"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup                 = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "7.0"

  # Provider Settings "block"
  provider_name               = "GCP"
  provider_instance_size_name = "M30"
}
```

### Example Multi Region cluster

```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id               = "<YOUR-PROJECT-ID>"
  name                     = "cluster-test-multi-region"
  num_shards               = 1
  cloud_backup             = true
  cluster_type             = "REPLICASET"

  # Provider Settings "block"
  provider_name               = "AWS"
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

```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id              = "<YOUR-PROJECT-ID>"
  name                    = "cluster-test-global"
  num_shards              = 1
  cloud_backup            = true
  cluster_type            = "GEOSHARDED"

  # Provider Settings "block"
  provider_name               = "AWS"
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
### Example AWS Shared Tier (M2/M5) cluster
```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id              = "<YOUR-PROJECT-ID>"
  name                    = "cluster-test-global"

  # Provider Settings "block"
  provider_name = "TENANT"
  backing_provider_name = "AWS"
  provider_region_name = "US_EAST_1"
  provider_instance_size_name = "M2"
}
```
### Example AWS Free Tier cluster
```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id              = "<YOUR-PROJECT-ID>"
  name                    = "cluster-test-global"

  # Provider Settings "block"
  provider_name = "TENANT"
  backing_provider_name = "AWS"
  provider_region_name = "US_EAST_1"
  provider_instance_size_name = "M0"
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

* `project_id` - (Required) The unique ID for the project to create the database user.
* `provider_name` - (Required) Cloud service provider on which the servers are provisioned.

    The possible values are:

    - `AWS` - Amazon AWS
    - `GCP` - Google Cloud Platform
    - `AZURE` - Microsoft Azure
    - `TENANT` - A multi-tenant deployment on one of the supported cloud service providers. Only valid when providerSettings.instanceSizeName is either M2 or M5.
* `name` - (Required) Name of the cluster as it appears in Atlas. Once the cluster is created, its name cannot be changed. **WARNING** Changing the name will result in destruction of the existing cluster and the creation of a new cluster.
* `provider_instance_size_name` - (Required) Atlas provides different instance sizes, each with a default storage capacity and RAM size. The instance size you select is used for all the data-bearing servers in your cluster. See [Create a Cluster](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/) `providerSettings.instanceSizeName` for valid values and default resources.

* `auto_scaling_disk_gb_enabled` - (Optional) Specifies whether disk auto-scaling is enabled. The default is true.
    - Set to `true` to enable disk auto-scaling.
    - Set to `false` to disable disk auto-scaling.

~> **IMPORTANT:** If `auto_scaling_disk_gb_enabled` is true, then Atlas will automatically scale disk size up and down.
This will cause the value of `disk_size_gb` returned to potentially be different than what is specified in the Terraform config and if one then applies a plan, not noting this, Terraform will scale the cluster disk size back to the original `disk_size_gb` value.
To prevent this a lifecycle customization should be used, i.e.:  
`lifecycle {
  ignore_changes = [disk_size_gb]
}`
After adding the `lifecycle` block to explicitly change `disk_size_gb` comment out the `lifecycle` block and run `terraform apply`. Please be sure to uncomment the `lifecycle` block once done to prevent any accidental changes.

-> **NOTE:** If `provider_name` is set to `TENANT`, the parameter `auto_scaling_disk_gb_enabled` will be ignored.

* `auto_scaling_compute_enabled` - (Optional) Specifies whether cluster tier auto-scaling is enabled. The default is false.
    - Set to `true` to enable cluster tier auto-scaling. If enabled, you must specify a value for `providerSettings.autoScaling.compute.maxInstanceSize`.
    - Set to `false` to disable cluster tier auto-scaling.

~> **IMPORTANT:** If `auto_scaling_compute_enabled` is true,  then Atlas will automatically scale up to the maximum provided and down to the minimum, if provided.
This will cause the value of `provider_instance_size_name` returned to potentially be different than what is specified in the Terraform config and if one then applies a plan, not noting this, Terraform will scale the cluster back to the original instanceSizeName value.
To prevent this a lifecycle customization should be used, i.e.:  
`lifecycle {
  ignore_changes = [provider_instance_size_name]
}`
But in order to explicitly change `provider_instance_size_name` comment the `lifecycle` block and run `terraform apply`. Please ensure to uncomment it to prevent any accidental changes.

* `auto_scaling_compute_scale_down_enabled` - (Optional) Set to `true` to enable the cluster tier to scale down. This option is only available if `autoScaling.compute.enabled` is `true`.
    - If this option is enabled, you must specify a value for `providerSettings.autoScaling.compute.minInstanceSize`

* `backup_enabled` - (Optional) Legacy Backup - Set to true to enable Atlas legacy backups for the cluster.
**Important** - MongoDB deprecated the Legacy Backup feature. Clusters that use Legacy Backup can continue to use it. MongoDB recommends using [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/).
    * New Atlas clusters of any type do not support this parameter. These clusters must use Cloud Backup, `cloud_backup`, to enable Cloud Backup.  If you create a new Atlas cluster and set `backup_enabled` to true, the Provider will respond with an error.  This change doesn’t affect existing clusters that use legacy backups.
    * Setting this value to false to disable legacy backups for the cluster will let Atlas delete any stored snapshots. In order to preserve the legacy backups snapshots, disable the legacy backups and enable the cloud backups in the single **terraform apply** action.
    ```
    backup_enabled = "false"
    cloud_backup = "true"
    ```
    * The default value is false.  M10 and above only.

* `retain_backups_enabled` - (Optional) Set to true to retain backup snapshots for the deleted cluster. M10 and above only. 
* `bi_connector_config` - (Optional) Specifies BI Connector for Atlas configuration on this cluster. BI Connector for Atlas is only available for M10+ clusters. See [BI Connector](#bi-connector) below for more details.
* `cluster_type` - (Required) Specifies the type of the cluster that you want to modify. You cannot convert a sharded cluster deployment to a replica set deployment.

    -> **WHEN SHOULD YOU USE CLUSTERTYPE?**
      When you set replication_specs, when you are deploying Global Clusters or when you are deploying non-Global replica sets and sharded clusters.

    Accepted values include:
      - `REPLICASET` Replica set
      - `SHARDED` Sharded cluster
      - `GEOSHARDED` Global Cluster

* `disk_size_gb` - (Optional - GCP/AWS Only) Capacity, in gigabytes, of the host’s root volume. Increase this number to add capacity, up to a maximum possible value of 4096 (i.e., 4 TB). This value must be a positive integer.
  * The minimum disk size for dedicated clusters is 10GB for AWS and GCP. If you specify diskSizeGB with a lower disk size, Atlas defaults to the minimum disk size value.
  * Note: The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require additional storage space beyond this limitation, consider upgrading your cluster to a higher tier.
  * Cannot be used with clusters with local NVMe SSDs
  * Cannot be used with Azure clusters
* `encryption_at_rest_provider` - (Optional) Possible values are AWS, GCP, AZURE or NONE.  Only needed if you desire to manage the keys, see [Encryption at Rest using Customer Key Management](https://docs.atlas.mongodb.com/security-aws-kms/) for complete documentation.  You must configure encryption at rest for the Atlas project before enabling it on any cluster in the project. For complete documentation on configuring Encryption at Rest, see Encryption at Rest using Customer Key Management. Requires M10 or greater. and for legacy backups, backup_enabled, to be false or omitted. **Note: Atlas encrypts all cluster storage and snapshot volumes, securing all cluster data on disk: a concept known as encryption at rest, by default**.   
* `tags` - (Optional) Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#tags).
* `labels` - (Optional) Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#labels). **DEPRECATED** Use `tags` instead.
* `mongo_db_major_version` - (Optional) Version of the cluster to deploy. Atlas supports the following MongoDB versions for M10+ clusters: `4.4`, `5.0`, `6.0` or `7.0`. If omitted, Atlas deploys a cluster that runs MongoDB 7.0. If `provider_instance_size_name`: `M0`, `M2` or `M5`, Atlas deploys MongoDB 5.0. Atlas always deploys the cluster with the latest stable release of the specified version. See [Release Notes](https://www.mongodb.com/docs/upcoming/release-notes/) for latest Current Stable Release.
* `num_shards` - (Optional) Selects whether the cluster is a replica set or a sharded cluster. If you use the replicationSpecs parameter, you must set num_shards.
* `pit_enabled` - (Optional) - Flag that indicates if the cluster uses Continuous Cloud Backup. If set to true, cloud_backup must also be set to true.
* `cloud_backup` - (Optional) Flag indicating if the cluster uses Cloud Backup for backups.

    If true, the cluster uses Cloud Backup for backups. If cloud_backup and backup_enabled are false, the cluster does not use Atlas backups.

    You cannot enable cloud backup if you have an existing cluster in the project with legacy backup enabled.

    ~> **IMPORTANT:** If setting to true for an existing cluster or imported cluster be sure to run terraform refresh after applying to enable modification of the Cloud Backup Snapshot Policy going forward.

* `backing_provider_name` - (Optional) Cloud service provider on which the server for a multi-tenant cluster is provisioned.

    This setting is only valid when providerSetting.providerName is TENANT and providerSetting.instanceSizeName is M2 or M5.

    The possible values are:

    - AWS - Amazon AWS
    - GCP - Google Cloud Platform
    - AZURE - Microsoft Azure

* `provider_disk_iops` - (Optional - AWS Only) The maximum input/output operations per second (IOPS) the system can perform. The possible values depend on the selected `provider_instance_size_name` and `disk_size_gb`.  This setting requires that `provider_instance_size_name` to be M30 or greater and cannot be used with clusters with local NVMe SSDs.  The default value for `provider_disk_iops` is the same as the cluster tier's Standard IOPS value, as viewable in the Atlas console.  It is used in cases where a higher number of IOPS is needed and possible.  If a value is submitted that is lower or equal to the default IOPS value for the cluster tier Atlas ignores the requested value and uses the default.  More details available under the providerSettings.diskIOPS parameter: [MongoDB API Clusters](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/)
  * You do not need to configure IOPS for a STANDARD disk configuration but only for a PROVISIONED configuration.

* `provider_disk_type_name` - (Optional - Azure Only) Azure disk type of the server’s root volume. If omitted, Atlas uses the default disk type for the selected providerSettings.instanceSizeName.  Example disk types and associated storage sizes: P4 - 32GB, P6 - 64GB, P10 - 128GB, P15 - 256GB, P20 - 512GB, P30 - 1024GB, P40 - 2048GB, P50 - 4095GB.  More information and the most update to date disk types/storage sizes can be located at https://docs.atlas.mongodb.com/reference/api/clusters-create-one/.
* `provider_encrypt_ebs_volume` - **(Deprecated) The Flag is always true.** Flag that indicates whether the Amazon EBS encryption feature encrypts the host's root volume for both data at rest within the volume and for data moving between the volume and the cluster. Note: This setting is always enabled for clusters with local NVMe SSDs. **Atlas encrypts all cluster storage and snapshot volumes, securing all cluster data on disk: a concept known as encryption at rest, by default.**.
* `provider_region_name` - (Optional) Physical location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases.  Requires the **Atlas region name**, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
    Do not specify this field when creating a multi-region cluster using the replicationSpec document or a Global Cluster with the replicationSpecs array.
* `provider_volume_type` - (AWS - Optional) The type of the volume. The possible values are: `STANDARD` and `PROVISIONED`.  `PROVISIONED` is ONLY required if setting IOPS higher than the default instance IOPS.
-> **NOTE:** `STANDARD` is not available for NVME clusters.

* `replication_factor` - (Deprecated) Number of replica set members. Each member keeps a copy of your databases, providing high availability and data redundancy. The possible values are 3, 5, or 7. The default value is 3.
* `provider_auto_scaling_compute_min_instance_size` - (Optional) Minimum instance size to which your cluster can automatically scale (e.g., M10). Required if `autoScaling.compute.scaleDownEnabled` is `true`.
* `provider_auto_scaling_compute_max_instance_size` - (Optional) Maximum instance size to which your cluster can automatically scale (e.g., M40). Required if `autoScaling.compute.enabled` is `true`.

* `replication_specs` - Configuration for cluster regions.  See [Replication Spec](#replication-spec) below for more details.
* `paused` (Optional) - Flag that indicates whether the cluster is paused or not. You can pause M10 or larger clusters.  You cannot initiate pausing for a shared/tenant tier cluster.  See [Considerations for Paused Clusters](https://docs.atlas.mongodb.com/pause-terminate-cluster/#considerations-for-paused-clusters)  
  **NOTE** Pause lasts for up to 30 days. If you don't resume the cluster within 30 days, Atlas resumes the cluster.  When the cluster resumption happens Terraform will flag the changed state.  If you wish to keep the cluster paused, reapply your Terraform configuration.   If you prefer to allow the automated change of state to unpaused use:
  `lifecycle {
  ignore_changes = [paused]
  }`
* `termination_protection_enabled` - Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
* `version_release_system` - (Optional) - Release cadence that Atlas uses for this cluster. This parameter defaults to `LTS`. If you set this field to `CONTINUOUS`, you must omit the `mongo_db_major_version` field. Atlas accepts:
  - `CONTINUOUS`:  Atlas creates your cluster using the most recent MongoDB release. Atlas automatically updates your cluster to the latest major and rapid MongoDB releases as they become available.
  - `LTS`: Atlas creates your cluster using the latest patch release of the MongoDB version that you specify in the mongoDBMajorVersion field. Atlas automatically updates your cluster to subsequent patch releases of this MongoDB version. Atlas doesn't update your cluster to newer rapid or major MongoDB releases as they become available.
* `timeouts`- (Optional) The duration of time to wait for Cluster to be created, updated, or deleted. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Cluster create & delete is `3h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).
* `accept_data_risks_and_force_replica_set_reconfig`- (Optional) If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set `accept_data_risks_and_force_replica_set_reconfig` to the current date. Learn more about Reconfiguring a Replica Set during a regional outage [here](https://dochub.mongodb.org/core/regional-outage-reconfigure-replica-set).

### Multi-Region Cluster

```terraform
//Example 3 Multi-Region block
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

```

**Replication Spec**  - Configuration block for multi-region cluster.

* `num_shards` - (Required) Number of shards up to 50 to deploy for a sharded cluster. The resource returns 1 to indicate a replica set and values of 2 and higher to indicate a sharded cluster. The returned value equals the number of shards in the cluster.
If you are upgrading a replica set to a sharded cluster, you cannot increase the number of shards in the same update request. You should wait until after the cluster has completed upgrading to sharded and you have reconnected all application clients to the MongoDB router before adding additional shards. Otherwise, your data might become inconsistent once MongoDB Cloud begins distributing data across shards. To learn more, see [Convert a replica set to a sharded cluster documentation](https://www.mongodb.com/docs/atlas/scale-cluster/#convert-a-replica-set-to-a-sharded-cluster) and [Convert a replica set to a sharded cluster tutorial](https://www.mongodb.com/docs/upcoming/tutorial/convert-replica-set-to-replicated-shard-cluster).
* `id` - (Optional) Unique identifer of the replication document for a zone in a Global Cluster.
* `regions_config` - (Optional) Physical location of the region. Each regionsConfig document describes the region’s priority in elections and the number and type of MongoDB nodes Atlas deploys to the region. You must order each regionsConfigs document by regionsConfig.priority, descending. See [Region Config](#region-config) below for more details.
* `zone_name` - (Optional) Name for the zone in a Global Cluster.


**Region Config**

* `region_name` - (Optional) Physical location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases.  Requires the **Atlas region name**, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
* `electable_nodes` - (Optional) Number of electable nodes for Atlas to deploy to the region. Electable nodes can become the primary and can facilitate local reads.
  * The total number of electableNodes across all replication spec regions  must total 3, 5, or 7.
  * Specify 0 if you do not want any electable nodes in the region.
  * You cannot create electable nodes in a region if `priority` is 0.
* `priority` - (Optional)  Election priority of the region. For regions with only read-only nodes, set this value to 0.
  * For regions where `electable_nodes` is at least 1, each region must have a priority of exactly one (1) less than the previous region. The first region must have a priority of 7. The lowest possible priority is 1.
  * The priority 7 region identifies the Preferred Region of the cluster. Atlas places the primary node in the Preferred Region. Priorities 1 through 7 are exclusive - no more than one region per cluster can be assigned a given priority.
  * Example: If you have three regions, their priorities would be 7, 6, and 5 respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be 4 and 3 respectively.
* `read_only_nodes` - (Optional) Number of read-only nodes for Atlas to deploy to the region. Read-only nodes can never become the primary, but can facilitate local-reads. Specify 0 if you do not want any read-only nodes in the region.
* `analytics_nodes` - (Optional) The number of analytics nodes for Atlas to deploy to the region. Analytics nodes are useful for handling analytic data such as reporting queries from BI Connector for Atlas. Analytics nodes are read-only, and can never become the primary. If you do not specify this option, no analytics nodes are deployed to the region.


### BI Connector

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
* **Note**  A minimum oplog retention is required when seeking to change a cluster's class to Local NVMe SSD. To learn more and for latest guidance see  [`oplogMinRetentionHours`](https://www.mongodb.com/docs/manual/core/replica-set-oplog/#std-label-replica-set-minimum-oplog-size) 
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

### Labels

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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - The cluster ID.
*  `mongo_db_version` - Version of MongoDB the cluster runs, in `major-version`.`minor-version` format.
* `id` -  The Terraform's unique identifier used internally for state management.
* `mongo_uri` - Base connection string for the cluster. Atlas only displays this field after the cluster is operational, not while it builds the cluster.
* `mongo_uri_updated` - Lists when the connection string was last updated. The connection string changes, for example, if you change a replica set to a sharded cluster.
* `mongo_uri_with_options` - connection string for connecting to the Atlas cluster. Includes the replicaSet, ssl, and authSource query parameters in the connection string with values appropriate for the cluster.
* `connection_strings` - Set of connection strings that your applications use to connect to this cluster. More info in [Connection-strings](https://docs.mongodb.com/manual/reference/connection-string/). Use the parameters in this object to connect your applications to this cluster. To learn more about the formats of connection strings, see [Connection String Options](https://docs.atlas.mongodb.com/reference/faq/connection-changes/). NOTE: Atlas returns the contents of this object after the cluster is operational, not while it builds the cluster.

   **NOTE** Connection strings must be returned as a list, therefore to refer to a specific attribute value add index notation. Example: mongodbatlas_cluster.cluster-test.connection_strings.0.standard_srv

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
* `container_id` - The Container ID is the id of the container created when the first cluster in the region (AWS/Azure) or project (GCP) was created.
* `srv_address` - Connection string for connecting to the Atlas cluster. The +srv modifier forces the connection to use TLS/SSL. See the mongoURI for additional options.
* `state_name` - Current state of the cluster. The possible states are:
    - IDLE
    - CREATING
    - UPDATING
    - DELETING
    - DELETED
    - REPAIRING

### Cloud Backup Policy

**WARNING:** This property is deprecated, use `mongodbatlas_cloud_backup_schedule` resource instead.

Cloud Backup Policy will be added if `cloud_backup` is enabled because MongoDB Atlas automatically creates a default policy, if not, returned values will be empty.

* `snapshot_backup_policy` - current snapshot schedule and retention settings for the cluster.

* `snapshot_backup_policy.#.cluster_id` - Unique identifier of the Atlas cluster.
* `snapshot_backup_policy.#.cluster_name` - Name of the Atlas cluster that contains the snapshot backup policy.
* `snapshot_backup_policy.#.next_snapshot` - UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
* `snapshot_backup_policy.#.reference_hour_of_day` - UTC Hour of day between 0 and 23 representing which hour of the day that Atlas takes a snapshot.
* `snapshot_backup_policy.#.reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the referenceHourOfDay that Atlas takes the snapshot.
* `snapshot_backup_policy.#.restore_window_days` - Specifies a restore window in days for the cloud provider backup to maintain.
* `snapshot_backup_policy.#.update_snapshots` - Specifies it's true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.

### Policies
**WARNING:** This property is deprecated, use `mongodbatlas_cloud_backup_schedule` resource instead.

* `snapshot_backup_policy.#.policies` - A list of policy definitions for the cluster.
* `snapshot_backup_policy.#.policies.#.id` - Unique identifier of the backup policy.

#### Policy Item
**WARNING:** This property is deprecated, use `mongodbatlas_cloud_backup_schedule` resource instead.

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
