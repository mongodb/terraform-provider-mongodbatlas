resource "mongodbatlas_advanced_cluster" "this" {
  project_id   = var.project_id
  name         = var.cluster_name
  cluster_type = "SHARDED"

  # v1.x schema (blocks)
  replication_specs { # first shard
    region_configs {
      auto_scaling {
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      analytics_auto_scaling {
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      electable_specs {
        instance_size = "M40"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M40"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "Zone 1"
  }

  replication_specs { # second shard
    region_configs {
      auto_scaling {
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      analytics_auto_scaling {
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      electable_specs {
        instance_size = "M40"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M40"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "Zone 1"
  }

  lifecycle { # avoid plans as autoscaling changes instance_size
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      replication_specs[0].region_configs[0].analytics_specs.instance_size,
      replication_specs[1].region_configs[0].electable_specs.instance_size,
      replication_specs[1].region_configs[0].analytics_specs.instance_size,
    ]
  }
}
