terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 2.0"
    }
  }
}

# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Create Atlas Advanced Cluster with optional auto-scaling
# Using use_effective_fields = true eliminates the need for lifecycle.ignore_changes
# even when auto-scaling is enabled, making this module work seamlessly in both scenarios
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true

  replication_specs = [
    for idx, spec in var.replication_specs : {
      zone_name = try(spec.zone_name, null)
      region_configs = [
        for region in spec.region_configs : {
          priority      = region.priority
          provider_name = region.provider_name
          region_name   = region.region_name

          # Auto-scaling configuration (optional)
          auto_scaling = var.enable_auto_scaling ? {
            disk_gb_enabled            = try(region.auto_scaling.disk_gb_enabled, false)
            compute_enabled            = try(region.auto_scaling.compute_enabled, false)
            compute_scale_down_enabled = try(region.auto_scaling.compute_scale_down_enabled, false)
            compute_min_instance_size  = try(region.auto_scaling.compute_min_instance_size, null)
            compute_max_instance_size  = try(region.auto_scaling.compute_max_instance_size, null)
          } : null

          # Analytics auto-scaling configuration (optional)
          analytics_auto_scaling = var.enable_analytics_auto_scaling ? {
            disk_gb_enabled            = try(region.analytics_auto_scaling.disk_gb_enabled, false)
            compute_enabled            = try(region.analytics_auto_scaling.compute_enabled, false)
            compute_scale_down_enabled = try(region.analytics_auto_scaling.compute_scale_down_enabled, false)
            compute_min_instance_size  = try(region.analytics_auto_scaling.compute_min_instance_size, null)
            compute_max_instance_size  = try(region.analytics_auto_scaling.compute_max_instance_size, null)
          } : null

          # Electable specs
          electable_specs = {
            instance_size = region.electable_specs.instance_size
            node_count    = region.electable_specs.node_count
            disk_iops     = try(region.electable_specs.disk_iops, null)
            disk_size_gb  = try(region.electable_specs.disk_size_gb, null)
          }

          # Analytics specs (optional)
          analytics_specs = try(region.analytics_specs, null) != null ? {
            instance_size = region.analytics_specs.instance_size
            node_count    = region.analytics_specs.node_count
            disk_iops     = try(region.analytics_specs.disk_iops, null)
            disk_size_gb  = try(region.analytics_specs.disk_size_gb, null)
          } : null

          # Read-only specs (optional)
          read_only_specs = try(region.read_only_specs, null) != null ? {
            instance_size = region.read_only_specs.instance_size
            node_count    = region.read_only_specs.node_count
            disk_iops     = try(region.read_only_specs.disk_iops, null)
            disk_size_gb  = try(region.read_only_specs.disk_size_gb, null)
          } : null
        }
      ]
    }
  ]

  tags = var.tags
}

# Data source to read effective values after Atlas auto-scales
# This is always available regardless of whether auto-scaling is enabled
data "mongodbatlas_advanced_cluster" "this" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
  depends_on = [mongodbatlas_advanced_cluster.this]
}
