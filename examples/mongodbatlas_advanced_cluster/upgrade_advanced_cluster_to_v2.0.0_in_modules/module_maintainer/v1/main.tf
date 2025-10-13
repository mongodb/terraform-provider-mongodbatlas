resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = var.project_id
  name           = var.name
  cluster_type   = "SHARDED"
  disk_size_gb   = var.disk_size_gb
  backup_enabled = true

  replication_specs {
    # v1.x requires num_shards for sharded/geosharded
    num_shards = 2
    region_configs {
      electable_specs {
        instance_size = var.instance_size
        node_count    = 3
      }
      provider_name = var.provider_name
      region_name   = var.region_name
      priority      = 7
    }
  }

  advanced_configuration {
    javascript_enabled = true
  }

  dynamic "tags" {
    for_each = var.tags
    content {
      key   = tags.key
      value = tags.value
    }
  }
}
