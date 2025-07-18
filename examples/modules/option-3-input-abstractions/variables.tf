variable "region_configs" {
  description = "List of region configurations for a replica set cluster. Each object defines a region."
  type = list(object({
    provider_name        = string
    region_name          = string
    instance_size        = string
    priority             = optional(number, 7) # required if you have more than one region
    ebs_volume_type      = optional(string)
    disk_size_gb         = optional(number)
    disk_iops            = optional(number)
    electable_node_count = number
    read_only_node_count = optional(number, 0)
    analytics_specs = optional(object({
      instance_size   = string
      ebs_volume_type = optional(string)
      disk_size_gb    = optional(number)
      disk_iops       = optional(number)
      node_count      = number
    }))
  }))
  default = []
}

variable "shards" {
  description = "List of shard configurations for a sharded cluster. Each object defines a shard and its regions."
  type = list(object({
    zone_name = optional(string)
    region_configs = list(object({
      provider_name        = string
      region_name          = string
      instance_size        = string
      priority             = optional(number, 7) # required if you have more than one region
      ebs_volume_type      = optional(string)
      disk_size_gb         = optional(number)
      disk_iops            = optional(number)
      electable_node_count = number
      read_only_node_count = optional(number, 0)
      analytics_specs = optional(object({
        instance_size   = string
        ebs_volume_type = optional(string)
        disk_size_gb    = optional(number)
        disk_iops       = optional(number)
        node_count      = number
      }))
    }))
  }))
  default = []
}

variable "auto_scaling" {
  description = "Configuration for auto-scaling."
  type = object({
    disk_gb_enabled            = bool
    compute_enabled            = bool
    compute_scale_down_enabled = optional(bool)
    compute_max_instance_size  = optional(string)
    compute_min_instance_size  = optional(string)
  })
  default = {
    disk_gb_enabled = false # defaults to false, default to true would imply being opinionated on a max_instance_size which can vary significantly
    compute_enabled = false
  }
}

variable "analytics_auto_scaling" {
  description = "Configuration for analytics auto-scaling."
  type = object({
    disk_gb_enabled            = bool
    compute_enabled            = bool
    compute_scale_down_enabled = optional(bool)
    compute_max_instance_size  = optional(string)
    compute_min_instance_size  = optional(string)
  })
  default = {
    disk_gb_enabled = false # defaults to false, default to true would imply being opinionated on a max_instance_size which can vary significantly
    compute_enabled = false
  }
}

variable "project_id" {
  description = "The MongoDB Atlas project ID."
  type        = string
}

variable "name" {
  description = "The name of the cluster."
  type        = string
}

variable "cluster_type" {
  description = "Type of cluster (REPLICASET, SHARDED, GEOSHARDED)."
  type        = string
  default     = "REPLICASET"
}

variable "mongo_db_major_version" {
  description = "The major MongoDB version."
  type        = string
  default     = "8.0"
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

variable "tags" { // TODO more opinionated abstraction will be done here
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
