# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Create Atlas Advanced Cluster with use_effective_fields
# MIGRATION CHANGE: Added use_effective_fields = true
# This eliminates the need for lifecycle.ignore_changes blocks
# The module now works seamlessly for both auto-scaling and non-auto-scaling scenarios
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true  # NEW: Enables effective fields behavior
  replication_specs    = var.replication_specs
  tags                 = var.tags

  # MIGRATION CHANGE: lifecycle.ignore_changes block has been removed
  # With use_effective_fields = true, Terraform automatically handles auto-scaling drift
}

# Data source to read effective (actual) values after Atlas scales the cluster
# MIGRATION CHANGE: This data source is new in v2
# It allows module users to see the actual provisioned specifications
# including changes made by Atlas auto-scaling
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true  # Must match the resource
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
