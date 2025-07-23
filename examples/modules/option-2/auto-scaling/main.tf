# Validate that only one of replication_specs, shards, or region_configs is defined
locals {
  # Count how many of the three options are defined (with length > 0)
  defined_count = (
    length(var.replication_specs) > 0 ? 1 : 0
    ) + (
    length(var.shards) > 0 ? 1 : 0
    ) + (
    length(var.region_configs) > 0 ? 1 : 0
  )

}

check "validate_only_one_defined" {
  assert {
    condition     = local.defined_count <= 1
    error_message = "Only one of replication_specs, shards, or region_configs can be defined"
  }
}

locals {

  # Build the specs from replication_specs (with autoscaling), shards (geo-sharded), or region_configs (replica set)
  effective_replication_specs = (
    length(var.replication_specs) > 0 ? tolist([
      for spec in var.replication_specs : {
        zone_name = try(spec.zone_name, null)
        region_configs = [
          for region in spec.region_configs : merge(
            region,
            {
              auto_scaling = merge(
                try(region.auto_scaling, {}),
                {
                  compute_enabled = true
                  disk_gb_enabled = true
                }
              ),
              analytics_auto_scaling = merge(
                try(region.analytics_auto_scaling, {}),
                {
                  compute_enabled = true
                  disk_gb_enabled = true
                }
              ),
              electable_specs = (
                region.electable_specs == null ? null : merge(
                  region.electable_specs,
                  {
                    instance_size = try(region.auto_scaling.compute_min_instance_size, null)
                  }
                )
              ),
              read_only_specs = (
                region.read_only_specs == null ? null : merge(
                  region.read_only_specs,
                  {
                    instance_size = try(region.auto_scaling.compute_min_instance_size, null)
                  }
                )
              ),
              analytics_specs = (
                region.analytics_specs == null ? null : merge(
                  region.analytics_specs,
                  {
                    instance_size = try(region.analytics_auto_scaling.compute_min_instance_size, null)
                  }
                )
              )
            }
          )
        ]
      }
      ]) : (length(var.shards) > 0 ? tolist([
        for shard in var.shards : {
          zone_name = shard.zone_name
          region_configs = [for region in shard.region_configs : {
            provider_name = region.provider_name
            region_name   = region.region_name
            priority      = region.priority
            electable_specs = {
              instance_size   = var.auto_scaling.compute_min_instance_size
              node_count      = region.electable_node_count
              ebs_volume_type = try(region.ebs_volume_type, null)
              disk_size_gb    = null
              disk_iops       = try(region.disk_iops, null)
            }
            read_only_specs = (
              try(region.read_only_node_count, 0) > 0 ? {
                instance_size   = var.auto_scaling.compute_min_instance_size
                node_count      = region.read_only_node_count
                ebs_volume_type = try(region.ebs_volume_type, null)
                disk_size_gb    = null
                disk_iops       = try(region.disk_iops, null)
              } : null
            )
            auto_scaling = merge({
              compute_enabled = true
              disk_gb_enabled = true
            }, var.auto_scaling) # all autoscaling configs are the same cluster wide, this how API currently works
            analytics_auto_scaling = merge({
              compute_enabled = true
              disk_gb_enabled = true
            }, var.analytics_auto_scaling) # all analytics autoscaling configs are the same cluster wide, this how API currently works
            analytics_specs = (
              region.analytics_specs != null ? {
                instance_size   = var.analytics_auto_scaling.compute_min_instance_size
                node_count      = region.analytics_specs.node_count
                ebs_volume_type = try(region.analytics_specs.ebs_volume_type, null)
                disk_size_gb    = null
                disk_iops       = try(region.analytics_specs.disk_iops, null)
              } : null
            )
          }]
        }
        ]) : tolist([
        {
          zone_name = null
          region_configs = [for region in var.region_configs : {
            provider_name = region.provider_name
            region_name   = region.region_name
            priority      = region.priority
            electable_specs = {
              instance_size   = var.auto_scaling.compute_min_instance_size
              node_count      = region.electable_node_count
              ebs_volume_type = try(region.ebs_volume_type, null)
              disk_size_gb    = null
              disk_iops       = try(region.disk_iops, null)
            }
            read_only_specs = ( # read_only_specs uses same compute and storage configs as electable_specs, this is how API currently works
              try(region.read_only_node_count, 0) > 0 ? {
                instance_size   = var.auto_scaling.compute_min_instance_size
                node_count      = region.read_only_node_count
                ebs_volume_type = try(region.ebs_volume_type, null)
                disk_size_gb    = null
                disk_iops       = try(region.disk_iops, null)
              } : null
            )
            auto_scaling = merge({
              compute_enabled = true
              disk_gb_enabled = true
            }, var.auto_scaling) # all autoscaling configs are the same cluster wide, this how API currently works
            analytics_auto_scaling = merge({
              compute_enabled = true
              disk_gb_enabled = true
            }, var.analytics_auto_scaling) # all analytics autoscaling configs are the same cluster wide, this how API currently works
            analytics_specs = (
              region.analytics_specs != null ? {
                instance_size   = var.analytics_auto_scaling.compute_min_instance_size
                node_count      = region.analytics_specs.node_count
                ebs_volume_type = try(region.analytics_specs.ebs_volume_type, null)
                disk_size_gb    = null
                disk_iops       = try(region.analytics_specs.disk_iops, null)
              } : null
            )
          }]
        }
    ]))
  )
}

resource "mongodbatlas_advanced_cluster" "this" {
  project_id             = var.project_id
  name                   = var.name
  cluster_type           = var.cluster_type
  mongo_db_major_version = var.mongo_db_major_version

  replication_specs = local.effective_replication_specs

  accept_data_risks_and_force_replica_set_reconfig = var.accept_data_risks_and_force_replica_set_reconfig
  advanced_configuration                           = var.advanced_configuration
  backup_enabled                                   = var.backup_enabled
  bi_connector_config                              = var.bi_connector_config
  config_server_management_mode                    = var.config_server_management_mode
  delete_on_create_timeout                         = var.delete_on_create_timeout
  encryption_at_rest_provider                      = var.encryption_at_rest_provider
  global_cluster_self_managed_sharding             = var.global_cluster_self_managed_sharding
  paused                                           = var.paused
  pinned_fcv                                       = var.pinned_fcv
  pit_enabled                                      = var.pit_enabled
  redact_client_log_data                           = var.redact_client_log_data
  replica_set_scaling_strategy                     = var.replica_set_scaling_strategy
  retain_backups_enabled                           = var.retain_backups_enabled
  root_cert_type                                   = var.root_cert_type
  tags                                             = var.tags
  termination_protection_enabled                   = var.termination_protection_enabled
  timeouts                                         = var.timeouts
  version_release_system                           = var.version_release_system



  lifecycle {
    # Terraform cannot make the ignore_changes block fully dynamic based on input variables or locals. The list must be static and known at plan time.
    # This static list supports up to 3 shards (replication specs) with up to 3 regions
    ignore_changes = [
      // replication_specs[0]
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      replication_specs[0].region_configs[0].read_only_specs.instance_size,
      replication_specs[0].region_configs[0].analytics_specs.instance_size,
      replication_specs[0].region_configs[0].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[0].read_only_specs.disk_size_gb,
      replication_specs[0].region_configs[0].analytics_specs.disk_size_gb,

      replication_specs[0].region_configs[1].electable_specs.instance_size,
      replication_specs[0].region_configs[1].read_only_specs.instance_size,
      replication_specs[0].region_configs[1].analytics_specs.instance_size,
      replication_specs[0].region_configs[1].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[1].read_only_specs.disk_size_gb,
      replication_specs[0].region_configs[1].analytics_specs.disk_size_gb,

      replication_specs[0].region_configs[2].electable_specs.instance_size,
      replication_specs[0].region_configs[2].read_only_specs.instance_size,
      replication_specs[0].region_configs[2].analytics_specs.instance_size,
      replication_specs[0].region_configs[2].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[2].read_only_specs.disk_size_gb,
      replication_specs[0].region_configs[2].analytics_specs.disk_size_gb,

      // replication_specs[1]
      replication_specs[1].region_configs[0].electable_specs.instance_size,
      replication_specs[1].region_configs[0].read_only_specs.instance_size,
      replication_specs[1].region_configs[0].analytics_specs.instance_size,
      replication_specs[1].region_configs[0].electable_specs.disk_size_gb,
      replication_specs[1].region_configs[0].read_only_specs.disk_size_gb,
      replication_specs[1].region_configs[0].analytics_specs.disk_size_gb,

      replication_specs[1].region_configs[1].electable_specs.instance_size,
      replication_specs[1].region_configs[1].read_only_specs.instance_size,
      replication_specs[1].region_configs[1].analytics_specs.instance_size,
      replication_specs[1].region_configs[1].electable_specs.disk_size_gb,
      replication_specs[1].region_configs[1].read_only_specs.disk_size_gb,
      replication_specs[1].region_configs[1].analytics_specs.disk_size_gb,

      replication_specs[1].region_configs[2].electable_specs.instance_size,
      replication_specs[1].region_configs[2].read_only_specs.instance_size,
      replication_specs[1].region_configs[2].analytics_specs.instance_size,
      replication_specs[1].region_configs[2].electable_specs.disk_size_gb,
      replication_specs[1].region_configs[2].read_only_specs.disk_size_gb,
      replication_specs[1].region_configs[2].analytics_specs.disk_size_gb,

      // replication_specs[2]
      replication_specs[2].region_configs[0].electable_specs.instance_size,
      replication_specs[2].region_configs[0].read_only_specs.instance_size,
      replication_specs[2].region_configs[0].analytics_specs.instance_size,
      replication_specs[2].region_configs[0].electable_specs.disk_size_gb,
      replication_specs[2].region_configs[0].read_only_specs.disk_size_gb,
      replication_specs[2].region_configs[0].analytics_specs.disk_size_gb,

      replication_specs[2].region_configs[1].electable_specs.instance_size,
      replication_specs[2].region_configs[1].read_only_specs.instance_size,
      replication_specs[2].region_configs[1].analytics_specs.instance_size,
      replication_specs[2].region_configs[1].electable_specs.disk_size_gb,
      replication_specs[2].region_configs[1].read_only_specs.disk_size_gb,
      replication_specs[2].region_configs[1].analytics_specs.disk_size_gb,

      replication_specs[2].region_configs[2].electable_specs.instance_size,
      replication_specs[2].region_configs[2].read_only_specs.instance_size,
      replication_specs[2].region_configs[2].analytics_specs.instance_size,
      replication_specs[2].region_configs[2].electable_specs.disk_size_gb,
      replication_specs[2].region_configs[2].read_only_specs.disk_size_gb,
      replication_specs[2].region_configs[2].analytics_specs.disk_size_gb
    ]
  }
}