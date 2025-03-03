resource "mongodbatlas_cluster" "this" {
  lifecycle {
    precondition {
      condition     = !(var.auto_scaling_disk_gb_enabled && var.disk_size > 0)
      error_message = "Must use either auto_scaling_disk_gb_enabled or disk_size, not both."
    }
  }

  project_id = var.project_id
  name       = var.cluster_name

  auto_scaling_disk_gb_enabled = var.auto_scaling_disk_gb_enabled
  cluster_type                 = var.cluster_type
  disk_size_gb                 = var.disk_size
  mongo_db_major_version       = var.mongo_db_major_version
  provider_instance_size_name  = var.instance_size
  provider_name                = var.provider_name

  dynamic "tags" {
    for_each = var.tags
    content {
      key   = tags.key
      value = tags.value
    }
  }

  dynamic "replication_specs" {
    for_each = var.replication_specs
    content {
      num_shards = replication_specs.value.num_shards
      zone_name  = replication_specs.value.zone_name

      dynamic "regions_config" {
        for_each = replication_specs.value.regions_config
        content {
          electable_nodes = regions_config.value.electable_nodes
          priority        = regions_config.value.priority
          read_only_nodes = regions_config.value.read_only_nodes
          region_name     = regions_config.value.region_name
        }
      }
    }
  }
}
