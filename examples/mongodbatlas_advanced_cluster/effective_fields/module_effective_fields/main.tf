# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Create Atlas Advanced Cluster with effective fields
# use_effective_fields enables:
# - Spec attributes stay constant (match your configuration)
# - Atlas-computed values available via data source effective_* attributes
# - No plan drift when Atlas auto-scales
# - No lifecycle.ignore_changes blocks needed
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true
  replication_specs    = var.replication_specs
  tags                 = var.tags
}

# Data source to read effective specs
# Exposes actual provisioned values via effective_* attributes
# (effective_electable_specs, effective_analytics_specs, effective_read_only_specs)
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
