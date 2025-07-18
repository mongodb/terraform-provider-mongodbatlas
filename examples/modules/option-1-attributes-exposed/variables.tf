variable "accept_data_risks_and_force_replica_set_reconfig" {
  description = "If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forcedreconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set **acceptDataRisksAndForceReplicaSetReconfig** to the current date."
  type        = string
  nullable    = true
  default     = null
}

variable "advanced_configuration" {
  description = "Additional settings for an Atlas cluster."
  type = object({
    change_stream_options_pre_and_post_images_expire_after_seconds = optional(number)
    custom_openssl_cipher_config_tls12                             = optional(list(string))
    default_max_time_ms                                            = optional(number)
    default_read_concern                                           = optional(string)
    default_write_concern                                          = optional(string, "majority")
    fail_index_key_too_long                                        = optional(bool)
    javascript_enabled                                             = optional(bool, false)
    minimum_enabled_tls_protocol                                   = optional(string, "TLS1_2")
    no_table_scan                                                  = optional(bool)
    oplog_min_retention_hours                                      = optional(number)
    oplog_size_mb                                                  = optional(number)
    sample_refresh_interval_bi_connector                           = optional(number)
    sample_size_bi_connector                                       = optional(number)
    tls_cipher_config_mode                                         = optional(string, "DEFAULT")
    transaction_lifetime_limit_seconds                             = optional(number)
  })
  nullable = true
  default = {
    default_write_concern        = "majority"
    javascript_enabled           = false
    minimum_enabled_tls_protocol = "TLS1_2"
    tls_cipher_config_mode       = "DEFAULT"
  }
}

variable "backup_enabled" {
  description = "Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses [Cloud Backups](https://docs.atlas.mongodb.com/backup/cloud-backup/overview/) for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups."
  type        = bool
  nullable    = true
  default     = null
}

variable "bi_connector_config" {
  description = "Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster."
  type = object({
    enabled         = optional(bool)
    read_preference = optional(string)
  })
  nullable = true
  default  = null
}

variable "cluster_type" {
  description = "Configuration of nodes that comprise thecluster."
  type        = string
}

variable "config_server_management_mode" {
  description = <<-EOT
Config Server Management Mode for creating or updating a sharded cluster.

When configured as ATLAS_MANAGED, atlas may automatically switch the cluster's config server type for optimal performance and savings.

When configured as FIXED_TO_DEDICATED, the cluster will always use a dedicated config server.
EOT

  type     = string
  nullable = true
  default  = null
}

variable "delete_on_create_timeout" {
  description = "Flag that indicates whether to delete the cluster if the cluster creation times out. Default is false."
  type        = bool
  nullable    = true
  default     = null
}

variable "disk_size_gb" {
  description = <<-EOT
Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.

 This value must be equal for all shards and node types.

 This value is not configurable on M0/M2/M5 clusters.

 MongoDB Cloud requires thisparameter if you set **replicationSpecs**.

 If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value. 

 Storage charge calculations depend on whether you choose the default value or a custom value.

 The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
EOT

  type     = number
  nullable = true
  default  = null
}

variable "encryption_at_rest_provider" {
  description = "Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster **replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize** setting must be `M10` or higher and `\"backupEnabled\" : false` or omittedentirely."
  type        = string
  nullable    = true
  default     = null
}

variable "global_cluster_self_managed_sharding" {
  description = <<-EOT
Set this field to configure the Sharding Management Mode when creating a new Global Cluster.

When set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.

When set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.

This setting cannot be changed once the cluster is deployed.
EOT

  type     = bool
  nullable = true
  default  = null
}

variable "labels" {
  description = <<-EOT
Map of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.

Cluster labels are deprecated andwill be removed in a future release. We strongly recommend that you use [resource tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas) instead.
EOT

  type     = map(any)
  nullable = true
  default  = null
}

variable "mongo_db_major_version" {
  description = <<-EOT
MongoDB major version of the cluster.

On creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLTSVersions).

 On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, theMongoDB version can be downgraded to the previous major version.
EOT

  type     = string
  nullable = true
  default  = null
}

variable "name" {
  description = "Human-readable label that identifies this cluster."
  type        = string
  nullable    = true
}

variable "paused" {
  description = "Flag that indicates whether the cluster is paused."
  type        = bool
  nullable    = true
  default     = null
}

variable "pinned_fcv" {
  description = "Pins the Feature Compatibility Version (FCV) to the current MongoDB version with a provided expiration date. To unpin the FCV the `pinned_fcv` attribute must be removed. This operation can take several minutes as the request processes through the MongoDB data plane. Once FCV is unpinned it will not be possible to downgrade the `mongo_db_major_version`. It is advised that updates to `pinned_fcv` are done isolated from other cluster changes. If a plan contains multiple changes, the FCV change will be applied first. If FCV is unpinned past the expiration date the `pinned_fcv` attribute must be removed. The following [knowledge hub article](https://kb.corp.mongodb.com/article/000021785/) and [FCV documentation](https://www.mongodb.com/docs/atlas/tutorial/major-version-change/#manage-feature-compatibility--fcv--during-upgrades) can be referenced for moredetails."
  type = object({
    expiration_date = string
  })
  nullable = true
  default  = null
}

variable "pit_enabled" {
  description = "Flag that indicates whether the cluster uses continuous cloud backups."
  type        = bool
  nullable    = true
  default     = null
}

variable "project_id" {
  description = <<-EOT
Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.

**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
EOT

  type     = string
}

variable "redact_client_log_data" {
  description = <<-EOT
Enable or disable log redaction.

This setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging.This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.

Use ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.

*Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.
EOT

  type     = bool
  nullable = true
  default  = null
}

variable "replica_set_scaling_strategy" {
  description = <<-EOT
Set this field to configure the replica set scaling mode for your cluster.

By default, Atlas scales under WORKLOAD_TYPE. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.

When configured as SEQUENTIAL, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitivesecondary reads.

When configured as NODE_TYPE, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.
EOT

  type     = string
  nullable = true
  default  = null
}

variable "replication_specs" {
  description = "List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations."
  type = list(object({
    num_shards = optional(number)
    region_configs = list(object({
      analytics_auto_scaling = optional(object({
        compute_enabled            = optional(bool)
        compute_max_instance_size  = optional(string)
        compute_min_instance_size  = optional(string)
        compute_scale_down_enabled = optional(bool)
        disk_gb_enabled            = optional(bool)
      }))
      analytics_specs = optional(object({
        disk_iops       = optional(number)
        disk_size_gb    = optional(number)
        ebs_volume_type = optional(string)
        instance_size   = optional(string)
        node_count      = optional(number)
      }))
      auto_scaling = optional(object({
        compute_enabled            = optional(bool)
        compute_max_instance_size  = optional(string)
        compute_min_instance_size  = optional(string)
        compute_scale_down_enabled = optional(bool)
        disk_gb_enabled            = optional(bool)
      }))
      backing_provider_name = optional(string)
      electable_specs = optional(object({
        disk_iops       = optional(number)
        disk_size_gb    = optional(number)
        ebs_volume_type = optional(string)
        instance_size   = optional(string)
        node_count      = optional(number)
      }))
      priority      = number
      provider_name = string
      read_only_specs = optional(object({
        disk_iops       = optional(number)
        disk_size_gb    = optional(number)
        ebs_volume_type = optional(string)
        instance_size   = optional(string)
        node_count      = optional(number)
      }))
      region_name = string
    }))
    zone_name = optional(string)
  }))
}

variable "retain_backups_enabled" {
  description = "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster."
  type        = bool
  nullable    = true
  default     = null
}

variable "root_cert_type" {
  description = "Root Certificate Authority that MongoDB Cloud cluster uses. MongoDB Cloud supports Internet Security Research Group."
  type        = string
  nullable    = true
  default     = null
}

variable "tags" {
  description = "Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster."
  type        = map(any)
  nullable    = true
  default     = null
}

variable "termination_protection_enabled" {
  description = "Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster."
  type        = bool
  nullable    = true
  default     = null
}

variable "timeouts" {
  type = object({
    create = optional(string)
    delete = optional(string)
    update = optional(string)
  })
  nullable = true
  default  = null
}

variable "version_release_system" {
  description = "Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify**mongoDBMajorVersion**."
  type        = string
  nullable    = true
  default     = null
}
