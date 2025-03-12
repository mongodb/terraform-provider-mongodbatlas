locals {
  disk_size = var.auto_scaling_disk_gb_enabled ? null : var.disk_size
  replication_specs = flatten([
    for old_spec in var.replication_specs : [
      for shard in range(old_spec.num_shards) : [
        {
          zone_name = old_spec.zone_name
          region_configs = tolist([
            for region in old_spec.regions_config : {
              region_name   = region.region_name
              provider_name = var.provider_name
              electable_specs = {
                instance_size = var.instance_size
                node_count    = region.electable_nodes
                disk_size_gb  = local.disk_size
              }
              priority = region.priority
              read_only_specs = region.read_only_nodes == 0 ? null : {
                instance_size = var.instance_size
                node_count    = region.read_only_nodes
                disk_size_gb  = local.disk_size
              }
              auto_scaling = var.auto_scaling_disk_gb_enabled ? {
                disk_gb_enabled = true
              } : null
            }
          ])
        }
      ]
    ]
    ]
  )
}

moved {
  from = mongodbatlas_cluster.this
  to   = mongodbatlas_advanced_cluster.this
}


resource "mongodbatlas_advanced_cluster" "this" {
  lifecycle {
    precondition {
      condition     = !(var.auto_scaling_disk_gb_enabled && var.disk_size > 0)
      error_message = "Must use either auto_scaling_disk_gb_enabled or disk_size, not both."
    }
  }

  project_id             = var.project_id
  name                   = var.cluster_name
  cluster_type           = var.cluster_type
  mongo_db_major_version = var.mongo_db_major_version
  replication_specs      = local.replication_specs
  tags                   = var.tags
}

data "mongodbatlas_cluster" "this" { # note the usage of `cluster` not `advanced_cluster`, this is to have outputs stay compatible with the v1 module
  name       = mongodbatlas_advanced_cluster.this.name
  project_id = mongodbatlas_advanced_cluster.this.project_id

  depends_on = [mongodbatlas_advanced_cluster.this]
}

# OLD cluster configuration:
# resource "mongodbatlas_cluster" "this" {
#   lifecycle {
#     precondition {
#       condition     = !(var.auto_scaling_disk_gb_enabled && var.disk_size > 0)
#       error_message = "Must use either auto_scaling_disk_gb_enabled or disk_size, not both."
#     }
#   }

#   project_id = var.project_id
#   name       = var.cluster_name

#   auto_scaling_disk_gb_enabled = var.auto_scaling_disk_gb_enabled
#   cluster_type                 = var.cluster_type
#   disk_size_gb                 = var.disk_size
#   mongo_db_major_version       = var.mongo_db_major_version
#   provider_instance_size_name  = var.instance_size
#   provider_name                = var.provider_name

#   dynamic "tags" {
#     for_each = var.tags
#     content {
#       key   = tags.key
#       value = tags.value
#     }
#   }

#   dynamic "replication_specs" {
#     for_each = var.replication_specs
#     content {
#       num_shards = replication_specs.value.num_shards
#       zone_name  = replication_specs.value.zone_name

#       dynamic "regions_config" {
#         for_each = replication_specs.value.regions_config
#         content {
#           electable_nodes = regions_config.value.electable_nodes
#           priority        = regions_config.value.priority
#           read_only_nodes = regions_config.value.read_only_nodes
#           region_name     = regions_config.value.region_name
#         }
#       }
#     }
#   }
# }
