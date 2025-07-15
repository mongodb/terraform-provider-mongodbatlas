# Main resource logic for cluster-abstraction module

locals {
  # If single_region is set, build the specs from it; else, empty list
  effective_replication_specs = (
    var.single_region != null ? [
      {
        zone_name      = null
        region_configs = [
          {
            provider_name   = var.single_region.provider_name
            region_name     = var.single_region.region_name
            priority        = 7
            electable_specs = {
              instance_size = var.single_region.instance_size
              node_count    = var.single_region.node_count
              ebs_volume_type = var.single_region.ebs_volume_type != null ? var.single_region.ebs_volume_type : null
              disk_size_gb    = var.single_region.disk_size_gb    != null ? var.single_region.disk_size_gb    : null
              disk_iops       = var.single_region.disk_iops       != null ? var.single_region.disk_iops       : null
            }
            read_only_specs = null
            auto_scaling    = var.auto_scaling
          }
        ]
      }
    ] : []
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
            instance_size = region_configs.value.electable_specs.instance_size
            node_count    = region_configs.value.electable_specs.node_count
            ebs_volume_type = region_configs.value.electable_specs.ebs_volume_type != null ? region_configs.value.electable_specs.ebs_volume_type : null
            disk_size_gb    = region_configs.value.electable_specs.disk_size_gb    != null ? region_configs.value.electable_specs.disk_size_gb    : null
            disk_iops       = region_configs.value.electable_specs.disk_iops       != null ? region_configs.value.electable_specs.disk_iops       : null
          }

          dynamic "read_only_specs" {
            for_each = region_configs.value.read_only_specs != null ? [region_configs.value.read_only_specs] : []
            content {
              instance_size = read_only_specs.value.instance_size
              node_count    = read_only_specs.value.node_count
            }
          }

          dynamic "auto_scaling" {
            for_each = region_configs.value.auto_scaling != null ? [region_configs.value.auto_scaling] : []
            content {
              disk_gb_enabled           = region_configs.value.auto_scaling.disk_gb_enabled
              compute_enabled           = region_configs.value.auto_scaling.compute_enabled
              compute_min_instance_size = region_configs.value.auto_scaling.compute_min_instance_size
              compute_max_instance_size = region_configs.value.auto_scaling.compute_max_instance_size
            }
          }
        }
      }
    }
  }
}
