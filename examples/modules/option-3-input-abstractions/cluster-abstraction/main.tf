# Main resource logic for cluster-abstraction module

locals {
  # Build the specs from either shards (geo-sharded) or replica_set_regions (replica set)
  effective_replication_specs = (
    length(var.shards) > 0 ? [
      for shard in var.shards : {
        zone_name = shard.zone_name
        region_configs = [for region in shard.region_configs : {
          provider_name = region.provider_name
          region_name   = region.region_name
          priority      = region.priority
          electable_specs = {
            instance_size   = region.instance_size
            node_count      = region.electable_node_count
            ebs_volume_type = try(region.ebs_volume_type, null)
            disk_size_gb    = try(region.disk_size_gb, null)
            disk_iops       = try(region.disk_iops, null)
          }
          read_only_specs = (
            try(region.read_only_node_count, 0) > 0 ? {
              instance_size   = region.instance_size
              node_count      = region.read_only_node_count
              ebs_volume_type = try(region.ebs_volume_type, null)
              disk_size_gb    = try(region.disk_size_gb, null)
              disk_iops       = try(region.disk_iops, null)
            } : null
          )
          auto_scaling           = var.auto_scaling
          analytics_auto_scaling = var.analytics_auto_scaling
          analytics_specs = (
            region.analytics_specs != null ? {
              instance_size   = region.analytics_specs.instance_size
              node_count      = region.analytics_specs.node_count
              ebs_volume_type = try(region.analytics_specs.ebs_volume_type, null)
              disk_size_gb    = try(region.analytics_specs.disk_size_gb, null)
              disk_iops       = try(region.analytics_specs.disk_iops, null)
            } : null
          )
        }]
      }
      ] : [
      {
        zone_name = null
        region_configs = [for region in var.region_configs : {
          provider_name = region.provider_name
          region_name   = region.region_name
          priority      = region.priority
          electable_specs = {
            instance_size   = region.instance_size
            node_count      = region.electable_node_count
            ebs_volume_type = try(region.ebs_volume_type, null)
            disk_size_gb    = try(region.disk_size_gb, null)
            disk_iops       = try(region.disk_iops, null)
          }
          read_only_specs = ( # read_only_specs uses same compute and storage configs as electable_specs, this is how API currently works
            try(region.read_only_node_count, 0) > 0 ? {
              instance_size   = region.instance_size
              node_count      = region.read_only_node_count
              ebs_volume_type = try(region.ebs_volume_type, null)
              disk_size_gb    = try(region.disk_size_gb, null)
              disk_iops       = try(region.disk_iops, null)
            } : null
          )
          auto_scaling           = var.auto_scaling           # all autoscaling configs are the same cluster wide, this how API currently works
          analytics_auto_scaling = var.analytics_auto_scaling # all analytics autoscaling configs are the same cluster wide, this how API currently works
          analytics_specs = (
            region.analytics_specs != null ? {
              instance_size   = region.analytics_specs.instance_size
              node_count      = region.analytics_specs.node_count
              ebs_volume_type = try(region.analytics_specs.ebs_volume_type, null)
              disk_size_gb    = try(region.analytics_specs.disk_size_gb, null)
              disk_iops       = try(region.analytics_specs.disk_iops, null)
            } : null
          )
        }]
      }
    ]
  )
}

resource "mongodbatlas_advanced_cluster" "this" {
  project_id             = var.project_id
  name                   = var.name
  cluster_type           = var.cluster_type
  mongo_db_major_version = var.mongo_db_major_version

  dynamic "replication_specs" {
    for_each = local.effective_replication_specs
    content {
      # zone_name is optional
      zone_name = lookup(replication_specs.value, "zone_name", null)
      dynamic "region_configs" {
        for_each = replication_specs.value.region_configs
        content {
          provider_name = region_configs.value.provider_name
          region_name   = region_configs.value.region_name
          priority      = region_configs.value.priority

          electable_specs {
            instance_size   = region_configs.value.electable_specs.instance_size
            node_count      = region_configs.value.electable_specs.node_count
            ebs_volume_type = region_configs.value.electable_specs.ebs_volume_type != null ? region_configs.value.electable_specs.ebs_volume_type : null
            disk_size_gb    = region_configs.value.electable_specs.disk_size_gb != null ? region_configs.value.electable_specs.disk_size_gb : null
            disk_iops       = region_configs.value.electable_specs.disk_iops != null ? region_configs.value.electable_specs.disk_iops : null
          }

          dynamic "read_only_specs" {
            for_each = region_configs.value.read_only_specs != null ? [region_configs.value.read_only_specs] : []
            content {
              instance_size   = read_only_specs.value.instance_size
              node_count      = read_only_specs.value.node_count
              ebs_volume_type = read_only_specs.value.ebs_volume_type != null ? read_only_specs.value.ebs_volume_type : null
              disk_size_gb    = read_only_specs.value.disk_size_gb != null ? read_only_specs.value.disk_size_gb : null
              disk_iops       = read_only_specs.value.disk_iops != null ? read_only_specs.value.disk_iops : null
            }
          }

          dynamic "analytics_specs" {
            for_each = region_configs.value.analytics_specs != null ? [region_configs.value.analytics_specs] : []
            content {
              instance_size   = analytics_specs.value.instance_size
              node_count      = analytics_specs.value.node_count
              ebs_volume_type = analytics_specs.value.ebs_volume_type != null ? analytics_specs.value.ebs_volume_type : null
              disk_size_gb    = analytics_specs.value.disk_size_gb != null ? analytics_specs.value.disk_size_gb : null
              disk_iops       = analytics_specs.value.disk_iops != null ? analytics_specs.value.disk_iops : null
            }
          }

          dynamic "auto_scaling" {
            for_each = region_configs.value.auto_scaling != null ? [region_configs.value.auto_scaling] : []
            content {
              disk_gb_enabled            = region_configs.value.auto_scaling.disk_gb_enabled
              compute_enabled            = region_configs.value.auto_scaling.compute_enabled
              compute_scale_down_enabled = region_configs.value.auto_scaling.compute_scale_down_enabled != null ? region_configs.value.auto_scaling.compute_scale_down_enabled : null
              compute_min_instance_size  = region_configs.value.auto_scaling.compute_min_instance_size != null ? region_configs.value.auto_scaling.compute_min_instance_size : null
              compute_max_instance_size  = region_configs.value.auto_scaling.compute_max_instance_size != null ? region_configs.value.auto_scaling.compute_max_instance_size : null
            }
          }

          dynamic "analytics_auto_scaling" {
            for_each = region_configs.value.analytics_auto_scaling != null ? [region_configs.value.analytics_auto_scaling] : []
            content {
              disk_gb_enabled            = region_configs.value.analytics_auto_scaling.disk_gb_enabled
              compute_enabled            = region_configs.value.analytics_auto_scaling.compute_enabled
              compute_scale_down_enabled = region_configs.value.analytics_auto_scaling.compute_scale_down_enabled != null ? region_configs.value.analytics_auto_scaling.compute_scale_down_enabled : null
              compute_min_instance_size  = region_configs.value.analytics_auto_scaling.compute_min_instance_size != null ? region_configs.value.analytics_auto_scaling.compute_min_instance_size : null
              compute_max_instance_size  = region_configs.value.analytics_auto_scaling.compute_max_instance_size != null ? region_configs.value.analytics_auto_scaling.compute_max_instance_size : null
            }
          }
        }
      }
    }
  }
}
