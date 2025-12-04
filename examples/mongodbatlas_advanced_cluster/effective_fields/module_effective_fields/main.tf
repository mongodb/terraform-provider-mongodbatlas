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
# Option 1 - MIGRATION/BACKWARD COMPATIBLE (current approach, omitting flag or setting to false):
# - replication_specs returns ACTUAL provisioned values (what's currently running)
# - Maintains compatibility with module_existing behavior
# - effective_*_specs also returns actual values (always available for dedicated clusters)
# - Best for migrating from lifecycle.ignore_changes approach
#
# Option 2 - RECOMMENDED FOR NEW MODULES (set use_effective_fields = true):
# - replication_specs returns CONFIGURED values (what you specified in .tf files)
# - effective_*_specs returns ACTUAL provisioned values (may differ due to auto-scaling)
# - Clear separation between intent (configured) and reality (effective)
# - Better visibility into both client-provided and Atlas-managed values
#
# This example uses Option 1 for backward compatibility during migration.
# To use Option 2, add: use_effective_fields = true
data "mongodbatlas_advanced_cluster" "this" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
  # use_effective_fields = true  # Uncomment for Option 2 (recommended for new modules)
  depends_on = [mongodbatlas_advanced_cluster.this]
}
