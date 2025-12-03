# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Create Atlas Advanced Cluster with auto-scaling
# This is the legacy approach requiring lifecycle.ignore_changes blocks
# to prevent Terraform from detecting drift when Atlas auto-scales the cluster
resource "mongodbatlas_advanced_cluster" "this" {
  project_id       = mongodbatlas_project.this.id
  name             = var.cluster_name
  cluster_type     = var.cluster_type
  replication_specs = var.replication_specs
  tags             = var.tags

  # lifecycle.ignore_changes is required when auto-scaling is enabled
  # to prevent Terraform from trying to revert Atlas-managed changes
  # This approach has limitations:
  # - Cannot be conditional based on auto-scaling configuration
  # - Module users cannot see actual provisioned values
  # - Requires listing all auto-scalable attributes
  lifecycle {
    ignore_changes = [
      # Ignore instance size changes from compute auto-scaling
      replication_specs[0].region_configs[0].electable_specs[0].instance_size,
      replication_specs[0].region_configs[0].analytics_specs[0].instance_size,
      replication_specs[0].region_configs[0].read_only_specs[0].instance_size,

      # Ignore disk size and IOPS changes from storage auto-scaling
      replication_specs[0].region_configs[0].electable_specs[0].disk_size_gb,
      replication_specs[0].region_configs[0].electable_specs[0].disk_iops,
      replication_specs[0].region_configs[0].analytics_specs[0].disk_size_gb,
      replication_specs[0].region_configs[0].analytics_specs[0].disk_iops,
      replication_specs[0].region_configs[0].read_only_specs[0].disk_size_gb,
      replication_specs[0].region_configs[0].read_only_specs[0].disk_iops,
    ]
  }
}
