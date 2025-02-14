locals {
  old_disk_size     = var.auto_scaling_disk_gb_enabled ? null : var.disk_size
  old_instance_size = coalesce(var.instance_size, "M10")
  old_provider_name = coalesce(var.provider_name, "AWS")
  replication_specs_old = flatten([
    for old_spec in var.replication_specs : [
      for shard in range(old_spec.num_shards) : [
        {
          zone_name = old_spec.zone_name
          region_configs = tolist([
            for region in old_spec.regions_config : {
              region_name   = region.region_name
              provider_name = local.old_provider_name
              electable_specs = {
                instance_size = local.old_instance_size
                node_count    = region.electable_nodes
                disk_size_gb  = local.old_disk_size
              }
              priority = region.priority
              read_only_specs = region.read_only_nodes == 0 ? null : {
                instance_size = local.old_instance_size
                node_count    = region.read_only_nodes
                disk_size_gb  = local.old_disk_size
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
  use_new_replication_specs = length(var.replication_specs_new) > 0
}
moved {
  from = mongodbatlas_cluster.this
  to   = mongodbatlas_advanced_cluster.this
}


resource "mongodbatlas_advanced_cluster" "this" {
  lifecycle {
    precondition {
      condition     = local.use_new_replication_specs || !(var.auto_scaling_disk_gb_enabled && var.disk_size > 0)
      error_message = "Must use either auto_scaling_disk_gb_enabled or disk_size, not both."
    }
    precondition {
      condition     = !((local.use_new_replication_specs && length(var.replication_specs) > 0) || (!local.use_new_replication_specs && length(var.replication_specs) == 0))
      error_message = "Must use either replication_specs_new or replication_specs, not both."
    }
  }

  project_id                  = var.project_id
  name                        = var.cluster_name
  cluster_type                = var.cluster_type
  mongo_db_major_version      = var.mongo_db_major_version
  replication_specs           = local.use_new_replication_specs ? var.replication_specs_new : local.replication_specs_old
  tags                        = var.tags
  redact_client_log_data      = true
  encryption_at_rest_provider = var.encryption_at_rest_provider

  advanced_configuration = {
    minimum_enabled_tls_protocol       = "TLS1_2"
    javascript_enabled                 = false
    tls_cipher_config_mode             = "CUSTOM"
    custom_openssl_cipher_config_tls12 = ["TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"]
  }
}

data "mongodbatlas_cluster" "this" {
  count      = local.use_new_replication_specs ? 0 : 1 # Not safe when Asymmetric Shards are used
  name       = mongodbatlas_advanced_cluster.this.name
  project_id = mongodbatlas_advanced_cluster.this.project_id

  depends_on = [mongodbatlas_advanced_cluster.this]
}

resource "mongodbatlas_search_deployment" "search_nodes" {
  count        = length(var.search_nodes_specs) > 0 ? 1 : 0
  project_id   = mongodbatlas_advanced_cluster.this.project_id
  cluster_name = mongodbatlas_advanced_cluster.this.name
  specs        = var.search_nodes_specs
}
