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

  replication_specs = local.effective_replication_specs
}
