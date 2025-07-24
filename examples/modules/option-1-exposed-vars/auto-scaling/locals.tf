locals {
  # Inject autoscaling enabled into all region_configs in replication_specs
  replication_specs_with_autoscaling = [
    for spec in var.replication_specs : {
      # Copy all fields from the original spec
      num_shards = try(spec.num_shards, null)
      zone_name  = try(spec.zone_name, null)
      region_configs = [
        for region in spec.region_configs : merge(
          region,
          {
            auto_scaling = merge(
              region.auto_scaling,
              {
                compute_enabled = true
                disk_gb_enabled = true
              }
            ),
            analytics_auto_scaling = (region.analytics_auto_scaling == null ? null : merge(
              region.analytics_auto_scaling,
              {
                compute_enabled = true
                disk_gb_enabled = true
              }
            )),
            electable_specs = (
              region.electable_specs == null ? null : merge(
                region.electable_specs,
                {
                  instance_size = region.auto_scaling.compute_min_instance_size
                }
              )
            ),
            read_only_specs = (
              region.read_only_specs == null ? null : merge(
                region.read_only_specs,
                {
                  instance_size = region.auto_scaling.compute_min_instance_size
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
  ]
}
