# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Create Atlas Advanced Cluster
# use_effective_fields = true enables clear separation of concerns:
# - Specification attributes remain exactly as defined in configuration
# - Atlas-computed values available separately in effective specs
# - No plan drift when Atlas auto-scales
# - Works seamlessly for both auto-scaling and non-auto-scaling scenarios
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true
  replication_specs    = var.replication_specs
  tags                 = var.tags
}

# Data source to read effective specs
# Exposes actual provisioned values, including any changes made by auto-scaling
data "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_advanced_cluster.this.project_id
  name                 = mongodbatlas_advanced_cluster.this.name
  use_effective_fields = true
  depends_on           = [mongodbatlas_advanced_cluster.this]
}
