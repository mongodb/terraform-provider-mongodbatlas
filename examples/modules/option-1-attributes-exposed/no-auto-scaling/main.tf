

resource "mongodbatlas_advanced_cluster" "this" {
  accept_data_risks_and_force_replica_set_reconfig = var.accept_data_risks_and_force_replica_set_reconfig
  advanced_configuration                           = var.advanced_configuration
  backup_enabled                                   = var.backup_enabled
  bi_connector_config                              = var.bi_connector_config
  cluster_type                                     = var.cluster_type
  config_server_management_mode                    = var.config_server_management_mode
  delete_on_create_timeout                         = var.delete_on_create_timeout
  disk_size_gb                                     = var.disk_size_gb
  encryption_at_rest_provider                      = var.encryption_at_rest_provider
  global_cluster_self_managed_sharding             = var.global_cluster_self_managed_sharding
  labels                                           = var.labels
  mongo_db_major_version                           = var.mongo_db_major_version
  name                                             = var.name
  paused                                           = var.paused
  pinned_fcv                                       = var.pinned_fcv
  pit_enabled                                      = var.pit_enabled
  project_id                                       = var.project_id
  redact_client_log_data                           = var.redact_client_log_data
  replica_set_scaling_strategy                     = var.replica_set_scaling_strategy
  replication_specs                                = var.replication_specs
  retain_backups_enabled                           = var.retain_backups_enabled
  root_cert_type                                   = var.root_cert_type
  tags                                             = var.tags
  termination_protection_enabled                   = var.termination_protection_enabled
  timeouts                                         = var.timeouts
  version_release_system                           = var.version_release_system

  // no lifecycle ingore, auto-scaling will causes non-empty plans when using this module
}

