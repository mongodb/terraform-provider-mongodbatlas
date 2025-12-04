# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Create Atlas Advanced Cluster with use_effective_fields
#
# use_effective_fields = true on the resource:
# - Eliminates need for lifecycle.ignore_changes blocks
# - Prevents plan drift when Atlas auto-scales the cluster
# - Spec attributes in resource stay constant (match your configuration)
#
# When auto-scaling is enabled, Atlas may adjust instance_size, disk_size_gb, and disk_iops
# regardless of whether compute or disk auto-scaling is enabled (for optimal performance).
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true
  replication_specs    = var.replication_specs
  tags                 = var.tags
}

# Data source to read cluster specifications
#
# IMPORTANT: The use_effective_fields flag on the data source controls what values are returned:
#
# PHASE 1 - BACKWARD COMPATIBLE MIGRATION (current approach, omitting flag or setting to false):
# - *_specs (electable_specs, analytics_specs, read_only_specs) return ACTUAL provisioned values
# - Maintains compatibility with module_existing behavior
# - effective_*_specs also returns actual values (always available for dedicated clusters)
# - Best for migrating from lifecycle.ignore_changes approach without breaking module users
#
# PHASE 2 - BREAKING CHANGE (set use_effective_fields = true, prepares for provider v3.x):
# - *_specs (electable_specs, analytics_specs, read_only_specs) return CONFIGURED values
# - effective_*_specs returns ACTUAL provisioned values (may differ due to auto-scaling)
# - Clear separation between intent (configured) and reality (effective)
# - BREAKING: Module users must switch from *_specs to effective_*_specs for actual values
# - Prepares for provider v3.x where this becomes default behavior
# - Recommended for new modules created from scratch
#
# This example uses Phase 1 for backward compatibility during migration.
# To implement Phase 2, add: use_effective_fields = true
data "mongodbatlas_advanced_cluster" "this" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
  # use_effective_fields = true  # Uncomment for Phase 2 (breaking change, prepares for v3.x)
  depends_on = [mongodbatlas_advanced_cluster.this]
}
