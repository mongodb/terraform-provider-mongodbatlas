locals {
  # Inject autoscaling enabled into all region_configs in replication_specs
  replication_specs_with_autoscaling = [
    for spec in var.replication_specs : {
      # Copy all fields from the original spec
      num_shards     = try(spec.num_shards, null)
      zone_name      = try(spec.zone_name, null)
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
            )
          }
        )
      ]
    }
  ]
}
