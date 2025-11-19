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
  replication_specs    = var.replication_specs
  tags                 = var.tags
}

# Data source to read effective values after Atlas auto-scales
# This is always available regardless of whether auto-scaling is enabled
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
