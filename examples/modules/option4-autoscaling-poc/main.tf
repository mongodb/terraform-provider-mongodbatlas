
locals {
  mongodbatlas_advanced_cluster_varsx = {
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
  }
  mongodbatlas_advanced_cluster_vars = {
    auto_scaling            = var.auto_scaling
    auto_scaling_analytics  = var.auto_scaling_analytics
    instance_size           = var.instance_size
    instance_size_analytics = var.instance_size_analytics
    provider_name           = var.provider_name
    regions                 = var.regions
  }
}


data "external" "mongodbatlas_advanced_cluster" {
  program = ["python3", "${path.module}/mongodbatlas_advanced_cluster.py"]
  query = {
    input_json = jsonencode(merge(local.mongodbatlas_advanced_cluster_vars, local.mongodbatlas_advanced_cluster_varsx, local.existing_cluster))
  }
}



resource "mongodbatlas_advanced_cluster" "this" {
  lifecycle {
    precondition {
      condition     = length(data.external.mongodbatlas_advanced_cluster.result.error_message) == 0
      error_message = data.external.mongodbatlas_advanced_cluster.result.error_message
    }
  }

  cluster_type                                     = data.external.mongodbatlas_advanced_cluster.result.cluster_type
  name                                             = data.external.mongodbatlas_advanced_cluster.result.name
  project_id                                       = data.external.mongodbatlas_advanced_cluster.result.project_id
  replication_specs                                = jsondecode(data.external.mongodbatlas_advanced_cluster.result.replication_specs)
  accept_data_risks_and_force_replica_set_reconfig = data.external.mongodbatlas_advanced_cluster.result.accept_data_risks_and_force_replica_set_reconfig == "" ? null : data.external.mongodbatlas_advanced_cluster.result.accept_data_risks_and_force_replica_set_reconfig
  advanced_configuration                           = data.external.mongodbatlas_advanced_cluster.result.advanced_configuration == "" ? null : jsondecode(data.external.mongodbatlas_advanced_cluster.result.advanced_configuration)
  backup_enabled                                   = data.external.mongodbatlas_advanced_cluster.result.backup_enabled == "" ? null : data.external.mongodbatlas_advanced_cluster.result.backup_enabled
  bi_connector_config                              = data.external.mongodbatlas_advanced_cluster.result.bi_connector_config == "" ? null : jsondecode(data.external.mongodbatlas_advanced_cluster.result.bi_connector_config)
  config_server_management_mode                    = data.external.mongodbatlas_advanced_cluster.result.config_server_management_mode == "" ? null : data.external.mongodbatlas_advanced_cluster.result.config_server_management_mode
  delete_on_create_timeout                         = data.external.mongodbatlas_advanced_cluster.result.delete_on_create_timeout == "" ? null : data.external.mongodbatlas_advanced_cluster.result.delete_on_create_timeout
  disk_size_gb                                     = data.external.mongodbatlas_advanced_cluster.result.disk_size_gb == "" ? null : data.external.mongodbatlas_advanced_cluster.result.disk_size_gb
  encryption_at_rest_provider                      = data.external.mongodbatlas_advanced_cluster.result.encryption_at_rest_provider == "" ? null : data.external.mongodbatlas_advanced_cluster.result.encryption_at_rest_provider
  global_cluster_self_managed_sharding             = data.external.mongodbatlas_advanced_cluster.result.global_cluster_self_managed_sharding == "" ? null : data.external.mongodbatlas_advanced_cluster.result.global_cluster_self_managed_sharding
  labels                                           = data.external.mongodbatlas_advanced_cluster.result.labels == "" ? null : jsondecode(data.external.mongodbatlas_advanced_cluster.result.labels)
  mongo_db_major_version                           = data.external.mongodbatlas_advanced_cluster.result.mongo_db_major_version == "" ? null : data.external.mongodbatlas_advanced_cluster.result.mongo_db_major_version
  paused                                           = data.external.mongodbatlas_advanced_cluster.result.paused == "" ? null : data.external.mongodbatlas_advanced_cluster.result.paused
  pinned_fcv                                       = data.external.mongodbatlas_advanced_cluster.result.pinned_fcv == "" ? null : jsondecode(data.external.mongodbatlas_advanced_cluster.result.pinned_fcv)
  pit_enabled                                      = data.external.mongodbatlas_advanced_cluster.result.pit_enabled == "" ? null : data.external.mongodbatlas_advanced_cluster.result.pit_enabled
  redact_client_log_data                           = data.external.mongodbatlas_advanced_cluster.result.redact_client_log_data == "" ? null : data.external.mongodbatlas_advanced_cluster.result.redact_client_log_data
  replica_set_scaling_strategy                     = data.external.mongodbatlas_advanced_cluster.result.replica_set_scaling_strategy == "" ? null : data.external.mongodbatlas_advanced_cluster.result.replica_set_scaling_strategy
  retain_backups_enabled                           = data.external.mongodbatlas_advanced_cluster.result.retain_backups_enabled == "" ? null : data.external.mongodbatlas_advanced_cluster.result.retain_backups_enabled
  root_cert_type                                   = data.external.mongodbatlas_advanced_cluster.result.root_cert_type == "" ? null : data.external.mongodbatlas_advanced_cluster.result.root_cert_type
  tags                                             = data.external.mongodbatlas_advanced_cluster.result.tags == "" ? null : jsondecode(data.external.mongodbatlas_advanced_cluster.result.tags)
  termination_protection_enabled                   = data.external.mongodbatlas_advanced_cluster.result.termination_protection_enabled == "" ? null : data.external.mongodbatlas_advanced_cluster.result.termination_protection_enabled
  timeouts                                         = data.external.mongodbatlas_advanced_cluster.result.timeouts == "" ? null : jsondecode(data.external.mongodbatlas_advanced_cluster.result.timeouts)
  version_release_system                           = data.external.mongodbatlas_advanced_cluster.result.version_release_system == "" ? null : data.external.mongodbatlas_advanced_cluster.result.version_release_system
}

