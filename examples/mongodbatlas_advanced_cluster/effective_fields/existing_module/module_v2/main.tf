# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Create Atlas Advanced Cluster
# MIGRATION CHANGE: Added use_effective_fields = true
# - Specification attributes remain exactly as defined in configuration
# - Atlas-computed values available separately in effective specs
# - No plan drift when Atlas auto-scales
# - lifecycle.ignore_changes block no longer needed
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true  # NEW: Enables effective fields behavior
  replication_specs    = var.replication_specs
  tags                 = var.tags

  # MIGRATION CHANGE: lifecycle.ignore_changes block removed
}

# Data source to read effective specs
# MIGRATION CHANGE: New in v2
# Exposes actual provisioned values, including changes made by auto-scaling
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true  # Must match the resource
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
